package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/iskanye/utilities-payment-api-gateway/internal/lib/jwt"
	"github.com/iskanye/utilities-payment-api-gateway/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const PasswordLen = 10

func TestAuth_RegisterLogin_Success(t *testing.T) {
	s := suite.NewTest(t)

	email := gofakeit.Email()
	pass := randomPassword()

	// Register
	resp := s.Register(email, pass)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var jsonId map[string]int64
	err := json.NewDecoder(resp.Body).Decode(&jsonId)
	require.NoError(t, err)

	id := jsonId["id"]
	assert.NotEmpty(t, id)

	// Login
	resp = s.Login(email, pass)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var jsonToken map[string]string
	err = json.NewDecoder(resp.Body).Decode(&jsonToken)
	require.NoError(t, err)

	tokenStr := jsonToken["token"]
	tokenId, err := jwt.ValidateToken(tokenStr, s.Cfg.AuthSecret)

	require.NoError(t, err)
	assert.Equal(t, tokenId, id)
}

func BenchmarkLogin(b *testing.B) {
}

func randomPassword() string {
	return gofakeit.Password(true, true, true, true, false, PasswordLen)
}
