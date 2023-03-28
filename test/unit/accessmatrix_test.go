package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tickets-dao/foundation/v3/core"
	"github.com/tickets-dao/foundation/v3/core/acl"
	"github.com/tickets-dao/foundation/v3/core/types"
	"github.com/tickets-dao/foundation/v3/mock"
	"github.com/tickets-dao/foundation/v3/token"
)

type IssuerCheckerToken struct {
	token.BaseToken
}

const (
	getRightFn = "getRight"
)

func (ict *IssuerCheckerToken) QueryGetRight(ccname string, address *types.Address, role, operation string) (bool, error) {
	stub := ict.GetStub()
	r := acl.Role(role)
	params := []string{stub.GetChannelID(), ccname, r.String(), operation, address.String()}
	right, err := acl.GetAccountRight(ict.GetStub(), params)
	if err != nil {
		return false, err
	}

	if right.HaveRight {
		return true, nil
	}

	return false, err
}

func TestRights(t *testing.T) {
	ledgerMock := mock.NewLedger(t)
	issuer := ledgerMock.NewWallet()
	feeSetter := ledgerMock.NewWallet()

	ict := &IssuerCheckerToken{
		token.BaseToken{
			Name:     testTokenName,
			Symbol:   testTokenSymbol,
			Decimals: 8,
		},
	}

	ledgerMock.NewChainCode(
		testTokenCCName,
		ict,
		&core.ContractOptions{},
		issuer.Address(),
		feeSetter.Address(),
	)

	const (
		createOp = "createEmissionApp"
		acceptOp = "acceptEmissionApp"
		deleteOp = "deleteEmissionApp"
	)

	user := ledgerMock.NewWallet()

	t.Run("add right & check if it is granted for user and operation", func(t *testing.T) {
		err := issuer.AddAccountRight(&mock.Right{
			Channel:   testTokenCCName,
			Chaincode: testTokenCCName,
			Role:      acl.Issuer.String(),
			Operation: createOp,
			Address:   user.Address(),
		})
		assert.NoError(t, err)

		isIssuer := issuer.Invoke(testTokenCCName, getRightFn,
			testTokenCCName, user.Address(), acl.Issuer.String(), createOp)
		assert.Equal(t, "true", isIssuer)
	})

	t.Run("multi-emission, non-permitted operation", func(t *testing.T) {
		isIssuer := issuer.Invoke(testTokenCCName, getRightFn,
			testTokenCCName, issuer.Address(), acl.Issuer.String(), deleteOp)
		assert.Equal(t, "false", isIssuer)
	})

	t.Run("remove right & check it is removed", func(t *testing.T) {
		err := issuer.RemoveAccountRight(&mock.Right{
			Channel:   testTokenCCName,
			Chaincode: testTokenCCName,
			Role:      acl.Issuer.String(),
			Operation: createOp,
			Address:   user.Address(),
		})
		assert.NoError(t, err)
		isIssuer := issuer.Invoke(testTokenCCName, getRightFn,
			testTokenCCName, user.Address(), acl.Issuer.String(), createOp)
		assert.Equal(t, "false", isIssuer)
	})

	t.Run("check double setting right", func(t *testing.T) {
		err := issuer.AddAccountRight(&mock.Right{
			Channel:   testTokenCCName,
			Chaincode: testTokenCCName,
			Role:      acl.Issuer.String(),
			Operation: acceptOp,
			Address:   user.Address(),
		})
		assert.NoError(t, err)

		err = issuer.AddAccountRight(&mock.Right{
			Channel:   testTokenCCName,
			Chaincode: testTokenCCName,
			Role:      acl.Issuer.String(),
			Operation: acceptOp,
			Address:   user.Address(),
		})
		assert.NoError(t, err)

		err = issuer.RemoveAccountRight(&mock.Right{
			Channel:   testTokenCCName,
			Chaincode: testTokenCCName,
			Role:      acl.Issuer.String(),
			Operation: acceptOp,
			Address:   user.Address(),
		})
		assert.NoError(t, err)

		isIssuer := issuer.Invoke(testTokenCCName, getRightFn,
			testTokenCCName, user.Address(), acl.Issuer.String(), acceptOp)
		assert.Equal(t, "false", isIssuer)
	})

	// TODO need check empty address & bad address string
}
