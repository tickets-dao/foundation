package unit

import (
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tickets-dao/foundation/v3/core"
	"github.com/tickets-dao/foundation/v3/core/types"
	"github.com/tickets-dao/foundation/v3/core/types/big"
	"github.com/tickets-dao/foundation/v3/mock"
	"github.com/tickets-dao/foundation/v3/token"
	"golang.org/x/crypto/sha3"
)

// implemented through query requests, an error comes through tx
// relate github.com/tickets-dao/foundation/-/issues/44 https://github.com/tickets-dao/foundation/-/issues/45

func (tt *TestToken) QueryAllowedBalanceAdd(token string, address *types.Address, amount *big.Int, reason string) error {
	return tt.AllowedBalanceAdd(token, address, amount, reason)
}

func (tt *TestToken) QueryAllowedBalanceSub(token string, address *types.Address, amount *big.Int, reason string) error {
	return tt.AllowedBalanceSub(token, address, amount, reason)
}

func (tt *TestToken) QueryAllowedBalanceLock(token string, address *types.Address, amount *big.Int) error {
	return tt.AllowedBalanceLock(token, address, amount)
}

func (tt *TestToken) QueryAllowedBalanceUnLock(token string, address *types.Address, amount *big.Int) error {
	return tt.AllowedBalanceUnLock(token, address, amount)
}

func (tt *TestToken) QueryAllowedBalanceTransferLocked(token string, from *types.Address, to *types.Address, amount *big.Int, reason string) error {
	return tt.AllowedBalanceTransferLocked(token, from, to, amount, reason)
}

func (tt *TestToken) QueryAllowedBalanceBurnLocked(token string, address *types.Address, amount *big.Int, reason string) error {
	return tt.AllowedBalanceBurnLocked(token, address, amount, reason)
}

func (tt *TestToken) QueryAllowedBalanceGetAll(address *types.Address) (map[string]string, error) {
	return tt.AllowedBalanceGetAll(address)
}

// TestAllowedIndustrialBalanceTransfer - Checking that allowed industrial balance can be transferred
func TestAllowedIndustrialBalanceTransfer(t *testing.T) {
	ledgerMock := mock.NewLedger(t)
	owner := ledgerMock.NewWallet()
	feeAddressSetter := ledgerMock.NewWallet()
	feeSetter := ledgerMock.NewWallet()
	user := ledgerMock.NewWallet()

	vt := token.VT{
		BaseToken: token.BaseToken{
			Name:     "Validation Token",
			Symbol:   "VT",
			Decimals: 8,
		},
	}

	ledgerMock.NewChainCode("vt", &vt, &core.ContractOptions{}, owner.Address(), feeSetter.Address(), feeAddressSetter.Address())

	const (
		ba1 = "BA02_GOLDBARLONDON.01"
		ba2 = "BA02_GOLDBARLONDON.02"
	)

	owner.AddAllowedBalance("vt", ba1, 100000000)
	owner.AddAllowedBalance("vt", ba2, 100000000)
	owner.AllowedBalanceShouldBe("vt", ba1, 100000000)
	owner.AllowedBalanceShouldBe("vt", ba2, 100000000)

	industrialAssets := []*types.MultiSwapAsset{
		{
			Group:  ba1,
			Amount: "50000000",
		},
		{
			Group:  ba2,
			Amount: "100000000",
		},
	}

	rawGA, err := json.Marshal(industrialAssets)
	assert.NoError(t, err)

	owner.SignedInvoke("vt", "allowedIndustrialBalanceTransfer", user.Address(), string(rawGA), "ref")
	owner.AllowedBalanceShouldBe("vt", ba1, 50000000)
	owner.AllowedBalanceShouldBe("vt", ba2, 0)
	user.AllowedBalanceShouldBe("vt", ba1, 50000000)
	user.AllowedBalanceShouldBe("vt", ba2, 100000000)
}

