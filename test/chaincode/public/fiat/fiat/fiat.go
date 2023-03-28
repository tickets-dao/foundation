package main

import (
	"log"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/tickets-dao/foundation/v3/core"
	"github.com/tickets-dao/foundation/v3/token"
)

func main() {
	fiat := NewFiatToken(token.BaseToken{
		Name:            "FIAT",
		Symbol:          "FIAT",
		Decimals:        8,
		UnderlyingAsset: "US Dollars",
	})
	cc, err := core.NewChainCode(fiat, "org0", &core.ContractOptions{
		DisabledFunctions: []string{"TxBuyToken", "TxBuyBack"},
		TxTTL:             30,
		NonceTTL:          10, //  для тестов зададим 10 секунд
	})
	if err != nil {
		log.Fatal(err)
	}
	if err = shim.Start(cc); err != nil {
		log.Fatal(err)
	}
}
