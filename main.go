package main

import (
	"log/slog"
	// "os"
	// "time"

	// tele "gopkg.in/telebot.v3"

	"fmt"
	translation "lang_learning_bot/translation"
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

	// =====================================================

	var textToTranslate string = "Hello, world!"

	api := translation.NewDeepLTranslator(translation.DLEn)

	translatedText, err := api.Translate(textToTranslate)
	if err != nil {
		slog.Error(err.Error())
	}

	slog.Info(fmt.Sprintf("Translated %s into %s", textToTranslate, translatedText))

	//=====================================================

	// var translator translation.Translator = translation.NewLingueeTranslator(translation.LingueeEn)

	// translatedText, _ := translator.Translate("insult")

	// fmt.Println(translatedText)

	//=====================================================

	// var translator translation.TermTranslator = translation.NewYandexDictionaryTranslator(translation.YDEn)

	// var wordToTranslate string = "carry on"
	// translatedText, err := translator.Translate(wordToTranslate)

	// if err != nil {
	// 	slog.Error(err.Error())
	// }

	// slog.Info(fmt.Sprintf("Translated %s into %s", wordToTranslate, translatedText))

}
