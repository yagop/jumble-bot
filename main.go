package main

import (
	"github.com/tucnak/telebot"
	"log"
	"os"
	"time"
	"strconv"
	"net/http"
	"io"
	"fmt"
)

var BOT_TOKEN string
var ADMIN_ID int
var TORRENT_DOWNLOAD_PATH string
var err error

func main() {
	BOT_TOKEN := os.Getenv("BOT_TOKEN")
	ADMIN_ID, err = strconv.Atoi(os.Getenv("ADMIN_ID"))
	TORRENT_DOWNLOAD_PATH = os.Getenv("TORRENT_DOWNLOAD_PATH")
	
	if err != nil {
		log.Fatalln("Error parsing ADMIN_ID_str")
		log.Fatalln(err)
	}

	bot, err := telebot.NewBot(BOT_TOKEN)
	if err != nil {
		log.Fatalln(err)
	}

	bot.Messages = make(chan telebot.Message, 100)
	go messages(bot)
	bot.Start(1 * time.Second)

}

func messages(bot *telebot.Bot) {
	for message := range bot.Messages {
		torrentMessageHandler(&message, bot)
		autoKicker(&message, bot)
	}
}

func downloadFile(url string, download_path string) (int64, error) {
	out, err := os.Create(download_path)
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer resp.Body.Close()
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	log.Println("Written bytes", n)
	return n, nil
}

func autoKicker(msg *telebot.Message, bot *telebot.Bot) {
	if (msg.IsService() && msg.UserJoined != telebot.User{}) {
		log.Println("Service join message", msg.UserJoined)
		log.Println("Service join message chat", msg.Chat)
		
		if msg.Chat.ID == -21005536 {
			log.Println("User ID:", msg.UserJoined.ID)
		}
	}
}

func torrentMessageHandler(msg *telebot.Message, bot *telebot.Bot) {
	if msg.IsPersonal() && msg.Sender.ID != ADMIN_ID {
		log.Println("Unauthorized message", msg.Sender)
	} else {
		log.Println("Authorized message", msg.Sender, msg.Text)
		// Check file type
		if (msg.Document != telebot.Document{}) {
			log.Println("Document message", msg.Document)
			if msg.Document.Mime == "application/x-bittorrent" {
				log.Println("Bittorrent file", msg.Document.File)
				url, err := bot.GetFileDirectURL(msg.Document.File.FileID)
				if err != nil {
					log.Println(err)
					return
				}
				log.Println("Torrent url", url)
				timestamp := strconv.FormatInt(time.Now().Unix(), 10)
				download_path := TORRENT_DOWNLOAD_PATH + "/" + timestamp + "." + msg.Document.FileName
				downloadFile(url, download_path)
			}
		}
	}
}