// TestAllowedBalanceAdd - Checking that allowed balance can be added
func TestAllowedBalanceAdd(t *testing.T) { //nolint:dupl
	ledgerMock := mock.NewLedger(t)
	owner := ledgerMock.NewWallet()
	cc := &TestToken{
		token.BaseToken{
			Symbol: "CC",
		},
	}
	ledgerMock.NewChainCode("cc", cc, nil, owner.Address())

	vt := &TestToken{
		token.BaseToken{
			Symbol: "VT",
		},
	}
	ledgerMock.NewChainCode("vt", vt, nil, owner.Address())

	user1 := ledgerMock.NewWallet()
	user1.AddBalance("cc", 1000)

	swapKey := "123"
	hashed := sha3.Sum256([]byte(swapKey))
	swapHash := hex.EncodeToString(hashed[:])

	txID := user1.SignedInvoke("cc", "swapBegin", "CC", "VT", "450", swapHash)
	user1.BalanceShouldBe("cc", 550)
	ledgerMock.WaitSwapAnswer("vt", txID, time.Second*5)

	user1.Invoke("vt", "swapDone", txID, swapKey)
	user1.AllowedBalanceShouldBe("vt", "CC", 450)

	t.Run("Allowed balance add", func(t *testing.T) {
		owner.Invoke("vt", "allowedBalanceAdd", "CC", user1.Address(), "50", "add balance")
		balance := owner.Invoke("vt", "allowedBalanceOf", user1.Address(), "CC")
		assert.Equal(t, "\"500\"", balance)
	})
}

// TestAllowedBalanceSub - Checking that balance sub is working
func TestAllowedBalanceSub(t *testing.T) { //nolint:dupl
	ledgerMock := mock.NewLedger(t)
	owner := ledgerMock.NewWallet()
	cc := &TestToken{
		token.BaseToken{
			Symbol: "CC",
		},
	}
	ledgerMock.NewChainCode("cc", cc, nil, owner.Address())

	vt := &TestToken{
		token.BaseToken{
			Symbol: "VT",
		},
	}
	ledgerMock.NewChainCode("vt", vt, nil, owner.Address())

	user1 := ledgerMock.NewWallet()
	user1.AddBalance("cc", 1000)

	swapKey := "123"
	hashed := sha3.Sum256([]byte(swapKey))
	swapHash := hex.EncodeToString(hashed[:])

	txID := user1.SignedInvoke("cc", "swapBegin", "CC", "VT", "450", swapHash)
	user1.BalanceShouldBe("cc", 550)
	ledgerMock.WaitSwapAnswer("vt", txID, time.Second*5)

	user1.Invoke("vt", "swapDone", txID, swapKey)
	user1.AllowedBalanceShouldBe("vt", "CC", 450)

	t.Run("Allowed balance sub", func(t *testing.T) {
		owner.Invoke("vt", "allowedBalanceSub", "CC", user1.Address(), "50", "sub balance")
		balance := owner.Invoke("vt", "allowedBalanceOf", user1.Address(), "CC")
		assert.Equal(t, "\"400\"", balance)
	})
}

// TestAllowedBalanceLock - Checking that allowed balance can be locked
func TestAllowedBalanceLock(t *testing.T) {
	ledgerMock := mock.NewLedger(t)
	owner := ledgerMock.NewWallet()
	cc := &TestToken{
		token.BaseToken{
			Symbol: "CC",
		},
	}
	ledgerMock.NewChainCode("cc", cc, nil, owner.Address())

	vt := &TestToken{
		token.BaseToken{
			Symbol: "VT",
		},
	}
	ledgerMock.NewChainCode("vt", vt, nil, owner.Address())

	user1 := ledgerMock.NewWallet()
	user1.AddBalance("cc", 1000)

	swapKey := "123"
	hashed := sha3.Sum256([]byte(swapKey))
	swapHash := hex.EncodeToString(hashed[:])

	txID := user1.SignedInvoke("cc", "swapBegin", "CC", "VT", "450", swapHash)
	user1.BalanceShouldBe("cc", 550)
	ledgerMock.WaitSwapAnswer("vt", txID, time.Second*5)

	user1.Invoke("vt", "swapDone", txID, swapKey)
	user1.AllowedBalanceShouldBe("vt", "CC", 450)

	t.Run("Allowed balance lock", func(t *testing.T) {
		owner.Invoke("vt", "allowedBalanceLock", "CC", user1.Address(), "50")
		balance := owner.Invoke("vt", "allowedBalanceOf", user1.Address(), "CC")
		assert.Equal(t, "\"400\"", balance)
	})
}

