package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// Hook defines the hook operations
type Hook struct {
	beego.Controller
}

// List gets the hook list
func (r *Hook) List() {
	user := r.Ctx.Input.Param(":user")

	logs.Debug("List Hook from '%s'", user)
	header := make(map[string]string)
	CtxDataWrap(r.Ctx, http.StatusOK, []byte("data"), header)
	return
}

// Create adds a hook to the user
func (r *Hook) Create() {
	user := r.Ctx.Input.Param(":user")
	logs.Debug("Create Hook on '%s'", user)
	CtxSuccessWrap(r.Ctx, http.StatusOK, "{}", nil)
}

// Get returns a hook information
func (r *Hook) Get() {
	user := r.Ctx.Input.Param(":user")
	id := r.Ctx.Input.Param(":id")
	logs.Debug("Get Hook '%s' from '%s'", id, user)

	CtxSuccessWrap(r.Ctx, http.StatusOK, "{}", nil)
}

// Delete removes a hook from a user
func (r *Hook) Delete() {
	user := r.Ctx.Input.Param(":user")
	id := r.Ctx.Input.Param(":id")
	logs.Debug("Delete Hook '%s' from '%s'", id, user)

	CtxSuccessWrap(r.Ctx, http.StatusOK, "{}", nil)
}
