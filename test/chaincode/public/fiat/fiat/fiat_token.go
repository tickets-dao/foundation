package main

import (
	"errors"

	"github.com/tickets-dao/foundation/v3/core/types"
	"github.com/tickets-dao/foundation/v3/core/types/big"
	"github.com/tickets-dao/foundation/v3/token"
)

// FiatToken - base struct
type FiatToken struct {
	token.BaseToken
}

// const emits = "emits"
// const lastNonce = "lastNonce" // default 0000000000

// NewFiatToken creates fiat token
func NewFiatToken(bt token.BaseToken) *FiatToken {
	return &FiatToken{bt}
}

// TxEmit - emits fiat token
func (mt *FiatToken) TxEmit(sender *types.Sender, address *types.Address, amount *big.Int) error {
	if !sender.Equal(mt.Issuer()) {
		return errors.New("unauthorized")
	}

	if amount.Cmp(big.NewInt(0)) == 0 {
		return errors.New("amount should be more than zero")
	}

	if err := mt.TokenBalanceAdd(address, amount, "txEmit"); err != nil {
		return err
	}
	return mt.EmissionAdd(amount)
}
