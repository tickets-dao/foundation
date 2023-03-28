package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tickets-dao/foundation/v3/core"
	"github.com/tickets-dao/foundation/v3/core/types/big"
	ma "github.com/tickets-dao/foundation/v3/mock"
	"github.com/tickets-dao/foundation/v3/proto"
	pb "google.golang.org/protobuf/proto"
)

type serieSetRate struct {
	tokenName string
	dealType  string
	currency  string
	rate      string
	errorMsg  string
}

// TestSetRate - positive test with valid parameters
func TestSetRate(t *testing.T) {
	t.Parallel()

	s := &serieSetRate{
		tokenName: "tt",
		dealType:  "distribute",
		currency:  "",
		rate:      "1",
		errorMsg:  "",
	}

	BaseTokenSetRateTest(t, s)
}

// TestSetRateAllParametersAreEmpty - negative test with all parameters are empty
// result - panic
func TestSetRateAllParametersAreEmpty(t *testing.T) {
	t.Parallel()

	s := &serieSetRate{
		tokenName: "",
		dealType:  "",
		currency:  "",
		rate:      "",
		errorMsg:  "",
	}

	t.Skip("reason: https://github.com/tickets-dao/foundation/-/issues/41")
	BaseTokenSetRateTest(t, s)
}

// TestSetRateToZero - negative test with invalid rate parameter set to zero
func TestSetRateToZero(t *testing.T) {
	t.Parallel()

	s := &serieSetRate{
		tokenName: "tt",
		dealType:  "distribute",
		currency:  "",
		rate:      "0",
		errorMsg:  "trying to set rate = 0",
	}

	BaseTokenSetRateTest(t, s)
}

// TestSetRateToString - negative test with invalid rate parameter set to string
func TestSetRateToString(t *testing.T) {
	t.Parallel()

	s := &serieSetRate{
		tokenName: "tt",
		dealType:  "distribute",
		currency:  "",
		rate:      "wonder",
		errorMsg:  "couldn't convert wonder to bigint",
	}

	BaseTokenSetRateTest(t, s)
}

// TestSetRateMinusValue - negative test with invalid rate parameter set to minus value
func TestSetRateMinusValue(t *testing.T) {
	t.Parallel()

	s := &serieSetRate{
		tokenName: "tt",
		dealType:  "distribute",
		currency:  "",
		rate:      "-3",
		errorMsg:  "value -3 should be positive",
	}

	BaseTokenSetRateTest(t, s)
}

// TestSetRateSetTokenNameToWrongStringParameter - negative test with invalid token Name parameter set to wrong string
// Panic
func TestSetRateSetTokenNameToWrongStringParameter(t *testing.T) {
	t.Parallel()

	s := &serieSetRate{
		tokenName: "wonder",
		dealType:  "distribute",
		currency:  "",
		rate:      "1",
		errorMsg:  "",
	}

	t.Skip("reason: https://github.com/tickets-dao/foundation/-/issues/41")
	BaseTokenSetRateTest(t, s)
}

// TestSetRateSetTokenNameToNumericParameter - negative test with invalid token Name parameter set to numeric
// Panic
func TestSetRateSetTokenNameToNumericParameter(t *testing.T) {
	t.Parallel()

	s := &serieSetRate{
		tokenName: "353",
		dealType:  "distribute",
		currency:  "",
		rate:      "1",
		errorMsg:  "",
	}

	t.Skip("reason: https://github.com/tickets-dao/foundation/-/issues/41")
	BaseTokenSetRateTest(t, s)
}

// TestSetRateSetDealTypeToWrongstringParameter - negative test with invalid deal Type parameter set to wrong string
// ??????? err = nill, but test should be failed
func TestSetRateSetDealTypeToWrongstringParameter(t *testing.T) {
	t.Parallel()

	s := &serieSetRate{
		tokenName: "tt",
		dealType:  "wonder",
		currency:  "",
		rate:      "1",
		errorMsg:  "",
	}

	BaseTokenSetRateTest(t, s)
}

// TestSetRateSetDealTypeToNumericParameter - negative test with invalid deal Type parameter set to numeric
// ??????? err = nill, but test should be failed
func TestSetRateSetDealTypeToNumericParameter(t *testing.T) {
	t.Parallel()

	s := &serieSetRate{
		tokenName: "tt",
		dealType:  "353",
		currency:  "",
		rate:      "1",
		errorMsg:  "",
	}

	BaseTokenSetRateTest(t, s)
}

