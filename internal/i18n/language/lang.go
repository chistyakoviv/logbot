package language

type Language int

const (
	En = iota + 1
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
	return Language(En).String()
}
