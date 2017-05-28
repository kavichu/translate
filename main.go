// Sample vision-quickstart uses the Google Cloud Vision API to label an image.
package main

import (
    "fmt"
    "log"

    // Imports the Google Cloud Translate client package.
    "cloud.google.com/go/translate"
    "golang.org/x/net/context"
    "golang.org/x/text/language"
    "encoding/json"
    "net/http"
)

type Translate struct {
    Text string
    Target_lang string
    Translations map[string]string
}

func handler (w http.ResponseWriter, req *http.Request) {
    decoder := json.NewDecoder(req.Body)
    var i18nText Translate
    err := decoder.Decode(&i18nText)
    if err != nil {
        panic(err)
    }
    print(i18nText.Text)
    defer req.Body.Close()

    ctx := context.Background()

    // Creates a client.
    client, err := translate.NewClient(ctx)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Detects the language
    langDetection, err := client.DetectLanguage(ctx, []string{i18nText.Text})
    if err != nil {
        log.Fatalf("Failed to detect the language: %v", err)
        return 
    }

    langs := []language.Tag {
        language.English,
        language.Spanish,
        language.German,
        language.BrazilianPortuguese,
        language.Chinese,
    }
    i18nText.Target_lang = langDetection[0][0].Language.String()

    resultTranslations := make(map[string]string)
    for _, lang := range langs {
        translations, err := client.Translate(ctx, []string{i18nText.Text}, lang, nil)
        if err != nil {
            log.Fatalf("Failed to translate text: %v", err)
        }else{
            fmt.Printf("Text: %v\n", i18nText.Text)
            fmt.Printf("Translation: %v\n", translations[0].Text)
            resultTranslations[lang.String()] = translations[0].Text
        }
    }

    i18nText.Translations = resultTranslations
    jsonData, err := json.Marshal(i18nText)
    if err != nil {
        log.Fatalf("Failed to parse json: %v", err)
    }

    w.Header().Set("Content-Type","application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonData)
}

func main() {
    http.HandleFunc("/translate_all", handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}