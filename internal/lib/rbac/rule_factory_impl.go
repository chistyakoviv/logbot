package rbac

type ruleFactory struct {
	constructors map[string]RuleConstructor
}

func NewRuleFactory() RuleFactoryInterface {
	return &ruleFactory{
		constructors: make(map[string]RuleConstructor),
	}
}

func (r *ruleFactory) Create(name string) RuleInterface {
	return r.constructors[name]()
}

func (r *ruleFactory) Add(name string, fn RuleConstructor) RuleFactoryInterface {
	r.constructors[name] = fn
	return r
}
