package rbac

type RuleFactoryInterface interface {
	Create(name string) RuleInterface
}
