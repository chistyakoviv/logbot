package i18n

import (
	"bytes"
	"fmt"
)

type i18nChain struct {
	i18n I18nInterface
	buf  bytes.Buffer
}

func NewI18nChain(i18n I18nInterface) I18nChainInterface {
	return &i18nChain{
		i18n: i18n,
		buf:  bytes.Buffer{},
	}
}

func (c *i18nChain) T(lang string, key string, opts ...Option) I18nChainInterface {
	c.buf.WriteString(c.i18n.T(lang, key, opts...))
	return c
}

func (c *i18nChain) String() string {
	return c.buf.String()
}

func (c *i18nChain) Append(s string) I18nChainInterface {
	c.buf.WriteString(s)
	return c
}

func (c *i18nChain) Appendf(format string, args ...any) I18nChainInterface {
	_, err := fmt.Fprintf(&c.buf, format, args...)
	if err != nil {
		// Branch can't be empty, but we need to ignore errors
		return c
	}
	return c
}
