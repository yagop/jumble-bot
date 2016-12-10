package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

// Function callback type declaration
type MessageCallbackFunction func(update *tgbotapi.Update, bot *tgbotapi.BotAPI)

// Struct containing the original message Id and a callback to be executed
type MessageReplyCall struct {
	messageId int
	callback  MessageCallbackFunction
}

type ReplyIdHandler struct {
	Callbacks []MessageReplyCall
}

// Register a callback to be executed when arrives a message with ID messageId
func (handler *ReplyIdHandler) add(messageId int, callaback MessageCallbackFunction) {
	log.Printf("Registered callback for message: %d", messageId)
	handler.Callbacks = append(handler.Callbacks, MessageReplyCall{messageId, callaback})
}

// Execute the callbacks registered on messageReplyCallRegister
func (handler *ReplyIdHandler) process(update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message.ReplyToMessage != nil {
		repliedMessageID := update.Message.ReplyToMessage.MessageID
		log.Printf("ReplyToMessage: %d", repliedMessageID)
		for _, callback := range handler.Callbacks {
			if repliedMessageID == callback.messageId {
				log.Printf("Calling callback: %x", callback.callback)
				callback.callback(update, bot)
			}
		}
	}
}
