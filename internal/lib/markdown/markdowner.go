package markdown

import "strings"

type MarkdownerInterface interface {
	Escape(s string) string
}

type markdown struct {
	replacer *strings.Replacer
}

func NewMarkdowner() MarkdownerInterface {
	return &markdown{
		replacer: strings.NewReplacer(
			"_", "\\_",
			"*", "\\*",
			// "[", "\\[",
			// "]", "\\]",
			// "(", "\\(",
			// ")", "\\)",
			// "~", "\\~",
			// "`", "\\`",
			// ">", "\\>",
			// "#", "\\#",
			// "+", "\\+",
			// "-", "\\-",
			// "=", "\\=",
			// "|", "\\|",
			// "{", "\\{",
			// "}", "\\}",
			// ".", "\\.",
			// "!", "\\!",
		),
	}
}

func (m *markdown) Escape(s string) string {
	return m.replacer.Replace(s)
}
