package middleware

import (
	"errors"
	"fmt"

	"encore.dev/beta/errs"
	"encore.dev/middleware"
	"encore.dev/rlog"
)

//encore:middleware global target=all
func LogAPIErrors(req middleware.Request, next middleware.Next) middleware.Response {
	response := next(req)
	if response.Err != nil {
		var encErr = &errs.Error{}
		if errors.As(response.Err, &encErr) {
			var errorAttrs []any
			if e, ok := encErr.Meta["error"]; ok {
				errorAttrs = append(errorAttrs, "error", e)
			}
			switch {
			case response.HTTPStatus >= 500:
				rlog.Error(encErr.Message, errorAttrs...)
			case response.HTTPStatus >= 400:
				rlog.Warn(encErr.Message, errorAttrs...)
			}
		} else {
			fmt.Println("not an enc error")
			switch {
			case response.HTTPStatus >= 500:
				rlog.Error("request failed", "error", response.Err)
			case response.HTTPStatus >= 300:
				rlog.Warn("request failed", "error", response.Err)
			}
		}
	}

	return response
}
