package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

type TomlConfig struct {
	BotToken              string
	AdminId               int
	Degug                 bool
	TorrentDownloadPath   string
	ChatIdToKickUsersFrom int64
}

var config TomlConfig
var replyHandler ReplyIdHandler

func main() {
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		log.Fatalln(err)
	}

	log.Printf("Bot token: %s\n", config.BotToken)
	log.Printf("Admin ID: %d\n", config.AdminId)
	log.Printf("TorrentDownloadPath: %s\n", config.TorrentDownloadPath)

	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	bot.Debug = true
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Bot name: %s\n", bot.Self.UserName)
	replyHandler = ReplyIdHandler{}

	botConfig := tgbotapi.UpdateConfig{0, 0, 30}

	updates, err := bot.GetUpdatesChan(botConfig)
	for update := range updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if update.Message != nil {
			torrentMessageHandler(&update, bot)
			autoKickHandler(&update, bot)
			replyHandler.process(&update, bot)
		}
	}
}

func autoKickHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	chat_id := update.Message.Chat.ID
	if chat_id == config.ChatIdToKickUsersFrom {
		if update.Message.NewChatMember != nil {
			new_chat_member := update.Message.NewChatMember
			log.Printf("New chat member on %d: %s", chat_id, new_chat_member.FirstName)
			text := fmt.Sprintf("%s we have moved to https://web.telegram.org/#/im?p=@tgbotchat", new_chat_member.FirstName)
			bot.send(tgbotapi.NewMessage(chat_id, text))
			resp, err := bot.KickChatMember(tgbotapi.ChatMemberConfig{chat_id, "", new_chat_member.ID})
			if err != nil {
				log.Fatal(err)
			}

			log.Print(resp.Result)
		}
	}
}

func torrentMessageHandler(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	fromId := update.Message.From.ID

	if fromId == config.AdminId {
		if update.Message.Text == "/loadtorrent" {

			msg := tgbotapi.NewMessage(int64(fromId), "Send me your Torrent file")
			msg.ReplyMarkup = tgbotapi.ForceReply{true, true}
			resp, err := bot.Send(msg)

			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Message sended: %d", resp.MessageID)
			replyHandler.add(resp.MessageID, func(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
				log.Printf("Received Reply from %s for %d", update.Message.From.UserName, resp.MessageID)
				if update.Message.From.ID == fromId {
					if update.Message.Document.MimeType == "application/x-bittorrent" {
						file_name := update.Message.Document.FileName
						log.Printf("Bittorrent file: %s\n", file_name)
						url, err := bot.GetFileDirectURL(update.Message.Document.FileID)
						if err != nil {
							log.Fatal(err)
						}
						log.Printf("File URL: %s\n", url)
						download_path := config.TorrentDownloadPath + "/" + file_name
						bytes, err := downloadFile(url, download_path)
						if err != nil {
							bot.Send(tgbotapi.NewMessage(int64(fromId), "Error downloading the Torrent file"))
						} else {
							text := fmt.Sprintf("Torrent `%s` (%d bytes) downloaded to `%s`", file_name, bytes, download_path)
							msg := tgbotapi.NewMessage(int64(fromId), text)
							msg.ParseMode = "markdown"
							_, err = bot.Send(msg)
							if err != nil {
								log.Fatal(err)
							}
						}
					}
				}
			})
		}
	}
}
