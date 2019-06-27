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
)

var token string
var options cli.Options
var err error
var sessions []*session.Session

func init() {
  flag.StringVar(&token, "t", "NOTOKEN", "Bot Token")
  flag.Parse()

  // If no flags are given, run CLI
  if token == "NOTOKEN" {
    options, err = cli.Start()
  } else {
    options = cli.Options{Token: token}
  }
}

func main() {
  logger := util.GetSugaredLogger()
  defer logger.Sync() // flushes buffer, if any

  // Catches any errors recorded on init and fails.
  if err != nil {
      logger.Fatalf(`Unable to start due to multiple errors: %v`, err)
  }

  logger.Info("Starting bot.")

  dg, err := discordgo.New("Bot " + token)
  if err != nil {
    logger.Fatalf("Error starting bot: %s", err)

  }

  dg.AddHandler(onReady)
  dg.AddHandler(onMessage)

  dg.Open()

  logger.Info("Bot started. Press ctrl+c to exit.")
  sc := make(chan os.Signal, 1)
  signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
  <-sc

  dg.Close()

}

func onReady(client *discordgo.Session, event *discordgo.Ready) {
  logger := util.GetSugaredLogger()
  logger.Info("Bot started and listening.")
  
  client.UserUpdate("", "", "", util.GetAvatar(), "")
  
  guilds, err := client.UserGuilds(100, "", "")

  if err != nil {
    logger.Fatalf("Error loading guilds: %s", err.Error())
  }

  for _, guild := range guilds {
    if guild.Owner {
      logger.Infof("Found existing session on startup: %s", guild.ID)
      // todo guild is session => load or delete
    }
  }

  // check if CLI has already requested a session
  if options.ServerName != "" {
    logger.Info("Using CLI session info.")
    var user *discordgo.User
    user = &discordgo.User{
      ID: options.UserID,
      Bot: false,
    }
    // Not sure if this is a good work-around, but it works....
    userpointer := *user
    session, err := session.Create(client, options.ServerName, &userpointer, options)
    if err != nil {
      // todo: error handling
    }
    sessions = append(sessions, session)
  }

}

func onMessage(client *discordgo.Session, message *discordgo.MessageCreate) {
  logger := util.GetSugaredLogger()
  logger.Debugf("[%s] %s: %s", message.Timestamp, message.Author.Username, message.Content)
  // todo: message handler
  if strings.HasPrefix(message.Content, "!session create") {
    name := strings.TrimSpace(strings.TrimPrefix(message.Content, "!session create"))
    if "" == name {
      name = "My Custom Guild"
    }
    session, err := session.Create(client, name, message.Author, options)

    if err != nil {
      // todo: error handling
    }

    sessions = append(sessions, session)
  }
}
