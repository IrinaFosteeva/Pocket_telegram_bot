package main

import (
	"github.com/IrinaFosteeva/TelegramBotPocket/pkg/config"
	"github.com/IrinaFosteeva/TelegramBotPocket/pkg/repository"
	"github.com/IrinaFosteeva/TelegramBotPocket/pkg/repository/boltDb"
	"github.com/IrinaFosteeva/TelegramBotPocket/pkg/server"
	"github.com/IrinaFosteeva/TelegramBotPocket/pkg/telegram"
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient(cfg.PocketConsumerKey)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDb(cfg)
	if err != nil {
		log.Fatal(err)
	}

	tokenRepository := boltDb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, cfg.TelegramBotUrl, cfg.Messages)

	authServer := server.NewAuthorizationServer(pocketClient, tokenRepository, cfg.AuthServerUrl)

	go func() {
		if err = telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err = authServer.Start(); err != nil {
		log.Fatal(err)

	}
}

func initDb(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return db, nil
}
