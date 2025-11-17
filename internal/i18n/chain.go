package i18n

import (
	"bytes"
)

type i18nChain struct {
	i18n II18n
	buf  bytes.Buffer
}

func NewI18nChain(i18n II18n) II18nChain {
	return &i18nChain{
		i18n: i18n,
		buf:  bytes.Buffer{},
	}
}

func (c *i18nChain) T(lang string, key string, opts ...Option) II18nChain {
	c.buf.WriteString(c.i18n.T(lang, key, opts...))
	return c
}

func (c *i18nChain) String() string {
	return c.buf.String()
}

func (c *i18nChain) Append(s string) II18nChain {
	c.buf.WriteString(s)
	return c
}
