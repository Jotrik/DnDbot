package main

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"reflect"
	_ "time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	var name tgbotapi.Update
	var pas tgbotapi.Update
	var user User
	var db *sql.DB

	db = Db_connect()
	defer db.Close()

	//bot
	bot, err := tgbotapi.NewBotAPI("1969368708:AAGYHyerxI9L4qmZ57zaeTAEU0DfvJVnGr8")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	//log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID

		if checkMesType(update) {

			switch update.Message.Text {
			case "/start":
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Hi!"))
			case "/registration":
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID,
					"Задайте логин (максимум 20 симваолов)"))

				if len(update.Message.Text) <= 20 {
					name = <-updates
					fmt.Println(name.Message.Text)
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Задайте пароль"))

				} else {
					msg_err := tgbotapi.NewMessage(update.Message.Chat.ID, "Слишком длинное имя")
					bot.Send(msg_err)
				}
				if checkMesType(update) {

					pas = <-updates
					NewUser(db, name.Message.Text, pas.Message.Text)
				} else {
					msg_err_pas := tgbotapi.NewMessage(update.Message.Chat.ID, "Используйте текст")
					bot.Send(msg_err_pas)
				}
			case "/login":
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Введите логин"))
				aname := <-updates
				user.name = aname.Message.Text
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Введите пароль"))
				apas := <-updates
				user.password = apas.Message.Text

				var tuser User
				var temp bool
				tuser, temp = Authentication(db, user)
				if temp == true {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Автроизация прошла успешно"))
				} else {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный логин или пароль"))
				}
				fmt.Println(tuser)

			}

		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Используйте текст")
			bot.Send(msg)
		}

	}
}

func checkMesType(update tgbotapi.Update) bool {
	if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
		return true
	} else {
		return false
	}
}
