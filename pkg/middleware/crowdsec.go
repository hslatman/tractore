package middleware

import "encore.dev/middleware"

//encore:middleware global target=all
func CrowdSecBouncer(req middleware.Request, next middleware.Next) middleware.Response {
	// TODO: no way to obtain client IP? So can't block based on client IP, unfortunately.
	return next(req)
}
