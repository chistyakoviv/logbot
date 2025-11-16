package i18n

import (
	"bytes"
	"fmt"
)

type I18nChain struct {
	i18n II18n
	buf  bytes.Buffer
}

func NewI18nChain(i18n II18n) II18nChain {
	return &I18nChain{
		i18n: i18n,
		buf:  bytes.Buffer{},
	}
}

func (c *I18nChain) T(lang string, key string, opts ...Option) II18nChain {
	fmt.Fprintf(&c.buf, "%s", c.i18n.T(lang, key, opts...))
	return c
}

func (c *I18nChain) String() string {
	return c.buf.String()
}
