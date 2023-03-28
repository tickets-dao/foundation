package unit

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tickets-dao/foundation/v3/core"
	"github.com/tickets-dao/foundation/v3/core/types"
	"github.com/tickets-dao/foundation/v3/core/types/big"
	"github.com/tickets-dao/foundation/v3/mock"
	"github.com/tickets-dao/foundation/v3/token"
)

func (tt *TestToken) TxIndustrialBalanceAdd(_ *types.Sender, token string, address *types.Address, amount *big.Int, reason string) error {
	return tt.IndustrialBalanceAdd(token, address, amount, reason)
}

func (tt *TestToken) QueryIndustrialBalanceGet(address *types.Address) (map[string]string, error) {
	return tt.IndustrialBalanceGet(address)
}

func (tt *TestToken) TxIndustrialBalanceSub(_ *types.Sender, token string, address *types.Address, amount *big.Int, reason string) error {
	return tt.IndustrialBalanceSub(token, address, amount, reason)
}

func (tt *TestToken) TxIndustrialBalanceTransfer(_ *types.Sender, token string, from *types.Address, to *types.Address, amount *big.Int, reason string) error {
	return tt.IndustrialBalanceTransfer(token, from, to, amount, reason)
}

func (tt *TestToken) TxIndustrialBalanceLock(_ *types.Sender, token string, address *types.Address, amount *big.Int) error {
	return tt.IndustrialBalanceLock(token, address, amount)
}

func (tt *TestToken) QueryIndustrialBalanceGetLocked(address *types.Address) (map[string]string, error) {
	return tt.IndustrialBalanceGetLocked(address)
}

func (tt *TestToken) TxIndustrialBalanceUnLock(_ *types.Sender, token string, address *types.Address, amount *big.Int) error {
	return tt.IndustrialBalanceUnLock(token, address, amount)
}

func (tt *TestToken) TxIndustrialBalanceTransferLocked(_ *types.Sender, token string, from *types.Address, to *types.Address, amount *big.Int, reason string) error {
	return tt.IndustrialBalanceTransferLocked(token, from, to, amount, reason)
}

func (tt *TestToken) TxIndustrialBalanceBurnLocked(_ *types.Sender, token string, address *types.Address, amount *big.Int, reason string) error {
	return tt.IndustrialBalanceBurnLocked(token, address, amount, reason)
}

// TestIndustrialBalanceAdd - Checking that industrial balance can be added
func TestIndustrialBalanceAdd(t *testing.T) {
	ledgerMock := mock.NewLedger(t)
	owner := ledgerMock.NewWallet()
	feeAddressSetter := ledgerMock.NewWallet()
	feeSetter := ledgerMock.NewWallet()
	user := ledgerMock.NewWallet()

	balanceAddAmount := "123"

	tt := &TestToken{
		token.BaseToken{
			Name:     testTokenWithGroup,
			Symbol:   testTokenSymbol,
			Decimals: 8,
		},
	}

	ledgerMock.NewChainCode(testTokenWithGroup, tt, &core.ContractOptions{}, owner.Address(), feeSetter.Address(), feeAddressSetter.Address())

	t.Run("Industrial balance add", func(t *testing.T) {
		owner.SignedInvoke(testTokenWithGroup, "industrialBalanceAdd", testTokenWithGroup, user.Address(), balanceAddAmount, "add balance")

		balanceResponse := owner.Invoke(testTokenWithGroup, "industrialBalanceGet", user.Address())
		balance, err := GetIndustrialBalanceFromResponseByGroup(balanceResponse, testGroup)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, balanceAddAmount, balance, "assert that balance equals "+balanceAddAmount)
	})
}

// TestIndustrialBalanceSub - Checking that industrial balance sub is working
func TestIndustrialBalanceSub(t *testing.T) {
	ledgerMock := mock.NewLedger(t)
	owner := ledgerMock.NewWallet()
	feeAddressSetter := ledgerMock.NewWallet()
	feeSetter := ledgerMock.NewWallet()
	user := ledgerMock.NewWallet()

	balanceAddAmount := "123"
	subAmount := "23"
	balanceAfterSubExpected := "100"

	tt := &TestToken{
		token.BaseToken{
			Name:     testTokenWithGroup,
			Symbol:   testTokenSymbol,
			Decimals: 8,
		},
	}

	ledgerMock.NewChainCode(testTokenWithGroup, tt, &core.ContractOptions{}, owner.Address(), feeSetter.Address(), feeAddressSetter.Address())

	owner.SignedInvoke(testTokenWithGroup, "industrialBalanceAdd", testTokenWithGroup, user.Address(), balanceAddAmount, "add balance for "+balanceAddAmount)
	owner.Invoke(testTokenWithGroup, "industrialBalanceGet", user.Address())

	t.Run("Industrial balance sub", func(t *testing.T) {
		owner.SignedInvoke(testTokenWithGroup, "industrialBalanceSub", testTokenWithGroup, user.Address(), subAmount, "sub balance for "+subAmount)
		balanceAfterSubResponse := owner.Invoke(testTokenWithGroup, "industrialBalanceGet", user.Address())
		balanceAfterSub, err := GetIndustrialBalanceFromResponseByGroup(balanceAfterSubResponse, testGroup)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, balanceAfterSubExpected, balanceAfterSub, "assert that balance equals "+balanceAfterSubExpected)
	})
}

