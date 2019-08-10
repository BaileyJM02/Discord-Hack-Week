package main

import (
  "flag"
  "github.com/bwmarrin/discordgo"
  "github.com/finione/Discord-Hack-Week/src/session"
  "github.com/finione/Discord-Hack-Week/src/util"
  "github.com/finione/Discord-Hack-Week/src/cli"
  "os"
  "os/signal"
  "strings"
  "syscall"
  "go.uber.org/zap"
)

var logger *zap.SugaredLogger
var token string
var options *cli.Options
var err error
var sessions map[string]*session.Session

func init() {
  flag.StringVar(&token, "t", "NOTOKEN", "Bot Token")
  flag.Parse()
  sessions = make(map[string]*session.Session)
}

func main() {
  // If no flags are given, run CLI
  if token == "NOTOKEN" {
    options, err = cli.Start()
  } else {
    options = &cli.Options{Token: token}
  }
  logger = util.GetSugaredLogger()
  defer logger.Sync() // flushes buffer, if any

  // Catches any errors recorded on cli.Start and fails.
  if err != nil {
      logger.Fatalf(`Unable to start due to multiple errors: %v`, err)
  }

  logger.Info("Starting bot.")

  dg, err := discordgo.New("Bot " + options.Token)
  if err != nil {
    logger.Fatalf("Error starting bot: %s", err)
  }

  dg.AddHandler(onReady)
  dg.AddHandler(onMessage)
  dg.AddHandler(onMember)

  dg.Open()

  logger.Info("Bot started. Press ctrl+c to exit.")
  sc := make(chan os.Signal, 1)
  signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
  <-sc
  logger.Info("Exiting.")
  
  defer dg.Close()

}

func onReady(client *discordgo.Session, event *discordgo.Ready) {
  logger.Info("Bot started and listening.")
  
  client.UserUpdate("", "", "", util.GetAvatar(), "")
  
  guilds, err := client.UserGuilds(100, "", "")

  if err != nil {
    logger.Fatalf("Error loading guilds: %s", err.Error())
  }

  for _, guild := range guilds {
    if guild.Owner {
      logger.Infof("Found existing session on startup: %s", guild.ID)
      // todo guild is session => load or delete, for now delete
      client.GuildDelete(guild.ID)
    }
  }

  // check if CLI has already requested a session
  if options.ServerName != "" {
    logger.Info("Using CLI session info.")
    // Not sure if this is a good work-around, but it works....
    session, err := session.Create(client, options.ServerName, options.UserID, "", options)
    if err != nil {
      // todo: error handling
    }

    sessions[session.Guild.ID] = session
  }

}

func onMessage(client *discordgo.Session, message *discordgo.MessageCreate) {
  logger.Debugf("[%s] %s: %s", message.Timestamp, message.Author.Username, message.Content)
  // todo: message handler
  if strings.HasPrefix(message.Content, "!session create") {
    name := strings.TrimSpace(strings.TrimPrefix(message.Content, "!session create"))
    if "" == name {
      name = "My Custom Guild"
    }
    
    // Create a session / guild
    session, err := session.Create(client, name, message.Author.ID, message.ChannelID, options)

    if err != nil {
      // todo: error handling
    }

    // Assign session to a map with key of guild so it can be accessed later
    sessions[session.Guild.ID] = session
  }
}


func onMember(client *discordgo.Session, member *discordgo.GuildMemberAdd) {
  // If the user who joins in the one that run the command / ID entered within CLI matches the current Guild Session
  if member.User.ID == sessions[member.GuildID].OwnerID {
    // Set the guild owner as user
    _, err := client.GuildEdit(member.GuildID, discordgo.GuildParams{OwnerID: member.User.ID})
    if err != nil {
      logger.Errorf("Error on transferring guild (%v): %v",member.GuildID, err)
    }
    // Get the bot to leave the guild
    err = client.GuildLeave(member.GuildID)
    if err != nil {
      logger.Errorf("Error on leaving guild (%v): %v",member.GuildID, err)
    }
    // Delete the session
    sessions[member.GuildID] = nil
  }
}
