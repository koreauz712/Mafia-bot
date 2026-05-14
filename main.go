package main

import (
	"log"
	"math/rand"
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

	bot, err := tgbotapi.NewBotAPI("8260077131:AAE66L2kNNoMAvQCDh3cfNd2X8KV9lfifNM")
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
				"🎭 Mafia Botga xush kelibsiz!\n\n/join - oyinga qo‘shilish\n/begin - o‘yinni boshlash")

		case "/join":

			alreadyJoined := false

			for _, p := range players {
				if p.ID == msg.From.ID {
					alreadyJoined = true
					break
				}
			}

			if alreadyJoined {
				send(bot, msg.Chat.ID,
					"❌ Siz allaqachon qo‘shilgansiz")
				continue
			}

			player := Player{
				ID:   msg.From.ID,
				Name: msg.From.FirstName,
			}

			players = append(players, player)

			send(bot, msg.Chat.ID,
				"✅ "+msg.From.FirstName+" oyinga qo‘shildi")

		case "/players":

			if len(players) == 0 {
				send(bot, msg.Chat.ID,
					"Hech kim qo‘shilmagan")
				continue
			}

			text := "👥 Playerlar:\n\n"

			for i, p := range players {
				text += string(rune(i+1+'0')) + ". " + p.Name + "\n"
			}

			send(bot, msg.Chat.ID, text)

		case "/begin":

			if len(players) < 3 {
				send(bot, msg.Chat.ID,
					"❌ Kamida 3 player kerak")
				continue
			}

			mafiaIndex := rand.Intn(len(players))

			for i := range players {

				if i == mafiaIndex {
					players[i].Role = "🔫 Mafia"
				} else {
					players[i].Role = "🙂 Citizen"
				}

				privateMsg := tgbotapi.NewMessage(
					players[i].ID,
					"🎭 Sizning rolingiz: "+players[i].Role,
				)

				bot.Send(privateMsg)
			}

			send(bot, msg.Chat.ID,
				"🎮 O‘yin boshlandi!\nRole private chatga yuborildi.")

		case "/reset":

			players = []Player{}

			send(bot, msg.Chat.ID,
				"♻️ Game reset qilindi")
		}
	}
}

func send(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)

	if err != nil {
		log.Println(err)
	}
}
