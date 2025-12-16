package parser

type StackParser interface {
	Parse([]byte, any) ([]byte, error)
}
