package telegram

import (
	"context"
	"fmt"
	"github.com/IrinaFosteeva/TelegramBotPocket/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) initAuthorizationProcess(message *tgbotapi.Message) error {
	authLink, err := b.generateAuthirizationLink(message.Chat.ID)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(b.messages.Start, authLink))
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) getAccessToken(chatId int64) (string, error) {
	return b.tokenRepository.Get(chatId, repository.AccessTokens)
}

func (b *Bot) generateAuthirizationLink(chatId int64) (string, error) {
	redirectUrl := b.generateRedirectUrl(chatId)
	requestToken, err := b.pocketClient.GetRequestToken(context.Background(), redirectUrl)
	if err != nil {
		return "", err
	}

	if err = b.tokenRepository.Save(chatId, requestToken, repository.RequestTokens); err != nil {
		return "", err
	}

	return b.pocketClient.GetAuthorizationURL(requestToken, redirectUrl)

}

func (b *Bot) generateRedirectUrl(chatId int64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.redirectURL, chatId)
}