// TestIndustrialBalanceTransfer - Checking that industrial balance transfer is working
func TestIndustrialBalanceTransfer(t *testing.T) {
	ledgerMock := mock.NewLedger(t)
	owner := ledgerMock.NewWallet()
	feeAddressSetter := ledgerMock.NewWallet()
	feeSetter := ledgerMock.NewWallet()

	balanceAddAmount := "123"
	transferAmount := "122"
	balanceAfterTransferUser1Expected := "1"
	balanceAfterTransferUser2Expected := "122"

	tt := &TestToken{
		token.BaseToken{
			Name:     testTokenWithGroup,
			Symbol:   testTokenSymbol,
			Decimals: 8,
		},
	}

	user1 := ledgerMock.NewWallet()
	user2 := ledgerMock.NewWallet()

	ledgerMock.NewChainCode(testTokenWithGroup, tt, &core.ContractOptions{}, owner.Address(), feeSetter.Address(), feeAddressSetter.Address())

	owner.SignedInvoke(testTokenWithGroup, "industrialBalanceAdd", testTokenWithGroup, user1.Address(), balanceAddAmount, "add balance for "+balanceAddAmount)
	owner.Invoke(testTokenWithGroup, "industrialBalanceGet", user1.Address())

	t.Run("Industrial balance transfer", func(t *testing.T) {
		user1.SignedInvoke(testTokenWithGroup, "industrialBalanceTransfer", testGroup, user1.Address(), user2.Address(), transferAmount, "transfer balance for "+transferAmount)

		balanceAfterSubUser2Response := owner.Invoke(testTokenWithGroup, "industrialBalanceGet", user2.Address())
		balanceAfterSubUser2, err := GetIndustrialBalanceFromResponseByGroup(balanceAfterSubUser2Response, testGroup)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, balanceAfterTransferUser2Expected, balanceAfterSubUser2, "assert that balance equals "+balanceAfterTransferUser2Expected)

		balanceAfterSubUser1Response := owner.Invoke(testTokenWithGroup, "industrialBalanceGet", user1.Address())
		balanceAfterSubUser1, err := GetIndustrialBalanceFromResponseByGroup(balanceAfterSubUser1Response, testGroup)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, balanceAfterTransferUser1Expected, balanceAfterSubUser1, "assert that balance equals "+balanceAfterTransferUser1Expected)
	})
}

// TestIndustrialBalanceLockAndGetLocked - Checking that industrial balance can be locked
func TestIndustrialBalanceLockAndGetLocked(t *testing.T) {
	ledgerMock := mock.NewLedger(t)
	owner := ledgerMock.NewWallet()

	balanceAddAmount := "1000"
	lockAmount := "500"

	tt := &TestToken{
		token.BaseToken{
			Name:     testTokenWithGroup,
			Symbol:   testTokenSymbol,
			Decimals: 8,
		},
	}

	user1 := ledgerMock.NewWallet()

	ledgerMock.NewChainCode(testTokenWithGroup, tt, &core.ContractOptions{}, owner.Address())
	owner.SignedInvoke(testTokenWithGroup, "industrialBalanceAdd", testTokenWithGroup, user1.Address(), balanceAddAmount, "add industrial balance for "+balanceAddAmount)

	balanceResponse := owner.Invoke(testTokenWithGroup, "industrialBalanceGet", user1.Address())
	balance, err := GetIndustrialBalanceFromResponseByGroup(balanceResponse, testGroup)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	assert.Equal(t, balanceAddAmount, balance, "assert that balance equals "+balanceAddAmount)

	t.Run("Industrial balance lock and get test", func(t *testing.T) {
		owner.SignedInvoke(testTokenWithGroup, "industrialBalanceLock", testTokenWithGroup, user1.Address(), lockAmount)

		balanceAfterLockResponse := owner.Invoke(testTokenWithGroup, "industrialBalanceGet", user1.Address())
		balanceAfterLock, err := GetIndustrialBalanceFromResponseByGroup(balanceAfterLockResponse, testGroup)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, lockAmount, balanceAfterLock, "assert that balance equals "+balanceAfterLock)

		lockedBalanceResponse := owner.Invoke(testTokenWithGroup, "industrialBalanceGetLocked", user1.Address())
		lockedBalance, err := GetIndustrialBalanceFromResponseByGroup(lockedBalanceResponse, testGroup)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, lockAmount, lockedBalance, "assert that locked balance equals "+lockedBalance)
	})
}

