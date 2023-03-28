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

type serieSetLimits struct {
	tokenName string
	dealType  string
	currency  string
	minLimit  string
	maxLimit  string
	errorMsg  string
}

// TestSetLimitsToUnlimited - positive test with maxLimit set to a valid value
func TestSetLimits(t *testing.T) {
	t.Parallel()

	s := &serieSetLimits{
		tokenName: "tt",
		dealType:  "distribute",
		currency:  "",
		minLimit:  "1",
		maxLimit:  "10",
		errorMsg:  "",
	}

	BaseTokenSetLimitsTest(t, s)
}

// TestSetLimitsToUnlimited - positive test with maxLimit set to a valid unlimited value
func TestSetLimitsToUnlimited(t *testing.T) {
	t.Parallel()

	s := &serieSetLimits{
		tokenName: "tt",
		dealType:  "distribute",
		currency:  "",
		minLimit:  "1",
		maxLimit:  "0",
		errorMsg:  "",
	}

	BaseTokenSetLimitsTest(t, s)
}

// TestSetLimitsSetMinLimitToZero - positive test with min limit parameter set to zero
func TestSetLimitsSetMinLimitToZero(t *testing.T) {
	t.Parallel()

	s := &serieSetLimits{
		tokenName: "tt",
		dealType:  "distribute",
		currency:  "",
		minLimit:  "0",
		maxLimit:  "3",
		errorMsg:  "",
	}

	BaseTokenSetLimitsTest(t, s)
}

// TestSetRateAllParametersAreEmpty - negative test with all parameters are empty
// result - panic
func TestSetLimitsAllParametersAreEmpty(t *testing.T) {
	t.Parallel()

	s := &serieSetLimits{
		tokenName: "",
		dealType:  "",
		currency:  "",
		minLimit:  "",
		maxLimit:  "",
		errorMsg:  "",
	}

	t.Skip("reason: https://github.com/tickets-dao/foundation/-/issues/42")
	BaseTokenSetLimitsTest(t, s)
}

// TestSetLimitsMinLimitGreaterMaxLimit - negative test when min limit is greater than max limit
func TestSetLimitsMinLimitGreaterMaxLimit(t *testing.T) {
	t.Parallel()

	s := &serieSetLimits{
		tokenName: "tt",
		dealType:  "distribute",
		currency:  "",
		minLimit:  "10",
		maxLimit:  "3",
		errorMsg:  "min limit is greater than max limit",
	}

	BaseTokenSetLimitsTest(t, s)
}

// TestSetLimitsMinLimitGreaterMaxLimit - negative test with invalid min limit parameter set to minus value
func TestSetLimitsMinLimitToMinusValue(t *testing.T) {
	t.Parallel()

	s := &serieSetLimits{
		tokenName: "tt",
		dealType:  "distribute",
		currency:  "",
		minLimit:  "-1",
		maxLimit:  "10",
		errorMsg:  "value -1 should be positive",
	}

	BaseTokenSetLimitsTest(t, s)
}

// TestSetLimitsMinLimitGreaterMaxLimit - negative test with invalid min limit parameter set to string
func TestSetLimitsMinLimitToString(t *testing.T) {
	t.Parallel()

	s := &serieSetLimits{
		tokenName: "tt",
		dealType:  "distribute",
		currency:  "",
		minLimit:  "wonder",
		maxLimit:  "10",
		errorMsg:  "couldn't convert wonder to bigint",
	}

	BaseTokenSetLimitsTest(t, s)
}

// TestSetLimitsMaxLimitGreaterMaxLimit - negative test with invalid max limit parameter set to minus value
func TestSetLimitsMaxLimitToMinusValue(t *testing.T) {
	t.Parallel()

	s := &serieSetLimits{
		tokenName: "tt",
		dealType:  "distribute",
		currency:  "",
		minLimit:  "1",
		maxLimit:  "-1",
		errorMsg:  "value -1 should be positive",
	}

	BaseTokenSetLimitsTest(t, s)
}

// TestSetLimitsMaxLimitGreaterMaxLimit - negative test with invalid max limit parameter set to string
func TestSetLimitsMaxLimitToString(t *testing.T) {
	t.Parallel()

	s := &serieSetLimits{
		tokenName: "tt",
		dealType:  "distribute",
		currency:  "",
		minLimit:  "1",
		maxLimit:  "wonder",
		errorMsg:  "couldn't convert wonder to bigint",
	}

	BaseTokenSetLimitsTest(t, s)
}

// TestSetLimitsSetTokenNameToWrongStringParameter - negative test with invalid token Name parameter set to wrong string
// Panic
func TestSetLimitsSetTokenNameToWrongStringParameter(t *testing.T) {
	t.Parallel()

	s := &serieSetLimits{
		tokenName: "wonder",
		dealType:  "distribute",
		currency:  "",
		minLimit:  "1",
		maxLimit:  "10",
		errorMsg:  "",
	}

	t.Skip("reason: https://github.com/tickets-dao/foundation/-/issues/42")
	BaseTokenSetLimitsTest(t, s)
}

// TestSetLimitsSetTokenNameToNumericParameter - negative test with invalid token Name parameter set to numeric
// Panic
func TestSetLimitsSetTokenNameToNumericParameter(t *testing.T) {
	t.Parallel()

	s := &serieSetLimits{
		tokenName: "353",
		dealType:  "distribute",
		currency:  "",
		minLimit:  "1",
		maxLimit:  "10",
		errorMsg:  "",
	}

	t.Skip("reason: https://github.com/tickets-dao/foundation/-/issues/42")
	BaseTokenSetLimitsTest(t, s)
}

