package rbac

type RuleConstructor func() RuleInterface

type RuleFactoryInterface interface {
	Create(fn RuleConstructor) RuleInterface
}
