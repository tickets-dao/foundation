package unit

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tickets-dao/foundation/v3/core"
	"github.com/tickets-dao/foundation/v3/core/types"
	"github.com/tickets-dao/foundation/v3/core/types/big"
	"github.com/tickets-dao/foundation/v3/mock"
	"github.com/tickets-dao/foundation/v3/token"
)

const (
	testTokenName      = "Testing Token"
	testTokenSymbol    = "TT"
	testTokenCCName    = "tt"
	testTokenWithGroup = "tt_testGroup"
	testGroup          = "testGroup"

	testMessageEmptyNonce = "\"0\""

	testMSPId             = "atomyzeMSP"
	testWrongMSPId        = "wrongMSP"
	testMessageWrongMSPId = "incorrect MSP Id"

	testGetNonceFnName = "getNonce"
)

type TestToken struct {
	token.BaseToken
}

func (tt *TestToken) TxEmissionAdd(sender *types.Sender, address *types.Address, amount *big.Int) error {
	if !sender.Equal(tt.Issuer()) {
		return errors.New("unauthorized")
	}

	if amount.Cmp(big.NewInt(0)) == 0 {
		return errors.New("amount should be more than zero")
	}
	if err := tt.TokenBalanceAdd(address, amount, "txEmit"); err != nil {
		return err
	}
	return tt.EmissionAdd(amount)
}

// TestGetEmptyNonce - Checking that new wallet have empty nonce
func TestGetEmptyNonce(t *testing.T) {
	ledgerMock := mock.NewLedger(t)
	owner := ledgerMock.NewWallet()
	feeAddressSetter := ledgerMock.NewWallet()
	feeSetter := ledgerMock.NewWallet()

	tt := &TestToken{
		token.BaseToken{
			Name:     testTokenName,
			Symbol:   testTokenSymbol,
			Decimals: 8,
		},
	}

	ledgerMock.NewChainCode(testTokenCCName, tt, &core.ContractOptions{}, owner.Address(), feeSetter.Address(), feeAddressSetter.Address())

	t.Run("Get nonce with new wallet", func(t *testing.T) {
		nonce := owner.Invoke(testTokenCCName, testGetNonceFnName, owner.Address())
		assert.Equal(t, nonce, testMessageEmptyNonce)
	})
}

// TestGetNonce - Checking that the nonce after some operation is not null
func TestGetNonce(t *testing.T) {
	ledgerMock := mock.NewLedger(t)
	owner := ledgerMock.NewWallet()
	feeAddressSetter := ledgerMock.NewWallet()
	feeSetter := ledgerMock.NewWallet()

	tt := &TestToken{
		token.BaseToken{
			Name:     testTokenName,
			Symbol:   testTokenSymbol,
			Decimals: 8,
		},
	}
	ledgerMock.NewChainCode(testTokenCCName, tt, &core.ContractOptions{}, owner.Address(), feeSetter.Address(), feeAddressSetter.Address())

	owner.SignedInvoke(testTokenCCName, "emissionAdd", owner.Address(), "1000")
	owner.BalanceShouldBe(testTokenCCName, 1000)

	t.Run("Get nonce with new wallet", func(t *testing.T) {
		nonce := owner.Invoke(testTokenCCName, testGetNonceFnName, owner.Address())
		assert.NotEqual(t, nonce, testMessageEmptyNonce)
	})
}

// TestInit - Checking that init with right mspId working
func TestInit(t *testing.T) {
	ledgerMock := mock.NewLedger(t)

	tt := &TestToken{
		token.BaseToken{
			Name:     testTokenName,
			Symbol:   testTokenSymbol,
			Decimals: 8,
		},
	}

	t.Run("Init new chaincode with right MSP Id", func(t *testing.T) {
		message := ledgerMock.NewChainCodeWithCustomMSP(testTokenCCName, tt, &core.ContractOptions{}, testMSPId)
		assert.Empty(t, message)
	})
}

// TestInitMSPWithWrongMSPId - Checking that init with wrong mspId can't be done
func TestInitMSPWithWrongMSPId(t *testing.T) {
	ledgerMock := mock.NewLedger(t)

	tt := &TestToken{
		token.BaseToken{
			Name:     testTokenName,
			Symbol:   testTokenSymbol,
			Decimals: 8,
		},
	}
	t.Run("Init new chaincode with wrong MSP Id", func(t *testing.T) {
		message := ledgerMock.NewChainCodeWithCustomMSP(testTokenCCName, tt, &core.ContractOptions{}, testWrongMSPId)
		assert.Equal(t, message, testMessageWrongMSPId)
	})
}
