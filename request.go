package clearbank

import (
	"net/http"
)

// request is a wrapper around http.Request.
type request struct {
	req              *http.Request
	expectedStatuses []int
	decodeTo         interface{}
}

func NewRequest(r *http.Request) *request {
	return &request{req: r}
}

func (r *request) ExpectStatus(expected ...int) {
	r.expectedStatuses = expected
}

func (r *request) DecodeTo(to interface{}) {
	r.decodeTo = to
}

// AddQueryParam adds a query parameter to the request.
func (r *request) AddQueryParam(key, value string) {
	r.AddQueryParams(map[string]string{key: value})
}

// AddQueryParams adds multiple query parameters to the request.
func (r *request) AddQueryParams(params map[string]string) {
	q := r.req.URL.Query()
	for k := range params {
		q.Add(k, params[k])
	}
	r.req.URL.RawQuery = q.Encode()
}
