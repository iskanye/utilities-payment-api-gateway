package suite

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

func (s *Suite) PayBill(
	billID int64,
) *http.Response {
	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("id", fmt.Sprint(billID))

	req, _ := http.NewRequestWithContext(s.ctx, "POST", "/bills/pay", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.AddHeader(req)

	s.e.ServeHTTP(w, req)

	return w.Result()
}