// TestSetRateCurrencyEqualToken - negative test with invalid currency parameter set to equals token
// wrong errorMSG, "is" unnecessary in this sentence.
func TestSetRateCurrencyEqualToken(t *testing.T) {
	t.Parallel()

	s := &serieSetRate{
		tokenName: "tt",
		dealType:  "distribute",
		currency:  "TT",
		rate:      "3",
		errorMsg:  "currency is equals token: it is impossible",
	}

	BaseTokenSetRateTest(t, s)
}

// TestSetRateSetCurrencyToMinusValue - negative test with invalid currency parameter set to minus value
// ??????? err = nill, but test should be failed
func TestSetRateSetCurrencyToMinusValue(t *testing.T) {
	t.Parallel()

	s := &serieSetRate{
		tokenName: "tt",
		dealType:  "distribute",
		currency:  "-10",
		rate:      "1",
		errorMsg:  "",
	}

	BaseTokenSetRateTest(t, s)
}

// TestSetRateSetCurrencyToWrongStringParameter - negative test with invalid currency parameter set to wrong string
// ??????? err = nill, but test should be failed
func TestSetRateSetCurrencyToWrongStringParameter(t *testing.T) {
	t.Parallel()

	s := &serieSetRate{
		tokenName: "tt",
		dealType:  "distribute",
		currency:  "wonder",
		rate:      "1",
		errorMsg:  "",
	}

	BaseTokenSetRateTest(t, s)
}

// TestSetRateSetCurrencyToNumericParameter - negative test with invalid currency parameter set to numeric
// ??????? err = nill, but test should be failed
func TestSetRateSetCurrencyToNumericParameter(t *testing.T) {
	t.Parallel()

	s := &serieSetRate{
		tokenName: "tt",
		dealType:  "distribute",
		currency:  "353",
		rate:      "1",
		errorMsg:  "",
	}

	BaseTokenSetRateTest(t, s)
}

// TestSetRateWrongAuthorized - negative test with invalid issuer
func TestSetRateWrongAuthorized(t *testing.T) {
	mock := ma.NewLedger(t)
	issuer := mock.NewWallet()
	outsider := mock.NewWallet()

	tt := &BaseToken{
		Name:     "Test Token",
		Symbol:   "TT",
		Decimals: 8,
	}

	mock.NewChainCode("tt", tt, &core.ContractOptions{}, issuer.Address())

	if err := outsider.RawSignedInvokeWithErrorReturned("tt", "setRate", "distribute", "", "1"); err != nil {
		assert.Equal(t, "unauthorized", err.Error())
	}
}

// TestSetRateWrongNumberParameters - negative test with incorrect number of parameters
func TestSetRateWrongNumberParameters(t *testing.T) {
	mock := ma.NewLedger(t)
	issuer := mock.NewWallet()

	tt := &BaseToken{
		Name:     "Test Token",
		Symbol:   "TT",
		Decimals: 8,
	}

	mock.NewChainCode("tt", tt, &core.ContractOptions{}, issuer.Address())

	if err := issuer.RawSignedInvokeWithErrorReturned("tt", "setRate", "distribute", "", "", ""); err != nil {
		assert.Equal(t, "incorrect number of keys or signs", err.Error())
	}
}

// BaseTokenSetRateTest - base test for checking the SetRate API
func BaseTokenSetRateTest(t *testing.T, ser *serieSetRate) {
	mock := ma.NewLedger(t)
	issuer := mock.NewWallet()
	var err error

	tt := &BaseToken{
		Name:     "Test Token",
		Symbol:   "TT",
		Decimals: 8,
	}

	mock.NewChainCode("tt", tt, &core.ContractOptions{}, issuer.Address())

	if err = issuer.RawSignedInvokeWithErrorReturned(ser.tokenName, "setRate", ser.dealType, ser.currency, ser.rate); err != nil {
		assert.Equal(t, ser.errorMsg, err.Error())
	} else {
		assert.NoError(t, err)

		data, err1 := issuer.Ledger().GetStub("tt").GetState("tokenMetadata")
		assert.NoError(t, err1)

		config := &proto.Token{}
		err2 := pb.Unmarshal(data, config)
		assert.NoError(t, err2)

		rate := config.Rates[0]
		actualRate := new(big.Int).SetBytes(rate.Rate)

		stringRate := actualRate.String()
		if ser.rate != stringRate {
			t.Errorf("Invalid rate: %s", stringRate)
		}
	}
}