// TestSetLimitsCurrencyEqualToken - negative test with invalid currency parameter set to equals token
func TestSetLimitsCurrencyEqualToken(t *testing.T) {
	t.Parallel()

	s := &serieSetLimits{
		tokenName: "tt",
		dealType:  "distribute",
		currency:  "TT",
		minLimit:  "1",
		maxLimit:  "10",
		errorMsg:  "unknown currency. Rate for deal type distribute and currency TT was not set",
	}

	BaseTokenSetLimitsTest(t, s)
}

// TestSetLimitsSetCurrencyToWrongStringParameter - negative test with invalid currency parameter set to wrong  string
func TestSetLimitsSetCurrencyToWrongStringParameter(t *testing.T) {
	t.Parallel()

	s := &serieSetLimits{
		tokenName: "tt",
		dealType:  "distribute",
		currency:  "wonder",
		minLimit:  "1",
		maxLimit:  "10",
		errorMsg:  "unknown currency. Rate for deal type distribute and currency wonder was not set",
	}

	BaseTokenSetLimitsTest(t, s)
}

// TestSetLimitsSetCurrencyToNumeric - negative test with invalid currency parameter set to Numeric
func TestSetLimitsSetCurrencyToNumeric(t *testing.T) {
	t.Parallel()

	s := &serieSetLimits{
		tokenName: "tt",
		dealType:  "distribute",
		currency:  "353",
		minLimit:  "1",
		maxLimit:  "10",
		errorMsg:  "unknown currency. Rate for deal type distribute and currency 353 was not set",
	}

	BaseTokenSetLimitsTest(t, s)
}

// TestSetLimitsSetDealTypeToWrongstringParameter - negative test with invalid deal Type parameter set to wrong string
// failed because the extra space between "currency" and "was"
func TestSetLimitsSetDealTypeToWrongstringParameter(t *testing.T) {
	t.Parallel()

	s := &serieSetLimits{
		tokenName: "tt",
		dealType:  "wonder",
		currency:  "",
		minLimit:  "1",
		maxLimit:  "10",
		errorMsg:  "unknown DealType. Rate for deal type wonder and currency  was not set",
	}

	BaseTokenSetLimitsTest(t, s)
}

// TestSetLimitsSetDealTypeToNumericParameter - negative test with invalid deal Type parameter set to numeric
// failed because the extra space between "currency" and "was"
func TestSetLimitsSetDealTypeToNumericParameter(t *testing.T) {
	t.Parallel()

	s := &serieSetLimits{
		tokenName: "tt",
		dealType:  "353",
		currency:  "",
		minLimit:  "1",
		maxLimit:  "10",
		errorMsg:  "unknown DealType. Rate for deal type 353 and currency  was not set",
	}

	BaseTokenSetLimitsTest(t, s)
}

// TestSetLimitsWrongNumberParameters - negative test with incorrect number of parameters
func TestSetLimitsWrongNumberParameters(t *testing.T) {
	mock := ma.NewLedger(t)
	issuer := mock.NewWallet()

	tt := &BaseToken{
		Name:     "Test Token",
		Symbol:   "TT",
		Decimals: 8,
	}

	mock.NewChainCode("tt", tt, &core.ContractOptions{}, issuer.Address())

	issuer.SignedInvoke("tt", "setRate", "distribute", "", "1")

	if err := issuer.RawSignedInvokeWithErrorReturned("tt", "setLimits", "distribute", "", "", "1", "10"); err != nil {
		assert.Equal(t, "incorrect number of keys or signs", err.Error())
	}
}

// BaseTokenSetLimitsTest - base test for checking the SetLimitse API
func BaseTokenSetLimitsTest(t *testing.T, ser *serieSetLimits) {
	mock := ma.NewLedger(t)
	issuer := mock.NewWallet()
	var err error

	tt := &BaseToken{
		Name:     "Test Token",
		Symbol:   "TT",
		Decimals: 8,
	}

	mock.NewChainCode("tt", tt, &core.ContractOptions{}, issuer.Address())

	issuer.SignedInvoke("tt", "setRate", "distribute", "", "1")

	if err = issuer.RawSignedInvokeWithErrorReturned(ser.tokenName, "setLimits", ser.dealType, ser.currency, ser.minLimit, ser.maxLimit); err != nil {
		assert.Equal(t, ser.errorMsg, err.Error())
	} else {
		assert.NoError(t, err)

		data, err1 := issuer.Ledger().GetStub("tt").GetState("tokenMetadata")
		assert.NoError(t, err1)

		config := &proto.Token{}
		err2 := pb.Unmarshal(data, config)
		assert.NoError(t, err2)

		rate := config.Rates[0]
		actualMinLimit := new(big.Int).SetBytes(rate.Min)
		actualMaxLimit := new(big.Int).SetBytes(rate.Max)
		stringMinLimit := actualMinLimit.String()
		stringMaxLimit := actualMaxLimit.String()

		if ser.minLimit != stringMinLimit {
			t.Errorf("Invalid Min Limit: %s", stringMinLimit)
		}
		if ser.maxLimit != stringMaxLimit {
			t.Errorf("Invalid Min Limit: %s", stringMaxLimit)
		}
	}
}
