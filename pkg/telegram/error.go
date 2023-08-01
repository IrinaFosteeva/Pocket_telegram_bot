package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	errInvalidUrl   = errors.New("the link is not valid")
	errUnauthorised = errors.New("user id not authorized")
	errUnableToSave = errors.New("unable to save the link")
)

func (b *Bot) handleError(chatId int64, err error) {
	msg := tgbotapi.NewMessage(chatId, b.messages.Default)

	switch err {
	case errInvalidUrl:
		msg.Text = b.messages.InvalidUrl
		b.bot.Send(msg)
	case errUnauthorised:
		msg.Text = b.messages.Unauthorised
		b.bot.Send(msg)
	case errUnableToSave:
		msg.Text = b.messages.UnableToSave
		b.bot.Send(msg)
	default:
		b.bot.Send(msg)
	}
}