// TestIndustrialBalanceUnLock - Checking that industrial balance can be unlocked
func TestIndustrialBalanceUnLock(t *testing.T) {
	ledgerMock := mock.NewLedger(t)
	owner := ledgerMock.NewWallet()

	balanceAddAmount := "1000"
	lockAmount := "500"

	tt := &TestToken{
		token.BaseToken{
			Name:     testTokenWithGroup,
			Symbol:   testTokenSymbol,
			Decimals: 8,
		},
	}

	user1 := ledgerMock.NewWallet()

	ledgerMock.NewChainCode(testTokenWithGroup, tt, &core.ContractOptions{}, owner.Address())
	owner.SignedInvoke(testTokenWithGroup, "industrialBalanceAdd", testTokenWithGroup, user1.Address(), balanceAddAmount, "add industrial balance for "+balanceAddAmount)

	balanceResponse := owner.Invoke(testTokenWithGroup, "industrialBalanceGet", user1.Address())
	balance, err := GetIndustrialBalanceFromResponseByGroup(balanceResponse, testGroup)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	assert.Equal(t, balanceAddAmount, balance, "assert that balance equals "+balanceAddAmount)

	owner.SignedInvoke(testTokenWithGroup, "industrialBalanceLock", testTokenWithGroup, user1.Address(), lockAmount)

	balanceAfterLockResponse := owner.Invoke(testTokenWithGroup, "industrialBalanceGet", user1.Address())
	balanceAfterLock, err := GetIndustrialBalanceFromResponseByGroup(balanceAfterLockResponse, testGroup)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	assert.Equal(t, lockAmount, balanceAfterLock, "assert that balance equals "+balanceAfterLock)

	lockedBalanceResponse := owner.Invoke(testTokenWithGroup, "industrialBalanceGetLocked", user1.Address())
	lockedBalance, err := GetIndustrialBalanceFromResponseByGroup(lockedBalanceResponse, testGroup)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	assert.Equal(t, lockAmount, lockedBalance, "assert that locked balance equals "+lockedBalance)

	t.Run("Industrial balance unlock test", func(t *testing.T) {
		owner.SignedInvoke(testTokenWithGroup, "industrialBalanceUnLock", testTokenWithGroup, user1.Address(), "300")
		lockedBalanceAfterUnlockResponse := owner.Invoke(testTokenWithGroup, "industrialBalanceGetLocked", user1.Address())
		lockedBalanceAfterUnlock, err := GetIndustrialBalanceFromResponseByGroup(lockedBalanceAfterUnlockResponse, testGroup)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "200", lockedBalanceAfterUnlock, "assert that locked balance equals "+lockedBalance)
	})
}

// TestIndustrialBalanceTransferLocked - Checking that locked industrial balance can be transfrred
func TestIndustrialBalanceTransferLocked(t *testing.T) {
	ledgerMock := mock.NewLedger(t)
	owner := ledgerMock.NewWallet()

	balanceAddAmount := "1000"
	lockAmount := "500"
	transferAmount := "300"

	tt := &TestToken{
		token.BaseToken{
			Name:     testTokenWithGroup,
			Symbol:   testTokenSymbol,
			Decimals: 8,
		},
	}

	user1 := ledgerMock.NewWallet()
	user2 := ledgerMock.NewWallet()

	ledgerMock.NewChainCode(testTokenWithGroup, tt, &core.ContractOptions{}, owner.Address())
	owner.SignedInvoke(testTokenWithGroup, "industrialBalanceAdd", testTokenWithGroup, user1.Address(), balanceAddAmount, "add industrial balance for "+balanceAddAmount)

	balanceResponse := owner.Invoke(testTokenWithGroup, "industrialBalanceGet", user1.Address())
	balance, err := GetIndustrialBalanceFromResponseByGroup(balanceResponse, testGroup)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	assert.Equal(t, balanceAddAmount, balance, "assert that balance equals "+balanceAddAmount)
	owner.SignedInvoke(testTokenWithGroup, "industrialBalanceLock", testTokenWithGroup, user1.Address(), lockAmount)

	balanceAfterLockResponse := owner.Invoke(testTokenWithGroup, "industrialBalanceGet", user1.Address())
	balanceAfterLock, err := GetIndustrialBalanceFromResponseByGroup(balanceAfterLockResponse, testGroup)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	assert.Equal(t, lockAmount, balanceAfterLock, "assert that balance equals "+balanceAfterLock)

	lockedBalanceResponse := owner.Invoke(testTokenWithGroup, "industrialBalanceGetLocked", user1.Address())
	lockedBalance, err := GetIndustrialBalanceFromResponseByGroup(lockedBalanceResponse, testGroup)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	assert.Equal(t, lockAmount, lockedBalance, "assert that locked balance equals "+lockedBalance)

	t.Run("Industrial balance transfer locked", func(t *testing.T) {
		owner.SignedInvoke(testTokenWithGroup, "industrialBalanceTransferLocked", testTokenWithGroup, user1.Address(), user2.Address(), transferAmount, "transfer locked")

		balanceAfterTransferUser1Response := owner.Invoke(testTokenWithGroup, "industrialBalanceGetLocked", user1.Address())
		balanceAfterTransferUser1, err := GetIndustrialBalanceFromResponseByGroup(balanceAfterTransferUser1Response, testGroup)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "200", balanceAfterTransferUser1, "assert that locked balance equals "+lockedBalance)

		balanceAfterTransferUser2Response := owner.Invoke(testTokenWithGroup, "industrialBalanceGet", user2.Address())
		balanceAfterTransferUser2, err := GetIndustrialBalanceFromResponseByGroup(balanceAfterTransferUser2Response, testGroup)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, transferAmount, balanceAfterTransferUser2, "assert that locked balance equals "+lockedBalance)
	})
}

