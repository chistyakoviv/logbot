package i18n

type messages map[string][]string

type Option func(*I18nOpts)

type II18nChain interface {
	T(lang string, key string, opts ...Option) II18nChain
	Append(s string) II18nChain
	String() string
}

type II18n interface {
	RegisterT(lang string, m messages)
	T(lang string, key string, opts ...Option) string
	DefaultLang() string
	Chain() II18nChain
	GetLangs() []string
}
