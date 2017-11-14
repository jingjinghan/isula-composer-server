package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
)

//TODO: more logs info

// CtxErrorWrap wraps the error http message
func CtxErrorWrap(ctx *context.Context, code int, err error, msg string) {
	ctx.Output.SetStatus(code)
	ctx.Output.Body([]byte(msg))

	if err != nil {
		logs.Trace("Failed to [%s] [%s] [%d]: [%v]", ctx.Input.Method(), ctx.Input.URI(), code, err)
	} else {
		logs.Trace("Failed to [%s] [%s] [%d]", ctx.Input.Method(), ctx.Input.URI(), code)
	}
}

// CtxSuccessWrap wraps the success http message
func CtxSuccessWrap(ctx *context.Context, code int, result interface{}, header map[string]string) {
	ctx.Output.SetStatus(code)
	for n, v := range header {
		ctx.Output.Header(n, v)
	}
	output, _ := json.Marshal(result)
	ctx.Output.Body(output)

	logs.Trace("Succeed in [%s] [%s].", ctx.Input.Method(), ctx.Input.URI())
}

// CtxDataWrap wraps the http data steam
func CtxDataWrap(ctx *context.Context, code int, result []byte, header map[string]string) {
	ctx.Output.SetStatus(code)
	for n, v := range header {
		ctx.Output.Header(n, v)
	}
	ctx.Output.Body(result)

	logs.Trace("Succeed in [%s] [%s].", ctx.Input.Method(), ctx.Input.URI())
}
