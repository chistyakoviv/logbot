package rbac

type ruleFactory struct{}

func NewRuleFactory() RuleFactoryInterface {
	return &ruleFactory{}
}

func (r *ruleFactory) Create(fn RuleConstructor) RuleInterface {
	return fn()
}
