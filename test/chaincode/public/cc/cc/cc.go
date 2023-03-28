package main

import (
	"log"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/tickets-dao/foundation/v3/core"
	"github.com/tickets-dao/foundation/v3/token"
)

type CcToken struct {
	token.BaseToken
}

func NewCcToken(bt token.BaseToken) *CcToken {
	return &CcToken{bt}
}

func main() {
	cct := NewCcToken(token.BaseToken{
		Name:            "Currency Coin",
		Symbol:          "CC",
		Decimals:        8,
		UnderlyingAsset: "US Dollars",
	})
	cc, err := core.NewChainCode(cct, "org0", nil)
	if err != nil {
		log.Fatal(err)
	}
	if err = shim.Start(cc); err != nil {
		log.Fatal(err)
	}
}
