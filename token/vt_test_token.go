package token

import (
	"errors"

	"github.com/tickets-dao/foundation/v3/core/types"
	"github.com/tickets-dao/foundation/v3/core/types/big"
)

type VT struct {
	BaseToken
}

func (vt *VT) TxEmitToken(sender *types.Sender, amount *big.Int) error {
	if !sender.Equal(vt.Issuer()) {
		return errors.New("unauthorized")
	}
	if err := vt.TokenBalanceAdd(vt.Issuer(), amount, "emitToken"); err != nil {
		return err
	}
	return vt.EmissionAdd(amount)
}
