package policy

import "github.com/gofier/framework/request/http/auth"

type key = string
type value = string
type IPolicy interface {
	Before(iu auth.IUser, routeParamMap map[key]value) bool
	Create(iu auth.IUser, routeParamMap map[key]value) bool
	Update(iu auth.IUser, routeParamMap map[key]value) bool
	Delete(iu auth.IUser, routeParamMap map[key]value) bool
	ForceDelete(iu auth.IUser, routeParamMap map[key]value) bool
	View(iu auth.IUser, routeParamMap map[key]value) bool
	Restore(iu auth.IUser, routeParamMap map[key]value) bool
}

type Action byte

const (
	ActionCreate Action = iota
	ActionUpdate
	ActionDelete
	ActionForceDelete
	ActionView
	ActionRestore
)

type Authorization struct {
	auth.RequestUser
}

func (a *Authorization) Authorize() {

}
