package tracker

import (
	"context"
	"fmt"

	"encore.app/mercure"
	"encore.app/pkg/events"
	"encore.app/pkg/signedurl"
	"encore.dev"
	"encore.dev/beta/errs"
	"encore.dev/metrics"
	"encore.dev/rlog"
)

var EmailsTracked = metrics.NewCounter[uint64]("emails_tracked", metrics.CounterConfig{})

//encore:service
type Service struct {
}

type TrackParams struct {
	ID        int    `query:"id"`
	Signature string `query:"signature"`
	Expiry    string `query:"expiry"`
}

type TrackResponse struct {
	ID          int    `header:"X-TRACTORE-ID"`
	ContentType string `header:"Content-Type"`
}

//encore:api public method=GET path=/t
func (s *Service) Track(ctx context.Context, params *TrackParams) (*TrackResponse, error) {
	id := params.ID
	signature := params.Signature
	expiry := params.Expiry

	rlog.Info("tracking request", "signature", signature)

	// always require a signature
	if signature == "" {
		return nil, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "signature missing",
		}
	}

	baseURL := encore.Meta().APIBaseURL.String()
	reconstructedURL := fmt.Sprintf("%s/t?id=%d&signature=%s&expiry=%s", baseURL, id, signature, expiry)

	if err := signedurl.Verify(reconstructedURL); err != nil {
		rlog.Debug("invalid signature", "signature", signature, "url", reconstructedURL)
		return nil, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: "invalid signature",
		}
	}

	e := &events.MercureMessage{
		ID:    id,
		State: "tracked",
	}

	if err := mercure.PublishTrack(ctx, e); err != nil {
		return nil, &errs.Error{
			Code:    errs.Internal,
			Message: "internal error",
			Meta: errs.Metadata{
				"error": err,
			},
		}
	}

	EmailsTracked.Increment()

	// TODO: use a raw endpoint instead, to return a transparent pixel?
	return &TrackResponse{ID: params.ID, ContentType: "application/octet-stream"}, nil
}

// NOTE: below code is the raw API endpoint that returns an actual 1px
// transparent GIF. Encore does not seem to support returning binary
// responses natively yet, but it looks like specifying a 1x1 px width
// and height does the job of the image not being shown if it can't be
// reached, so the actual binary response doesn't need to be returned,
// it seems. The code for rendering the pixel was largely taken from
// https://gist.github.com/RobinUS2/7639902.

// const transPixel = "\x47\x49\x46\x38\x39\x61\x01\x00\x01\x00\x80\x00\x00\x00\x00\x00\x00\x00\x00\x21\xF9\x04\x01\x00\x00\x00\x00\x2C\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x02\x44\x01\x00\x3B"

// //encore:api public raw method=GET path=/t
// func (s *Service) Track(w http.ResponseWriter, req *http.Request) {
// 	// TODO: check environment and host?

// 	q := req.URL.Query()
// 	id := q.Get("id") // TODO: validate ID; track it internally
// 	nid, err := strconv.Atoi(id)
// 	if err != nil {
// 		rlog.Error("failed parsing email ID", "error", err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write(nil)
// 	}

// 	ctx := context.Background()

// 	e := &events.TrackEmail{
// 		ID: nid,
// 	}

// 	if err := mercure.PublishTrack(ctx, e); err != nil {
// 		rlog.Error("publishing tracked email ID to mercure failed", "error", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write(nil)
// 	}

// 	w.Header().Set("Content-Type", "image/gif")
// 	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
// 	w.Header().Set("Expires", "Mon, 06 May 2024 01:01:01 GMT")
// 	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
// 	w.Header().Set("Pragma", "no-cache")
// 	w.WriteHeader(200)
// 	_, _ = fmt.Fprint(w, transPixel)
// }
