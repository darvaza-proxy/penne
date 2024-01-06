package horizon

import (
	"net/http"
)

var _ http.Handler = (*Horizon)(nil)

var forwardHTTPHeaders = []string{
	"Forwarded", // rfc7239
	"X-Forwarded-For",
	"X-Forwarded-Host",
	"X-Forwarded-Proto",
}

// ServeHTTP handles HTTP requests passed from another [Horizon].
func (z *Horizon) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var next http.Handler

	switch {
	case z.next != nil:
		// hand-over to the next horizon
		next = z.next
	default:
		// EOL
		next = z.nextH
	}

	next.ServeHTTP(rw, req)
}

// HorizonServeHTTP handles HTTP requests directly from the [http.Server] when the
// client belongs in the range.
//
// A Horizon that acts as entry point has to make sure security constraints
// are checked.
func (z *Horizon) HorizonServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if !z.allowForwarding {
		hdr := rw.Header()
		for _, h := range forwardHTTPHeaders {
			hdr.Del(h)
		}
	}

	z.ServeHTTP(rw, req)
}
