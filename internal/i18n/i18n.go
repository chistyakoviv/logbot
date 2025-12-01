package i18n

type messages map[string][]string

type Option func(*I18nOpts)

type I18nChainInterface interface {
	T(lang string, key string, opts ...Option) I18nChainInterface
	Append(s string) I18nChainInterface
	Appendf(format string, args ...any) I18nChainInterface
	String() string
}

type I18nInterface interface {
	RegisterT(lang string, m messages)
	T(lang string, key string, opts ...Option) string
	DefaultLang() string
	Chain() I18nChainInterface
	GetLangs() []string
	GetLangCode(lang string) int
}
