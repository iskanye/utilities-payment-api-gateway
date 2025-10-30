package suite

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

func (s *Suite) AddBill(
	token string,
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
	req.Header.Set("Authorization", "Bearer "+token)

	s.e.ServeHTTP(w, req)

	return w.Result()
}

func (s *Suite) GetBills(
	token string,
	user_id int64,
) *http.Response {
	w := httptest.NewRecorder()

	vals := url.Values{}
	vals.Add("user_id", fmt.Sprint(user_id))

	req, _ := http.NewRequestWithContext(s.ctx, "GET", "/bills?"+vals.Encode(), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	s.e.ServeHTTP(w, req)

	return w.Result()
}
