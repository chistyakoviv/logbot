package i18n

import (
	"github.com/chistyakoviv/logbot/internal/i18n/language"
	"github.com/chistyakoviv/logbot/internal/i18n/messages/en"
	"github.com/chistyakoviv/logbot/internal/utils"
)

type messages map[string][]string

type I18n struct {
	data map[string]messages
}

func New() *I18n {
	return &I18n{
		data: map[string]messages{
			"en": {
				"greeting":    en.Greeting,
				"help":        en.Help,
				"intro":       en.Intro,
				"description": en.Description,
			},
		},
	}
}

func (i *I18n) RegisterT(lang string, m messages) {
	i.data[lang] = m
}

func (i *I18n) T(lang string, key string) string {
	translation, ok := i.data[lang]
	if !ok {
		translation, ok = i.data[i.DefaultLang()]
		if !ok {
			return "{{language not supported}}"
		}
		_, ok = translation[key]
		if !ok {
			return "{{no translation specified}}"
		}
	}
	msgs, ok := translation[key]
	if !ok {
		return "{{no translation specified}}"
	}
	return utils.RandFromSlice(msgs)
}

func (i *I18n) DefaultLang() string {
	return language.Default()
}
