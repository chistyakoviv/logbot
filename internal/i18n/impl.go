package i18n

import (
	"fmt"

	"github.com/chistyakoviv/logbot/internal/i18n/language"
	"github.com/chistyakoviv/logbot/internal/i18n/messages/en"
	"github.com/chistyakoviv/logbot/internal/utils"
)

const errLanguageNotSupported = "{{language not supported}}"
const errNoTranslationSpecified = "{{no translation specified}}"

type I18nOpts struct {
	Args []any
}

func NewI18nOpts(opts ...Option) *I18nOpts {
	o := &I18nOpts{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

type I18n struct {
	data map[string]messages
}

func New() *I18n {
	return &I18n{
		data: map[string]messages{
			"en": {
				"greeting":                en.Greeting,
				"help":                    en.Help,
				"intro":                   en.Intro,
				"description":             en.Description,
				"subscribe_begin":         en.SubscribeBegin,
				"subscribe_invalid_token": en.SubscribeInvalidToken,
				"subscribe_token_exists":  en.SubscribeTokenExists,
				"subscribe_error":         en.SubscribeError,
				"subscribe_complete":      en.SubscribeComplete,
				"mention":                 en.Mention,
			},
		},
	}
}

func (i *I18n) RegisterT(lang string, m messages) {
	i.data[lang] = m
}

func (i *I18n) T(lang string, key string, opts ...Option) string {
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

	opts = append(opts, WithDefaultArgs())
	o := NewI18nOpts(opts...)

	return fmt.Sprintf(utils.RandFromSlice(msgs), o.Args...)
}

func (i *I18n) Chain() II18nChain {
	return NewI18nChain(i)
}

func (i *I18n) DefaultLang() string {
	return language.Default()
}

func WithArgs(args []any) Option {
	return func(o *I18nOpts) {
		o.Args = args
	}
}

func WithDefaultArgs() Option {
	return func(o *I18nOpts) {
		if o.Args == nil {
			o.Args = []any{}
		}
	}
}
