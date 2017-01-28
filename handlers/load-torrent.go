package handlers

import (
	"log"
  "fmt"
	"gopkg.in/telegram-bot-api.v4"

	"github.com/yagop/jumble-bot/config"
  "github.com/yagop/jumble-bot/utils"
)

func LoadTorrent(update *tgbotapi.Update, bot *tgbotapi.BotAPI, replyHandler *ReplyIdHandler, config *config.TomlConfig) {
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
			replyHandler.Add(resp.MessageID, func(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
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
						bytes, err := utils.DownloadFile(url, download_path)
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
