package translation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	DLRu = "RU"
	DLEn = "EN"
	DLDe = "DE"
)

const DeepLTranslateURL = "https://api-free.deepl.com/v2/translate"

type DeepLTranslator struct {
	Url     string
	Token   string
	SrcLang string
	DstLang string
}

func NewDeepLTranslator(srcLang string) *DeepLTranslator {
	deeplApiToken := os.Getenv("DEEPL_API_TOKEN")
	if deeplApiToken == "" {
		panic("Provide deepl api token")
	}

	return &DeepLTranslator{
		Url:     DeepLTranslateURL,
		Token:   deeplApiToken,
		SrcLang: srcLang,
		DstLang: DLRu,
	}
}

type DeepLResponse struct {
	Translations []struct {
		DetectedSourceLang string `json:"detected_source_language"`
		Text               string `json:"text"`
	} `json:"translations"`
}

type DeepLRequest struct {
	Text           []string `json:"text"`
	TargetLanguage string   `json:"target_lang"`
	SourceLanguage string   `json:"source_lang"`
}

func NewDeepLRequest(
	txt string,
	api *DeepLTranslator,
) DeepLRequest {
	return DeepLRequest{
		Text:           []string{txt},
		TargetLanguage: api.DstLang,
		SourceLanguage: api.SrcLang,
	}
}

func parseDeepLTranslationResponse(data []byte) (string, error) {
	var response DeepLResponse

	err := json.Unmarshal(data, &response)
	if err != nil {
		return "", nil
	}
	fmt.Println("Deepl API response", response)
	return response.Translations[0].Text, nil
}

func (api *DeepLTranslator) Translate(seq string) (*TranslatedSequence, error) {
	var err error

	defer func() {
		if err != nil {
			err = fmt.Errorf("DeepLTranslator: Error %w occured", err)
		}
	}()

	body, err := json.Marshal(NewDeepLRequest(seq, api))
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", api.Url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", api.Token))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(response.Body)
	defer response.Body.Close()

	fmt.Println(string(data))

	request.Body.Close()

	translatedText, err := parseDeepLTranslationResponse(data)
	if err != nil {
		return nil, err
	}

	return &TranslatedSequence{Text: seq, Translation: translatedText}, nil
}
