package tests

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/iskanye/utilities-payment-api-gateway/internal/lib/jwt"
	"github.com/iskanye/utilities-payment-api-gateway/tests/suite"
	"github.com/iskanye/utilities-payment/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	passwordLen = 10
	adminEmail  = "admin@admin.com"
	adminPass   = "admin"

	deltaDay = 86400
)

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

	tokenStr := s.DecodeToken(t, resp)
	tokenId, isAdmin, err := jwt.ValidateToken(tokenStr, s.Cfg.AuthSecret)

	require.NoError(t, err)
	assert.Equal(t, tokenId, id)
	assert.False(t, isAdmin)
}

func TestBilling_GetBill_Success(t *testing.T) {
	s := suite.NewTest(t)

	// Login
	resp := s.Login(adminEmail, adminPass)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotEmpty(t, resp.Body)

	token := s.DecodeToken(t, resp)

	// Create bill
	address := gofakeit.Address().Address
	amount := gofakeit.Number(100, 100000)
	userID := int64(gofakeit.Number(1, 100000))

	resp = s.AddBill(token, address, amount, userID)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotEmpty(t, resp.Body)

	var jsonBillId map[string]int64
	err := json.NewDecoder(resp.Body).Decode(&jsonBillId)
	require.NoError(t, err)

	billId := jsonBillId["id"]
	assert.NotEmpty(t, billId)

	// Get bill
	resp = s.GetBills(token, userID)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotEmpty(t, resp.Body)

	var jsonBills map[string][]models.Bill
	err = json.NewDecoder(resp.Body).Decode(&jsonBills)
	require.NoError(t, err)

	bill := jsonBills["bills"][0]

	assert.Equal(t, billId, bill.ID)
	assert.Equal(t, address, bill.Address)
	assert.Equal(t, amount, bill.Amount)
	assert.Equal(t, userID, bill.UserID)

	dueDate, err := time.Parse(time.RFC3339, bill.DueDate)
	require.NoError(t, err)
	assert.InDelta(t, time.Now().AddDate(0, s.Cfg.BillingTerm, 0).Unix(), dueDate.Unix(), deltaDay)
}

// Benchmarks

func BenchmarkAuth_Login(b *testing.B) {
	s := suite.NewBench(b)

	email := gofakeit.Email()
	pass := randomPassword()

	// Register
	s.Register(email, pass)

	for b.Loop() {
		s.Login(email, pass)
	}
}

func BenchmarkBilling_GetBills(b *testing.B) {
	s := suite.NewBench(b)

	// Login
	resp := s.Login(adminEmail, adminPass)
	token := s.DecodeToken(b, resp)

	// Create bill
	address := gofakeit.Address().Address
	amount := gofakeit.Number(100, 100000)
	userID := int64(gofakeit.Number(1, 100000))

	s.AddBill(token, address, amount, userID)

	for b.Loop() {
		s.GetBills(token, userID)
	}
}

func randomPassword() string {
	return gofakeit.Password(true, true, true, true, false, passwordLen)
}
