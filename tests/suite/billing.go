package suite

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

func (s *Suite) AddBill(
	address string,
	amount int,
	user_id int64,
) *http.Response {
	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("address", address)
	form.Add("amount", fmt.Sprint(amount))
	form.Add("user_id", fmt.Sprint(user_id))

	req, _ := http.NewRequestWithContext(s.ctx, "POST", "/admin/bills", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.AddHeader(req, s.UserID)

	s.e.ServeHTTP(w, req)

	return w.Result()
}

func (s *Suite) GetBills(
	user_id int64,
) *http.Response {
	w := httptest.NewRecorder()

	vals := url.Values{}
	vals.Add("user_id", fmt.Sprint(user_id))

	req, _ := http.NewRequestWithContext(s.ctx, "GET", "/bills", nil)
	s.AddHeader(req, user_id)

	s.e.ServeHTTP(w, req)

	return w.Result()
}

func (s *Suite) GetBill(
	bill_id int64,
) *http.Response {
	w := httptest.NewRecorder()

	req, _ := http.NewRequestWithContext(s.ctx, "GET", "/bills/"+fmt.Sprint(bill_id), nil)
	s.AddHeader(req, s.UserID)

	s.e.ServeHTTP(w, req)

	return w.Result()
}
