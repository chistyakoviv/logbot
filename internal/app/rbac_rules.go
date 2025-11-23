package app

import (
	"fmt"

	"github.com/chistyakoviv/logbot/internal/rbac"
)

type superuserRule struct{}

func NewSuperuserRule() rbac.RuleInterface {
	return &superuserRule{}
}

func (r *superuserRule) Execute(userId any, item rbac.ItemInterface, context rbac.RuleContext) bool {
	fmt.Println("Superuser rule executed")
	return true
}
