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

const DEEPL_TRANSLATE_URL = "https://api-free.deepl.com/v2/translate"

type DeelpTranslatorAPI struct {
	Url string
	Token string
	SrcLang string
	DstLang string
}

func NewDeeplTranslatorAPI(srcLang string) *DeelpTranslatorAPI {
	deeplApiToken := os.Getenv("DEEPL_API_TOKEN")
	if deeplApiToken == "" {
		panic("Provide deepl api token")
	}

	return &DeelpTranslatorAPI{
		Url: DEEPL_TRANSLATE_URL,
		Token: deeplApiToken,
		SrcLang: srcLang,
		DstLang: DLRu,
	}
}

type DeepLResponse struct {	
	Translations []struct{
		DetectedSourceLang string `json:"detected_source_language"`
		Text string `json:"text"`
	} `json:"translations"`
}

type DeepLRequest struct {
	Text []string `json:"text"`
	TargetLanguage string `json:"target_lang"`
	SourceLanguage string `json:"source_lang"`
}

func NewDeepLRequest(
	txt string,
	api *DeelpTranslatorAPI,
) DeepLRequest {
	return DeepLRequest{
		Text: []string{txt},
		TargetLanguage: api.DstLang,
		SourceLanguage: api.SrcLang,
	}
}

func parseDeeplTranslationResponse(data []byte) (string, error) {
	var response DeepLResponse

	err := json.Unmarshal(data, &response)
	if err != nil {
		return "", nil
	}
	fmt.Println("Deepl API response", response)
	return response.Translations[0].Text, nil
}

func (api *DeelpTranslatorAPI) Translate(seq string) (string, error) {
	var err error

	defer func(){
		if err != nil {
			err = fmt.Errorf("error %w occured", err)
		}
	}()

	body, err := json.Marshal(NewDeepLRequest(seq, api))
	if err != nil {
		return "", err
	}

	fmt.Println(string(body))

	request, err := http.NewRequest("POST", api.Url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	request.Header.Set("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", api.Token))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	data, err := io.ReadAll(response.Body)
	defer func(){
		response.Body.Close()
	}()
	
	fmt.Println(string(data))

	request.Body.Close()
	
	translatedText, err := parseDeeplTranslationResponse(data)
	if err != nil {
		return "", err
	}

	return translatedText, nil
}

