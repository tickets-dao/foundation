package main

import (
	"errors"
	"log"

	"github.com/tickets-dao/foundation/v3/core"
	"github.com/tickets-dao/foundation/v3/core/types"
	industrialtoken "github.com/tickets-dao/integration/chaincode/industrial/industrial_token"

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

// IT - industrial token base struct
type IT struct {
	industrialtoken.IndustrialToken
}

var groups = []industrialtoken.Group{{
	ID:       "202009",
	Emission: 10000000000000,
	Maturity: "21.09.2020 22:00:00",
	Note:     "Test note",
}, {
	ID:       "202010",
	Emission: 100000000000000,
	Maturity: "21.10.2020 22:00:00",
	Note:     "Test note",
}, {
	ID:       "202011",
	Emission: 200000000000000,
	Maturity: "21.11.2020 22:00:00",
	Note:     "Test note",
}, {
	ID:       "202012",
	Emission: 50000000000000,
	Maturity: "21.12.2020 22:00:00",
	Note:     "Test note",
},
}

// NBTxInitialize - initializes chaincode
func (it *IT) NBTxInitialize(sender *types.Sender) error {
	if !sender.Equal(it.Issuer()) {
		return errors.New("unauthorized")
	}

	return it.Initialize(groups)
}

func main() {
	ft := &IT{
		industrialtoken.IndustrialToken{
			Name:            "Test Industrial Token",
			Symbol:          "INDUSTRIAL",
			Decimals:        8,
			UnderlyingAsset: "TEST_UnderlyingAsset",
			DeliveryForm:    "TEST_DeliveryForm",
			UnitOfMeasure:   "TEST_IT",
			TokensForUnit:   "1",
			PaymentTerms:    "Non-prepaid",
			Price:           "Floating",
		},
	}

	cc, err := core.NewChainCode(ft, "org0", nil)
	if err != nil {
		log.Fatal(err)
	}
	if err = shim.Start(cc); err != nil {
		log.Fatal(err)
	}
}
