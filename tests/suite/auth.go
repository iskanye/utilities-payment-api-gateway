package suite

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/iskanye/utilities-payment-api-gateway/internal/lib/jwt"
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

func (s *Suite) GetUsers() *http.Response {
	w := httptest.NewRecorder()

	req, _ := http.NewRequestWithContext(s.ctx, "POST", "/admin/bills", nil)
	s.AddHeader(req)

	s.e.ServeHTTP(w, req)

	return w.Result()
}

func (s *Suite) DecodeToken(
	t require.TestingT,
	r *http.Response,
) (int64, bool) {
	var jsonToken map[string]string
	err := json.NewDecoder(r.Body).Decode(&jsonToken)
	require.NoError(t, err)

	s.token = jsonToken["token"]

	userID, isAdmin, err := jwt.ValidateToken(s.token, s.Cfg.AuthSecret)
	require.NoError(t, err)

	s.UserID = userID
	return s.UserID, isAdmin
}

func (s *Suite) AddHeader(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("UserID", fmt.Sprint(s.UserID))
}
