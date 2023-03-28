package token

import (
	"errors"
	"fmt"

	"github.com/tickets-dao/foundation/v3/core"
	"github.com/tickets-dao/foundation/v3/core/types"
	"github.com/tickets-dao/foundation/v3/core/types/big"
	"github.com/tickets-dao/foundation/v3/proto"
)

type Metadata struct {
	Name            string          `json:"name"`
	Symbol          string          `json:"symbol"`
	Decimals        uint            `json:"decimals"`
	UnderlyingAsset string          `json:"underlying_asset"` //nolint:tagliatelle
	Issuer          string          `json:"issuer"`
	Methods         []string        `json:"methods"`
	TotalEmission   *big.Int        `json:"total_emission"` //nolint:tagliatelle
	Fee             *Fee            `json:"fee"`
	Rates           []*MetadataRate `json:"rates"`
}

type MetadataRate struct {
	DealType string   `json:"deal_type"` //nolint:tagliatelle
	Currency string   `json:"currency"`
	Rate     *big.Int `json:"rate"`
	Min      *big.Int `json:"min"`
	Max      *big.Int `json:"max"`
}

type Fee struct {
	Address  string   `json:"address"`
	Currency string   `json:"currency"`
	Fee      *big.Int `json:"fee"`
	Floor    *big.Int `json:"floor"`
	Cap      *big.Int `json:"cap"`
}

// QueryMetadata returns Metadata
func (bt *BaseToken) QueryMetadata() (*Metadata, error) {
	if err := bt.loadConfigUnlessLoaded(); err != nil {
		return &Metadata{}, err
	}
	m := &Metadata{
		Name:            bt.Name,
		Symbol:          bt.Symbol,
		Decimals:        bt.Decimals,
		UnderlyingAsset: bt.UnderlyingAsset,
		Issuer:          bt.Issuer().String(),
		Methods:         bt.GetMethods(),
		TotalEmission:   new(big.Int).SetBytes(bt.config.TotalEmission),
		Fee:             &Fee{},
	}
	if types.IsValidAddressLen(bt.config.FeeAddress) {
		m.Fee.Address = types.AddrFromBytes(bt.config.FeeAddress).String()
	}
	if bt.config.Fee != nil {
		m.Fee.Currency = bt.config.Fee.Currency
		m.Fee.Fee = new(big.Int).SetBytes(bt.config.Fee.Fee)
		m.Fee.Floor = new(big.Int).SetBytes(bt.config.Fee.Floor)
		m.Fee.Cap = new(big.Int).SetBytes(bt.config.Fee.Cap)
	}
	for _, r := range bt.config.Rates {
		m.Rates = append(m.Rates, &MetadataRate{
			DealType: r.DealType,
			Currency: r.Currency,
			Rate:     new(big.Int).SetBytes(r.Rate),
			Min:      new(big.Int).SetBytes(r.Min),
			Max:      new(big.Int).SetBytes(r.Max),
		})
	}
	return m, nil
}

// QueryBalanceOf returns balance
func (bt *BaseToken) QueryBalanceOf(address *types.Address) (*big.Int, error) {
	return bt.TokenBalanceGet(address)
}

// QueryAllowedBalanceOf returns allowed balance
func (bt *BaseToken) QueryAllowedBalanceOf(address *types.Address, token string) (*big.Int, error) {
	return bt.AllowedBalanceGet(token, address)
}

// QueryDocumentsList - returns list of emission documents
func (bt *BaseToken) QueryDocumentsList() ([]core.Doc, error) {
	return core.DocumentsList(bt.GetStub())
}

// TxAddDocs - adds docs to a token
func (bt *BaseToken) TxAddDocs(sender *types.Sender, rawDocs string) error {
	if !sender.Equal(bt.Issuer()) {
		return errors.New("unathorized")
	}

	return core.AddDocs(bt.GetStub(), rawDocs)
}

// TxDeleteDoc - deletes doc from state
func (bt *BaseToken) TxDeleteDoc(sender *types.Sender, docID string) error {
	if !sender.Equal(bt.Issuer()) {
		return errors.New("unathorized")
	}

	return core.DeleteDoc(bt.GetStub(), docID)
}

// TxSetRate sets token rate to an asset for a type of deal
func (bt *BaseToken) TxSetRate(sender *types.Sender, dealType string, currency string, rate *big.Int) error {
	if !sender.Equal(bt.Issuer()) {
		return errors.New("unauthorized")
	}
	// TODO - check if it may be helpful in business logic
	if rate.Sign() == 0 {
		return errors.New("trying to set rate = 0")
	}
	if bt.Symbol == currency {
		return errors.New("currency is equals token: it is impossible")
	}
	if err := bt.loadConfigUnlessLoaded(); err != nil {
		return err
	}
	for i, r := range bt.config.Rates {
		if r.DealType == dealType && r.Currency == currency {
			bt.config.Rates[i].Rate = rate.Bytes()
			return bt.saveConfig()
		}
	}
	bt.config.Rates = append(bt.config.Rates, &proto.TokenRate{
		DealType: dealType,
		Currency: currency,
		Rate:     rate.Bytes(),
		Max:      new(big.Int).SetUint64(0).Bytes(), // todo maybe needs different solution
		Min:      new(big.Int).SetUint64(0).Bytes(),
	})
	return bt.saveConfig()
}

// TxSetLimits sets limits for a deal type and an asset
func (bt *BaseToken) TxSetLimits(sender *types.Sender, dealType string, currency string, min *big.Int, max *big.Int) error {
	if !sender.Equal(bt.Issuer()) {
		return errors.New("unauthorized")
	}
	if min.Cmp(max) > 0 && max.Cmp(big.NewInt(0)) > 0 {
		return errors.New("min limit is greater than max limit")
	}
	if err := bt.loadConfigUnlessLoaded(); err != nil {
		return err
	}
	unknownDealType := true
	for i, r := range bt.config.Rates {
		if r.DealType == dealType {
			unknownDealType = false
			if r.Currency == currency {
				bt.config.Rates[i].Max = max.Bytes()
				bt.config.Rates[i].Min = min.Bytes()
				return bt.saveConfig()
			}
		}
	}
	if unknownDealType {
		return fmt.Errorf("unknown DealType. Rate for deal type %s and currency %s was not set", dealType, currency)
	}
	return fmt.Errorf("unknown currency. Rate for deal type %s and currency %s was not set", dealType, currency)
}

// TxDeleteRate - deletes rate from state
func (bt *BaseToken) TxDeleteRate(sender *types.Sender, dealType string, currency string) error {
	if !sender.Equal(bt.Issuer()) {
		return errors.New("unauthorized")
	}
	if bt.Symbol == currency {
		return errors.New("currency is equals token: it is impossible")
	}
	if err := bt.loadConfigUnlessLoaded(); err != nil {
		return err
	}
	for i, r := range bt.config.Rates {
		if r.DealType == dealType && r.Currency == currency {
			bt.config.Rates = append(bt.config.Rates[:i], bt.config.Rates[i+1:]...)
			return bt.saveConfig()
		}
	}

	return nil
}
