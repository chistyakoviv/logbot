package rbac

type RuleInterface interface {
	Execute(userId any, item Item, context RuleContext) bool
}
