package handlers

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"log"

	"github.com/yagop/jumble-bot/config"
)

func AutoKick(update *tgbotapi.Update, bot *tgbotapi.BotAPI, config *config.TomlConfig) {
	chat_id := update.Message.Chat.ID
	if chat_id == config.ChatIdToKickUsersFrom {
		if update.Message.NewChatMember != nil {
			new_chat_member := update.Message.NewChatMember
			log.Printf("New chat member on %d: %s", chat_id, new_chat_member.FirstName)
			text := fmt.Sprintf("%s we have moved to https://web.telegram.org/#/im?p=@tgbotchat", new_chat_member.FirstName)
			bot.Send(tgbotapi.NewMessage(chat_id, text))
			resp, err := bot.KickChatMember(tgbotapi.ChatMemberConfig{chat_id, "", new_chat_member.ID})
			if err != nil {
				log.Printf("Error kicking user")
				log.Println(err)
			}

			log.Print(resp.Result)
		}
	}
}
