package main

import (
  "flag"
  "github.com/bwmarrin/discordgo"
  "go.uber.org/zap"
  "os"
  "os/signal"
  "syscall"
)

var token string
var sugar *zap.SugaredLogger

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
}

func onMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
  sugar.Debugf("[%s] %s: %s", message.Timestamp, message.Author.Username, message.Content)
  // todo: message handler
}