// TestAllowedBalanceUnLock - Checking that allowed balance can be unlocked
func TestAllowedBalanceUnLock(t *testing.T) {
	ledgerMock := mock.NewLedger(t)
	owner := ledgerMock.NewWallet()
	cc := &TestToken{
		token.BaseToken{
			Symbol: "CC",
		},
	}
	ledgerMock.NewChainCode("cc", cc, nil, owner.Address())

	vt := &TestToken{
		token.BaseToken{
			Symbol: "VT",
		},
	}
	ledgerMock.NewChainCode("vt", vt, nil, owner.Address())

	user1 := ledgerMock.NewWallet()
	user1.AddBalance("cc", 1000)

	swapKey := "123"
	hashed := sha3.Sum256([]byte(swapKey))
	swapHash := hex.EncodeToString(hashed[:])

	txID := user1.SignedInvoke("cc", "swapBegin", "CC", "VT", "450", swapHash)
	user1.BalanceShouldBe("cc", 550)
	ledgerMock.WaitSwapAnswer("vt", txID, time.Second*5)

	user1.Invoke("vt", "swapDone", txID, swapKey)
	user1.AllowedBalanceShouldBe("vt", "CC", 450)
	owner.Invoke("vt", "allowedBalanceLock", "CC", user1.Address(), "50")
	balance := owner.Invoke("vt", "allowedBalanceOf", user1.Address(), "CC")
	assert.Equal(t, "\"400\"", balance)

	t.Run("Allowed balance unlock", func(t *testing.T) {
		owner.Invoke("vt", "allowedBalanceUnLock", "CC", user1.Address(), "50")
		balance := owner.Invoke("vt", "allowedBalanceOf", user1.Address(), "CC")
		assert.Equal(t, "\"450\"", balance)
	})
}

// TestAllowedBalanceTransferLocked - Checking that allowed locked balance can be transferred
func TestAllowedBalanceTransferLocked(t *testing.T) {
	ledgerMock := mock.NewLedger(t)
	owner := ledgerMock.NewWallet()
	cc := &TestToken{
		token.BaseToken{
			Symbol: "CC",
		},
	}
	ledgerMock.NewChainCode("cc", cc, nil, owner.Address())

	vt := &TestToken{
		token.BaseToken{
			Symbol: "VT",
		},
	}
	ledgerMock.NewChainCode("vt", vt, nil, owner.Address())

	user1 := ledgerMock.NewWallet()
	user2 := ledgerMock.NewWallet()

	user1.AddBalance("cc", 1000)

	swapKey := "123"
	hashed := sha3.Sum256([]byte(swapKey))
	swapHash := hex.EncodeToString(hashed[:])

	txID := user1.SignedInvoke("cc", "swapBegin", "CC", "VT", "450", swapHash)
	user1.BalanceShouldBe("cc", 550)
	ledgerMock.WaitSwapAnswer("vt", txID, time.Second*5)

	user1.Invoke("vt", "swapDone", txID, swapKey)
	user1.AllowedBalanceShouldBe("vt", "CC", 450)
	owner.Invoke("vt", "allowedBalanceLock", "CC", user1.Address(), "50")
	balance := owner.Invoke("vt", "allowedBalanceOf", user1.Address(), "CC")
	assert.Equal(t, "\"400\"", balance)

	t.Run("Allowed balance transfer locked", func(t *testing.T) {
		owner.Invoke("vt", "allowedBalanceTransferLocked", "CC", user1.Address(), user2.Address(), "50", "transfer")
		user2.AllowedBalanceShouldBe("vt", "CC", 50)
	})
}

