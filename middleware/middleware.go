package middleware

import (
	"encoding/json"
	"errors"
	"test-server/utils"

	"github.com/valyala/fasthttp"
)

func Auth(next func(ctx *fasthttp.RequestCtx) ([]byte, error)) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Add("Content-Type", "application/json")
		token := string(ctx.Request.Header.Peek("Authorization"))
		if token == "" {
			ctx.Response.SetStatusCode(401)
			body, err := json.Marshal(utils.BadRequestError("Unauthorized"))
			if err != nil {
			}
			ctx.Response.SetBody(body)
		} else {
			userId, valid := utils.CheckToken(token)
			if valid {
				ctx.SetUserValue("session_user", userId)
				result, err := next(ctx)
				ctx.Response.SetStatusCode(200)
				if err != nil {
					var httpError *utils.HttpError
					if errors.As(err, &httpError) {
						ctx.Response.SetStatusCode(httpError.Code)
						body, err := json.Marshal(&err)
						if err != nil {
						}
						ctx.Response.SetBody(body)
					} else {
						ctx.Response.SetStatusCode(500)
						body, err := json.Marshal(utils.BadRequestError(err.Error()))
						if err != nil {
						}
						ctx.Response.SetBody(body)
					}
				} else {
					if result != nil {
						ctx.Response.SetBody(result)
					}
				}
			} else {
				ctx.Response.SetStatusCode(401)
				body, err := json.Marshal(utils.BadRequestError("Unauthorized"))
				if err != nil {
				}
				ctx.Response.SetBody(body)
			}
		}
	}
}

func Base(next func(ctx *fasthttp.RequestCtx) ([]byte, error)) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		result, err := next(ctx)
		ctx.Response.SetStatusCode(200)
		ctx.Response.Header.Add("Content-Type", "application/json")
		if err != nil {
			var httpError *utils.HttpError
			if errors.As(err, &httpError) {
				ctx.Response.SetStatusCode(httpError.Code)
				body, err := json.Marshal(&err)
				if err != nil {
				}
				ctx.Response.SetBody(body)
			} else {
				ctx.Response.SetStatusCode(500)
				body, err := json.Marshal(utils.BadRequestError(err.Error()))
				if err != nil {
				}
				ctx.Response.SetBody(body)
			}
		} else {
			if result != nil {
				ctx.Response.SetBody(result)
			}
		}
	}
}
