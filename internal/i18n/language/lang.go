package language

type Language int

const (
	En = iota
	Ru
)

const (
	en = "en"
	ru = "ru"
)

func (lang Language) String() string {

	switch lang {
	case En:
		return en
	case Ru:
		return ru
	default:
		return en
	}
}

func Default() string {
	return Language(0).String()
}
