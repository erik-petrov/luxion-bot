package main

import (
	"context"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const uri = "mongodb://user:pass@sample.host:27017/?maxPoolSize=20&w=majority"

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("LUXION-BOT"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://koteman123:" + os.Getenv("password") + "@cluster0.83jlpjn.mongodb.net/?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	coll := client.Database("myDB").Collection("luxion-users")

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			if update.Message.Text == "/start" {
				if (coll.FindOne(ctx, bson.D{{"ID", update.Message.Chat.ID}}) == nil) {
					doc := bson.D{{"ID", update.Message.Chat.ID}, {"send", false}}
					coll.InsertOne(ctx, doc)
				}
				msg.Text = "Добро пожаловать в Люксион-бота, впишите /get чтобы узнать текущие последние предложения Люксиона."
			} else if update.Message.Text == "/get" {
				data := Scrape()
				text := Join(data)
				msg.Text = "Текущие предложения Люксиона: " + text
			} else {
				msg.Text = "ниче не понял"
			}
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}
