package token

import (
	"errors"

	"github.com/tickets-dao/foundation/v3/core/types"
	"github.com/tickets-dao/foundation/v3/core/types/big"
)

func (bt *BaseToken) CheckLimitsAndPrice(method string, amount *big.Int, currency string) (*big.Int, error) {
	rate, exists, err := bt.GetRateAndLimits(method, currency)
	if err != nil {
		return big.NewInt(0), err
	}
	if !exists {
		return big.NewInt(0), errors.New("impossible to buy for this currency")
	}
	if !rate.InLimit(amount) {
		return big.NewInt(0), errors.New("amount out of limits")
	}
	return rate.CalcPrice(amount, RateDecimal), nil
}

func (bt *BaseToken) TxBuyToken(sender *types.Sender, amount *big.Int, currency string) error {
	if sender.Equal(bt.Issuer()) {
		return errors.New("impossible operation")
	}

	if amount.Cmp(big.NewInt(0)) == 0 {
		return errors.New("amount should be more than zero")
	}

	price, err := bt.CheckLimitsAndPrice("buyToken", amount, currency)
	if err != nil {
		return err
	}

	if err = bt.AllowedBalanceTransfer(currency, sender.Address(), bt.Issuer(), price, "buyToken"); err != nil {
		return err
	}

	if err = bt.TokenBalanceTransfer(bt.Issuer(), sender.Address(), amount, "buyToken"); err != nil {
		return err
	}
	return nil
}

func (bt *BaseToken) TxBuyBack(sender *types.Sender, amount *big.Int, currency string) error {
	if sender.Equal(bt.Issuer()) {
		return errors.New("impossible operation")
	}

	if amount.Cmp(big.NewInt(0)) == 0 {
		return errors.New("amount should be more than zero")
	}

	price, err := bt.CheckLimitsAndPrice("buyBack", amount, currency)
	if err != nil {
		return err
	}

	if err = bt.AllowedBalanceTransfer(currency, bt.Issuer(), sender.Address(), price, "buyBack"); err != nil {
		return err
	}

	if err = bt.TokenBalanceTransfer(sender.Address(), bt.Issuer(), amount, "buyBack"); err != nil {
		return err
	}
	return nil
}