// TestIndustrialBalanceBurnLocked - Checking that locked industrial balance can be burned
func TestIndustrialBalanceBurnLocked(t *testing.T) {
	ledgerMock := mock.NewLedger(t)
	owner := ledgerMock.NewWallet()

	balanceAddAmount := "1000"
	lockAmount := "500"

	tt := &TestToken{
		token.BaseToken{
			Name:     testTokenWithGroup,
			Symbol:   testTokenSymbol,
			Decimals: 8,
		},
	}

	user1 := ledgerMock.NewWallet()

	ledgerMock.NewChainCode(testTokenWithGroup, tt, &core.ContractOptions{}, owner.Address())
	owner.SignedInvoke(testTokenWithGroup, "industrialBalanceAdd", testTokenWithGroup, user1.Address(), balanceAddAmount, "add industrial balance for "+balanceAddAmount)

	balanceResponse := owner.Invoke(testTokenWithGroup, "industrialBalanceGet", user1.Address())
	balance, err := GetIndustrialBalanceFromResponseByGroup(balanceResponse, testGroup)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	assert.Equal(t, balanceAddAmount, balance, "assert that balance equals "+balanceAddAmount)

	owner.SignedInvoke(testTokenWithGroup, "industrialBalanceLock", testTokenWithGroup, user1.Address(), lockAmount)

	balanceAfterLockResponse := owner.Invoke(testTokenWithGroup, "industrialBalanceGet", user1.Address())
	balanceAfterLock, err := GetIndustrialBalanceFromResponseByGroup(balanceAfterLockResponse, testGroup)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	assert.Equal(t, lockAmount, balanceAfterLock, "assert that balance equals "+balanceAfterLock)

	lockedBalanceResponse := owner.Invoke(testTokenWithGroup, "industrialBalanceGetLocked", user1.Address())
	lockedBalance, err := GetIndustrialBalanceFromResponseByGroup(lockedBalanceResponse, testGroup)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	assert.Equal(t, lockAmount, lockedBalance, "assert that locked balance equals "+lockedBalance)

	t.Run("Industrial balance burn locked", func(t *testing.T) {
		owner.SignedInvoke(testTokenWithGroup, "industrialBalanceBurnLocked", testTokenWithGroup, user1.Address(), "300", "burn locked")

		lockedBalanceAfterBurnResponse := owner.Invoke(testTokenWithGroup, "industrialBalanceGetLocked", user1.Address())
		lockedBalanceAfterBurn, err := GetIndustrialBalanceFromResponseByGroup(lockedBalanceAfterBurnResponse, testGroup)
		if err != nil {
			assert.FailNow(t, err.Error())
		}
		assert.Equal(t, "200", lockedBalanceAfterBurn, "assert that locked balance equals "+lockedBalance)
	})
}

func GetIndustrialBalanceFromResponseByGroup(response string, group string) (string, error) {
	var balanceMap map[string]string
	err := json.Unmarshal([]byte(response), &balanceMap)
	if err != nil {
		return "", err
	}
	bl := balanceMap[group]
	if bl == "" {
		return "", errors.New("cant find balance for group: " + group)
	}
	return bl, nil
}
