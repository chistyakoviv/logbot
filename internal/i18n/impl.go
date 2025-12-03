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
				"greeting":    en.Greeting,
				"help":        en.Help,
				"intro":       en.Intro,
				"description": en.Description,

				"access_denied": en.AccessDenied,

				"subscribe_begin":         en.SubscribeBegin,
				"subscribe_invalid_token": en.SubscribeInvalidToken,
				"subscribe_token_exists":  en.SubscribeTokenExists,
				"subscribe_error":         en.SubscribeError,
				"subscribe_complete":      en.SubscribeComplete,

				"mention": en.Mention,

				"cancel_command":            en.CancelCommand,
				"cancel_command_error":      en.CancelCommandError,
				"cancel_no_current_command": en.CancelNoCurrentCommand,

				"unsubscribe_begin":            en.UnsubscribeBegin,
				"unsubscribe_invalid_token":    en.UnsubscribeInvalidToken,
				"unsubscribe_token_not_exists": en.UnsubscribeTokenNotExists,
				"unsubscribe_error":            en.UnsubscribeError,
				"unsubscribe_complete":         en.UnsubscribeComplete,

				"setlang_select_language":  en.SetLangSelectLanguage,
				"setlang_success":          en.SetLangSuccess,
				"setlang_unknown_language": en.SetLangUnknownLanguage,
				"setlang_same_language":    en.SetLangSameLanguage,
				"setlang_error":            en.SetLangError,

				"addlabels_enter_mentions":            en.AddLabelsEnterMentions,
				"addlabels_enter_labels":              en.AddLabelsEnterLabels,
				"addlabels_error":                     en.AddLabelsError,
				"addlabels_no_mentions_error":         en.AddLabelsNoMentionsError,
				"addlabels_save_mentions_error":       en.AddLabelsSaveMentionsError,
				"addlabels_no_labels_error":           en.AddLabelsNoLabels,
				"addlabels_failed_apply_labels_error": en.AddLabelsFailedApplyLabelsError,
				"addlabels_success":                   en.AddLabelsSuccess,

				"rmlabels_enter_mentions":             en.RmLabelsEnterMentions,
				"rmlabels_enter_labels":               en.RmLabelsEnterLabels,
				"rmlabels_error":                      en.RmLabelsError,
				"rmlabels_no_mentions_error":          en.RmLabelsNoMentionsError,
				"rmlabels_save_mentions_error":        en.RmLabelsSaveMentionsError,
				"rmlabels_no_labels_error":            en.RmLabelsNoLabels,
				"rmlabels_failed_remove_labels_error": en.RmLabelsFailedRemoveLabelsError,
				"rmlabels_success":                    en.RmLabelsSuccess,

				"labels_assigned": en.LabelsAssigned,
				"labels_empty":    en.LabelsEmpty,
				"labels_error":    en.LabelsError,

				"collapse_select_period": en.CollapseSelectPeriod,
				"collapse_period_set":    en.CollapsePeriodSet,

				"silence_select_period": en.SilenceSelectPeriod,
				"silence_period_set":    en.SilencePeriodSet,

				"callback_data_parse_error": en.CallbackDataParseError,
				"callback_failed_to_answer": en.CallbackFailedToAnswer,
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

func (i *I18n) Chain() I18nChainInterface {
	return NewI18nChain(i)
}

func (i *I18n) DefaultLang() string {
	return language.Default()
}

func (i *I18n) GetLangs() []string {
	return language.GetAll()
}

func (i *I18n) GetLangCode(lang string) int {
	return language.GetCode(lang)
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
