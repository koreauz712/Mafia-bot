package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Player struct {
	ID   int64
	Name string
	Role string
}

var players []Player

func main() {

	token := os.Getenv("BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Bot started: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	rand.Seed(time.Now().UnixNano())

	for update := range updates {

		if update.Message == nil {
			continue
		}

		msg := update.Message

		switch msg.Text {

		case "/start":
			send(bot, msg.Chat.ID,
				"Welcome to Mafia Bot!")

		case "/join":

			player := Player{
				ID:   msg.From.ID,
				Name: msg.From.FirstName,
			}

			players = append(players, player)

			send(bot, msg.Chat.ID,
				msg.From.FirstName+" joined the game")

		case "/begin":

			if len(players) < 3 {
				send(bot, msg.Chat.ID,
					"Need at least 3 players")
				continue
			}

			mafiaIndex := rand.Intn(len(players))

			for i := range players {

				if i == mafiaIndex {
					players[i].Role = "Mafia"
				} else {
					players[i].Role = "Citizen"
				}

				privateMsg := tgbotapi.NewMessage(
					players[i].ID,
					"Your role: "+players[i].Role,
				)

				bot.Send(privateMsg)
			}

			send(bot, msg.Chat.ID,
				"Game started!")
		}
	}
}

func send(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}
