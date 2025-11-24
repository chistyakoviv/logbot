package rbac

const (
	AND = iota
	OR
)

type compositeRule struct {
	operator  int
	ruleNames []string
}

func NewCompositeRule(operator int, ruleNames []string) RuleInterface {
	if operator != AND && operator != OR {
		panic("unknown operator for composite rule")
	}
	return &compositeRule{
		operator:  operator,
		ruleNames: ruleNames,
	}
}

func (c *compositeRule) Execute(userId any, item ItemInterface, context RuleContext) bool {
	if len(c.ruleNames) == 0 {
		return true
	}

	var result bool
	for _, ruleName := range c.ruleNames {
		result = context.CreateRule(ruleName).Execute(userId, item, context)

		if c.operator == AND && !result {
			return false
		}

		if c.operator == OR && result {
			return true
		}
	}

	return c.operator == AND
}
