package rbac

type RuleInterface interface {
	Execute(userId any, item ItemInterface, context RuleContext) bool
}