// TestAllowedBalanceBurnLocked - Checking that allowed balance can be burned
func TestAllowedBalanceBurnLocked(t *testing.T) {
	ledgerMock := mock.NewLedger(t)
	owner := ledgerMock.NewWallet()
	cc := &TestToken{
		token.BaseToken{
			Symbol: "CC",
		},
	}
	ledgerMock.NewChainCode("cc", cc, nil, owner.Address())

	vt := &TestToken{
		token.BaseToken{
			Symbol: "VT",
		},
	}
	ledgerMock.NewChainCode("vt", vt, nil, owner.Address())

	user1 := ledgerMock.NewWallet()

	user1.AddBalance("cc", 1000)

	swapKey := "123"
	hashed := sha3.Sum256([]byte(swapKey))
	swapHash := hex.EncodeToString(hashed[:])

	txID := user1.SignedInvoke("cc", "swapBegin", "CC", "VT", "450", swapHash)
	user1.BalanceShouldBe("cc", 550)
	ledgerMock.WaitSwapAnswer("vt", txID, time.Second*5)

	user1.Invoke("vt", "swapDone", txID, swapKey)
	user1.AllowedBalanceShouldBe("vt", "CC", 450)
	owner.Invoke("vt", "allowedBalanceLock", "CC", user1.Address(), "50")
	balance := owner.Invoke("vt", "allowedBalanceOf", user1.Address(), "CC")
	assert.Equal(t, "\"400\"", balance)

	t.Run("Allowed balance burn locked", func(t *testing.T) {
		owner.Invoke("vt", "allowedBalanceBurnLocked", "CC", user1.Address(), "50", "transfer")
	})
}

// TestAllowedBalancesGetAll - Checking that all allowed balances can be geted
func TestAllowedBalancesGetAll(t *testing.T) {
	ledgerMock := mock.NewLedger(t)
	owner := ledgerMock.NewWallet()
	cc := &TestToken{
		token.BaseToken{
			Symbol: "CC",
		},
	}
	ledgerMock.NewChainCode("cc", cc, nil, owner.Address())

	vt := &TestToken{
		token.BaseToken{
			Symbol: "VT",
		},
	}
	ledgerMock.NewChainCode("vt", vt, nil, owner.Address())

	nt := &TestToken{
		token.BaseToken{
			Symbol: "NT",
		},
	}
	ledgerMock.NewChainCode("nt", nt, nil, owner.Address())

	user1 := ledgerMock.NewWallet()
	user1.AddBalance("cc", 1000)

	swapKey := "123"
	hashed := sha3.Sum256([]byte(swapKey))
	swapHash := hex.EncodeToString(hashed[:])

	txID := user1.SignedInvoke("cc", "swapBegin", "CC", "VT", "450", swapHash)
	user1.BalanceShouldBe("cc", 550)
	ledgerMock.WaitSwapAnswer("vt", txID, time.Second*5)
	user1.Invoke("vt", "swapDone", txID, swapKey)

	txID2 := user1.SignedInvoke("cc", "swapBegin", "CC", "VT", "150", swapHash)
	user1.BalanceShouldBe("cc", 400)
	ledgerMock.WaitSwapAnswer("vt", txID2, time.Second*5)
	user1.Invoke("vt", "swapDone", txID2, swapKey)

	t.Run("Allowed balances get all", func(t *testing.T) {
		balance := owner.Invoke("vt", "allowedBalanceGetAll", user1.Address())
		assert.Equal(t, "{\"CC\":\"600\"}", balance)
	})
}
