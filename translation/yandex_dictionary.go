package translation

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	YDEn = "en"
	YDDe = "de"
	YDRu = "ru"
)

const YandexDictionaryTranslateURL = "https://dictionary.yandex.net/api/v1/dicservice.json/lookup"

type YandexDictionaryTranslator struct {
	SrcLang string
	DstLang string
	APIKey  string
}

func NewYandexDictionaryTranslator(SrcLang string) *YandexDictionaryTranslator {
	ydApiToken := os.Getenv("YD_API_KEY")
	if ydApiToken == "" {
		panic("Provide Yandex Dictionary API key")
	}

	return &YandexDictionaryTranslator{
		SrcLang: SrcLang,
		DstLang: YDRu,
		APIKey:  ydApiToken,
	}
}

type YandexDictionaryLookupResponse struct {
	Def []struct {
		OriginalText string `json:"text"`
		Pos          string `json:"pos"`
		Translations []struct {
			Text string `json:"text"`
		} `json:"tr"`
	} `json:"def"`
}

func (t *YandexDictionaryTranslator) generateYandexDictionaryUrl(word string) string {
	word = strings.ReplaceAll(word, " ", "%20")
	query := strings.Join(
		[]string{
			fmt.Sprintf("key=%s", t.APIKey),
			fmt.Sprintf("lang=%s-%s", t.SrcLang, t.DstLang),
			fmt.Sprintf("text=%s", word),
		},
		"&",
	)
	return fmt.Sprintf("%s?%s", YandexDictionaryTranslateURL, query)
}

func convertYandexTranslationToTranslatedTerm(word string, yt *YandexDictionaryLookupResponse) *TranslatedTerm {
	termTranslation := make([]TermTranslation, len(yt.Def))
	for idx, definition := range yt.Def {
		termTranslation[idx].Pos = definition.Pos
		termTranslation[idx].Vars = make([]string, len(definition.Translations))
		for jdx, tr := range definition.Translations {
			termTranslation[idx].Vars[jdx] = tr.Text
		}
	}
	return &TranslatedTerm{
		Text:         word,
		Translations: termTranslation,
	}
}

func (t *YandexDictionaryTranslator) Translate(word string) (*TranslatedTerm, error) {
	var err error

	defer func() {
		if err != nil {
			err = fmt.Errorf("YandexDictionaryTranslator: Error %w occured", err)
		}
	}()

	url := t.generateYandexDictionaryUrl(word)

	fmt.Printf("YDUrl = %s\n", url)

	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var translations YandexDictionaryLookupResponse

	err = json.Unmarshal(data, &translations)
	if err != nil {
		return nil, nil
	}

	fmt.Printf("Yandex Translation: %s - %s\n", word, translations)

	return convertYandexTranslationToTranslatedTerm(word, &translations), nil
}
