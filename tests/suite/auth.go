package suite

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/stretchr/testify/require"
)

func (s *Suite) Register(
	email string,
	password string,
) *http.Response {
	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("email", email)
	form.Add("password", password)

	req, _ := http.NewRequestWithContext(s.ctx, "POST", "/users/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	s.e.ServeHTTP(w, req)

	return w.Result()
}

func (s *Suite) Login(
	email string,
	password string,
) *http.Response {
	w := httptest.NewRecorder()

	form := url.Values{}
	form.Add("email", email)
	form.Add("password", password)

	req, _ := http.NewRequestWithContext(s.ctx, "POST", "/users/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	s.e.ServeHTTP(w, req)

	return w.Result()
}

func (s *Suite) DecodeToken(
	t require.TestingT,
	r *http.Response,
) string {
	var jsonToken map[string]string
	err := json.NewDecoder(r.Body).Decode(&jsonToken)
	require.NoError(t, err)

	s.token = jsonToken["token"]
	return s.token
}
