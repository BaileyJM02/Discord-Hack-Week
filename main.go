package main

import (
  "flag"
  "github.com/bwmarrin/discordgo"
  "go.uber.org/zap"
  "os"
  "os/signal"
  "strings"
  "swiss.dev/Discord-Hack-Week/src/session"
  "syscall"
)

var token string
var sugar *zap.SugaredLogger
var sessions []*session.Session

func init() {
  flag.StringVar(&token, "t", "", "Bot Token")
  flag.Parse()

}

func main() {
  logger, _ := zap.NewDevelopment()
  defer logger.Sync() // flushes buffer, if any
  sugar = logger.Sugar()


  if token == "" {
    sugar.Fatal("No token provided. Please use -t <bot token>")
  }

  sugar.Info("Starting bot.")

  dg, err := discordgo.New("Bot " + token)
  if err != nil {
    sugar.Fatalf("Error starting bot: %s", err)

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

func onReady(session *discordgo.Session, event *discordgo.Ready) {
  sugar.Info("Bot started and listening.")

  guilds, err := session.UserGuilds(100, "", "")

  if err != nil {
    sugar.Fatalf("Error loading guilds: %s", err.Error())
  }

  for _, guild := range guilds {
    if guild.Owner {
      sugar.Infof("Found existing session on startup: %s", guild.ID)
      // todo guild is session => load or delete
    }
  }

}

func onMessage(client *discordgo.Session, message *discordgo.MessageCreate) {
  sugar.Debugf("[%s] %s: %s", message.Timestamp, message.Author.Username, message.Content)
  // todo: message handler
  if strings.HasPrefix(message.Content, "!session create") {
    name := strings.TrimSpace(strings.TrimPrefix(message.Content, "!session create"))
    if "" == name {
      name = "My Custom Guild"
    }
    session, err := session.Create(client, name, message.Author)

    if err != nil {
      // todo: error handling
    }

    sessions = append(sessions, session)
  }
}
