package rbac

type RuleConstructor func() RuleInterface

type RuleFactoryInterface interface {
	Create(name string) RuleInterface
	Add(name string, fn RuleConstructor) RuleFactoryInterface
}
