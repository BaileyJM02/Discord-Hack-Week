package main

import (
	"flag"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var token string
var sugar *zap.SugaredLogger

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()

}

func main() {
	logger, _ := zap.NewProduction()
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

}

func onReady(session *discordgo.Session) {
	sugar.Info("Bot started and listening.")
}

func onMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	sugar.Debugf("[%s] %s: %s", message.Timestamp, message.Author.Username, message.Content)
	// todo: message handler
}
