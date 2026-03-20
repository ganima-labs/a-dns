package i18n

import (
	"embed"
	"encoding/json"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localesFS embed.FS

var (
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
	initOnce  sync.Once
)

func init() {
	_ = Init("fr")
}

func Init(lang string) error {
	var initErr error
	initOnce.Do(func() {
		bundle = i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

		files, err := localesFS.ReadDir("locales")
		if err != nil {
			initErr = err
			return
		}

		for _, f := range files {
			if err := loadLocale(f.Name()); err != nil {
				initErr = err
				return
			}
		}
	})
	if initErr != nil {
		return initErr
	}

	localizer = i18n.NewLocalizer(bundle, lang, "en")
	return nil
}

func loadLocale(filename string) error {
	data, err := localesFS.ReadFile("locales/" + filename)
	if err != nil {
		return err
	}
	_, err = bundle.ParseMessageFileBytes(data, filename)
	return err
}

func T(id string) string {
	if localizer == nil {
		return id
	}
	msg, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: id})
	if err != nil {
		return id
	}
	return msg
}

func TWithData(id string, data map[string]interface{}) string {
	if localizer == nil {
		return id
	}
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: data,
	})
	if err != nil {
		return id
	}
	return msg
}

func SetLanguage(lang string) {
	if bundle != nil {
		localizer = i18n.NewLocalizer(bundle, lang, "en")
	}
}

func GetSupportedLanguages() []string {
	return []string{"en", "fr", "es"}
}
