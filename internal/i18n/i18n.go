package i18n

import (
	"fmt"

	"github.com/chistyakoviv/logbot/internal/i18n/language"
	"github.com/chistyakoviv/logbot/internal/i18n/messages/en"
	"github.com/chistyakoviv/logbot/internal/utils"
)

const errLanguageNotSupported = "{{language not supported}}"
const errNoTranslationSpecified = "{{no translation specified}}"

type messages map[string][]string

type I18n struct {
	data map[string]messages
}

func New() *I18n {
	return &I18n{
		data: map[string]messages{
			"en": {
				"greeting":              en.Greeting,
				"help":                  en.Help,
				"intro":                 en.Intro,
				"description":           en.Description,
				"subscribe_begin":       en.SubscribeBegin,
				"subscribe_empty_token": en.SubscribeEmptyToken,
				"subscribe_error":       en.SubscribeError,
				"subscribe_complete":    en.SubscribeComplete,
			},
		},
	}
}

func (i *I18n) RegisterT(lang string, m messages) {
	i.data[lang] = m
}

func (i *I18n) T(lang string, key string, args ...interface{}) string {
	translation, ok := i.data[lang]
	if !ok {
		translation, ok = i.data[i.DefaultLang()]
		if !ok {
			return errLanguageNotSupported
		}
		_, ok = translation[key]
		if !ok {
			return errNoTranslationSpecified
		}
	}
	msgs, ok := translation[key]
	if !ok {
		return errNoTranslationSpecified
	}
	return fmt.Sprintf(utils.RandFromSlice(msgs), args...)
}

func (i *I18n) DefaultLang() string {
	return language.Default()
}
