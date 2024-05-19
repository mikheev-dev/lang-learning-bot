package main

import (
	// "log/slog"
	// "os"
	// "time"

	// tele "gopkg.in/telebot.v3"

	"fmt"
	translation "lang_learning_bot/translation"
	"log/slog"
)

func main() {
	// pref := tele.Settings{
	// 	Token:  os.Getenv("ENGLISH_ODYSSEY_BOT_TOKEN"),
	// 	Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	// }

	// b, err := tele.NewBot(pref)
	// if err != nil {
	// 	slog.Error("Error: ", err)
	// 	return
	// }

	// b.Handle("/hello", func(c tele.Context) error {
	// 	return c.Send("Hello!")
	// })

	// b.Start()

	var textToTranslate string = "Hello, world!"

	api := translation.NewDeeplTranslatorAPI(translation.DLEn)

	translatedText, err := api.Translate(textToTranslate)
	if err != nil {
		slog.Error(err.Error())
	}

	slog.Info(fmt.Sprintf("Translated %s into %s", textToTranslate, translatedText))
}
