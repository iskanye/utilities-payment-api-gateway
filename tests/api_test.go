package tests

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/iskanye/utilities-payment-api-gateway/tests/suite"
	"github.com/iskanye/utilities-payment-utils/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	passwordLen = 10

	adminEmail = "admin@admin.com"
	adminPass  = "admin"
	adminID    = 1

	billsN = 10

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

	tokenId, isAdmin := s.DecodeToken(t, resp)

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

	s.DecodeToken(t, resp)

	// Create bill
	address := gofakeit.Address().Address
	amount := gofakeit.Number(100, 100000)
	userID := int64(gofakeit.Number(1, 100000))

	resp = s.AddBill(address, amount, userID)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotEmpty(t, resp.Body)

	var jsonBillId map[string]int64
	err := json.NewDecoder(resp.Body).Decode(&jsonBillId)
	require.NoError(t, err)

	billID := jsonBillId["id"]
	assert.NotEmpty(t, billID)

	// Get bill
	resp = s.GetBill(billID)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotEmpty(t, resp.Body)

	var bill models.Bill
	err = json.NewDecoder(resp.Body).Decode(&bill)
	require.NoError(t, err)

	assert.Equal(t, billID, bill.ID)
	assert.Equal(t, address, bill.Address)
	assert.Equal(t, amount, bill.Amount)
	assert.Equal(t, userID, bill.UserID)

	dueDate, err := time.Parse(time.RFC3339, bill.DueDate)
	require.NoError(t, err)
	assert.InDelta(t, time.Now().AddDate(0, s.Cfg.BillingTerm, 0).Unix(), dueDate.Unix(), deltaDay)
}

func TestBilling_GetBills_Success(t *testing.T) {
	s := suite.NewTest(t)

	// Register new user
	email := gofakeit.Email()
	pass := randomPassword()
	resp := s.Register(email, pass)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotEmpty(t, resp.Body)

	var respJson map[string]int64
	err := json.NewDecoder(resp.Body).Decode(&respJson)
	require.NoError(t, err)

	userID := respJson["id"]
	require.NotEmpty(t, userID)

	// Login
	resp = s.Login(adminEmail, adminPass)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotEmpty(t, resp.Body)

	s.DecodeToken(t, resp)

	// Create bill
	var testBills [billsN]struct {
		id      int64
		address string
		amount  int
	}

	for i := range billsN {
		address := gofakeit.Address().Address
		amount := gofakeit.Number(100, 100000)

		testBills[i].address = address
		testBills[i].amount = amount

		resp = s.AddBill(address, amount, userID)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.NotEmpty(t, resp.Body)

		var jsonBillId map[string]int64
		err := json.NewDecoder(resp.Body).Decode(&jsonBillId)
		require.NoError(t, err)

		billId := jsonBillId["id"]
		require.NotEmpty(t, billId)

		testBills[i].id = billId
	}

	// Login
	resp = s.Login(email, pass)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotEmpty(t, resp.Body)

	s.DecodeToken(t, resp)

	// Get bill
	resp = s.GetBills()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotEmpty(t, resp.Body)

	var jsonBills map[string][]models.Bill
	err = json.NewDecoder(resp.Body).Decode(&jsonBills)
	require.NoError(t, err)

	bills := jsonBills["bills"]
	require.Equal(t, billsN, len(bills))

	for i := range bills {
		testBill := testBills[i]
		bill := bills[i]

		assert.Equal(t, testBill.id, bill.ID)
		assert.Equal(t, testBill.address, bill.Address)
		assert.Equal(t, testBill.amount, bill.Amount)
		assert.Equal(t, userID, bill.UserID)

		dueDate, err := time.Parse(time.RFC3339, bill.DueDate)
		require.NoError(t, err)
		assert.InDelta(t, time.Now().AddDate(0, s.Cfg.BillingTerm, 0).Unix(), dueDate.Unix(), deltaDay)
	}
}

func TestPayment_PayBill_Success(t *testing.T) {
	s := suite.NewTest(t)

	// Login
	resp := s.Login(adminEmail, adminPass)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotEmpty(t, resp.Body)

	s.DecodeToken(t, resp)

	// Create bill
	address := gofakeit.Address().Address
	amount := gofakeit.Number(100, 100000)
	userID := int64(gofakeit.Number(1, 100000))

	resp = s.AddBill(address, amount, userID)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotEmpty(t, resp.Body)

	var jsonBillId map[string]int64
	err := json.NewDecoder(resp.Body).Decode(&jsonBillId)
	require.NoError(t, err)

	billId := jsonBillId["id"]
	assert.NotEmpty(t, billId)

	// Pay bill
	resp = s.PayBill(billId)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
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

func BenchmarkAuth_GetUsers(b *testing.B) {
	s := suite.NewBench(b)

	// Login
	resp := s.Login(adminEmail, adminPass)
	s.DecodeToken(b, resp)

	for b.Loop() {
		s.GetUsers()
	}
}

func BenchmarkBilling_GetBills(b *testing.B) {
	s := suite.NewBench(b)

	// Login
	resp := s.Login(adminEmail, adminPass)
	s.DecodeToken(b, resp)

	// Create bill
	address := gofakeit.Address().Address
	amount := gofakeit.Number(100, 100000)

	s.AddBill(address, amount, adminID)
	s.AddBill(address, amount, adminID)
	s.AddBill(address, amount, adminID)

	for b.Loop() {
		s.GetBills()
	}
}

func randomPassword() string {
	return gofakeit.Password(true, true, true, true, false, passwordLen)
}
