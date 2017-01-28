package main

import (
	"github.com/BurntSushi/toml"
	"gopkg.in/telegram-bot-api.v4"
	"log"

	"github.com/yagop/jumble-bot/config"
	"github.com/yagop/jumble-bot/handlers"
)

var tomlConfig config.TomlConfig
var replyHandler handlers.ReplyIdHandler

func main() {
	if _, err := toml.DecodeFile("config.toml", &tomlConfig); err != nil {
		log.Fatalln(err)
	}

	log.Printf("Bot token: %s\n", tomlConfig.BotToken)
	log.Printf("Admin ID: %d\n", tomlConfig.AdminId)
	log.Printf("TorrentDownloadPath: %s\n", tomlConfig.TorrentDownloadPath)

	bot, err := tgbotapi.NewBotAPI(tomlConfig.BotToken)
	bot.Debug = true
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Bot name: %s\n", bot.Self.UserName)
	replyHandler = handlers.ReplyIdHandler{}

	botConfig := tgbotapi.UpdateConfig{0, 0, 30}

	updates, err := bot.GetUpdatesChan(botConfig)
	for update := range updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if update.Message != nil {
			handlers.LoadTorrent(&update, bot, &replyHandler, &tomlConfig)
			handlers.AutoKick(&update, bot, &tomlConfig)
			replyHandler.Process(&update, bot)
		}
	}
}
