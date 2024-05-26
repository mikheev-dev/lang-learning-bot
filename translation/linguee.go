package translation

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

const LingueeTranslateURL = "https://www.linguee.com/"

const (
	LingueeRu = "RU"
	LingueeEn = "EN"
	LingueeDe = "DE"
)

var LingueeMap map[string]string = map[string]string{
	LingueeRu: "russian",
	LingueeEn: "english",
	LingueeDe: "deutsch",
}

type LingueeTranslator struct {
	SrcLang string
	DstLang string
}

func NewLingueeTranslator(SrcLang string) *LingueeTranslator {
	return &LingueeTranslator{
		SrcLang,
		LingueeRu,
	}
}

func (t *LingueeTranslator) generateLingueeUrl(word string) string {
	params := []string{
		fmt.Sprintf("source=%s", t.SrcLang),
		fmt.Sprintf("query=%s", word),
		"ajax=1",
	}

	query := strings.Join(params, "&")

	dictionary := fmt.Sprintf("%s-%s", LingueeMap[t.SrcLang], LingueeMap[t.DstLang])

	return fmt.Sprintf("%s/%s/search?%s", LingueeTranslateURL, dictionary, query)
}

func (t *LingueeTranslator) Translate(word string) (string, error) {
	var err error

	defer func() {
		if err != nil {
			err = fmt.Errorf("LingueeTranslator: Error %w occured", err)
		}
	}()

	url := t.generateLingueeUrl(word)

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	pageToParse, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(pageToParse), nil
}
