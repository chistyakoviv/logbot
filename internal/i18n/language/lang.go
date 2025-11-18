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

const UnknownLanguage = 0

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

func GetAll() []string {
	return []string{ru, en}
}

func GetCode(lang string) int {
	switch lang {
	case en:
		return En
	case ru:
		return Ru
	default:
		return UnknownLanguage
	}
}
