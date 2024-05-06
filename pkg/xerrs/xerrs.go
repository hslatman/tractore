package xerrs

import (
	"encore.dev/beta/errs"
	"encore.dev/rlog"
)

func Internal(err error) *errs.Error {
	rlog.Error("internal error", "error", err)
	return &errs.Error{
		Code:    errs.Internal,
		Message: "internal error",
		Meta: errs.Metadata{
			"error": err,
		},
	}
}
