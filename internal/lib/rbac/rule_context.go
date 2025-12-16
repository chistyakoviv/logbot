package rbac

type RuleContext struct {
	ruleFactory RuleFactoryInterface
	parameters  RuleContextParameters
}

type RuleContextParameters map[string]any

func NewRuleContext(ruleFactory RuleFactoryInterface, parameters RuleContextParameters) RuleContext {
	return RuleContext{
		ruleFactory: ruleFactory,
		parameters:  parameters,
	}
}

func (r *RuleContext) GetParameters() RuleContextParameters {
	return r.parameters
}

func (r *RuleContext) GetParameterValue(name string) any {
	_, ok := r.parameters[name]
	if !ok {
		return nil
	}
	return r.parameters[name]
}

func (r *RuleContext) HasParameter(name string) bool {
	_, ok := r.parameters[name]
	return ok
}

func (r *RuleContext) CreateRule(name string) RuleInterface {
	return r.ruleFactory.Create(name)
}
