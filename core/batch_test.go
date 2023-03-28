package core

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric/common/util"
	"github.com/stretchr/testify/assert"
	"github.com/tickets-dao/foundation/v3/mock/stub"
	"github.com/tickets-dao/foundation/v3/proto"
	pb "google.golang.org/protobuf/proto"
)

const (
	testFunctionBatch = "testFnWithArgs"
	allowedMspID      = "allowedMspID"
)

var (
	testChaincodeName = "chaincode"

	args     = []string{"arg0", "arg1", "arg2", "arg3", "arg4"}
	testArgs = []string{"4аап@*", "hyexc566", "kiubfvr$", ";3вкпп", "ж?отов;", "!шжуця", "gfgt^"}

	creatorSKI    = []byte("creatorSKI")
	hexCreatorSKI = hex.EncodeToString(creatorSKI)
	sender        = &proto.Address{
		UserID:  "UserId",
		Address: []byte("Address"),
	}

	txID            = "TestTxID"
	txIDBytes       = []byte(txID)
	testEncodedTxID = hex.EncodeToString(txIDBytes)
)

func (*testContract) TxArgs(arg0 string, arg1 string, arg2 string, arg3 string, arg4 string, arg5 string, arg6 string) error {
	return nil
}

type serieBatches struct {
	FnName    string
	testID    string
	errorMsg  string
	timestamp *timestamp.Timestamp
}

// TestSaveAndLoadToBatchWithWrongArgs - negative test with wrong Args in saveToBatch
// Test must be failed, but it is passed
func TestSaveAndLoadToBatchWithWrongArgs(t *testing.T) {
	t.Parallel()

	ts := util.CreateUtcTimestamp()

	s := &serieBatches{
		FnName:    "args",
		testID:    testEncodedTxID,
		errorMsg:  "",
		timestamp: ts,
	}

	SaveAndLoadToBatchTest(t, s)
}

// TestSaveAndLoadToBatchWithWrongFnParameter - negative test with wrong Fn Name for Args in saveToBatch
// Test must be failed, but it is passed
func TestSaveAndLoadToBatchWithWrongFnParameter(t *testing.T) {
	t.Parallel()

	ts := util.CreateUtcTimestamp()

	s := &serieBatches{
		FnName:    "testFnWithArgs",
		testID:    testEncodedTxID,
		errorMsg:  "",
		timestamp: ts,
	}

	t.Skip("reason: https://github.com/tickets-dao/foundation/-/issues/49")
	SaveAndLoadToBatchTest(t, s)
}

// TestSaveAndLoadToBatchWithWrongID - negative test with wrong ID for loadToBatch
func TestSaveAndLoadToBatchWithWrongID(t *testing.T) {
	t.Parallel()

	ts := util.CreateUtcTimestamp()

	s := &serieBatches{
		FnName:    "args",
		testID:    "wonder",
		errorMsg:  "transaction wonder not found",
		timestamp: ts,
	}

	SaveAndLoadToBatchTest(t, s)
}

// SaveAndLoadToBatchTest - basic test to check Args in saveToBatch and loadFromBatch
func SaveAndLoadToBatchTest(t *testing.T, ser *serieBatches) {
	chainCode, errChainCode := NewChainCode(&testContract{}, allowedMspID, nil)
	assert.NoError(t, errChainCode)

	chainCode.init = &proto.InitArgs{}
	mockStub := stub.NewMockStub(testChaincodeName, chainCode)

	mockStub.TxID = testEncodedTxID
	mockStub.MockTransactionStart(testEncodedTxID)
	if ser.timestamp != nil {
		mockStub.TxTimestamp = ser.timestamp
	}

	batchTimestamp, err := mockStub.GetTxTimestamp()
	assert.NoError(t, err)

	errSave := chainCode.saveToBatch(mockStub, ser.FnName, creatorSKI, sender, testArgs, uint64(batchTimestamp.Seconds))
	assert.NoError(t, errSave)
	mockStub.MockTransactionEnd(testEncodedTxID)
	state, err := mockStub.GetState(fmt.Sprintf("\u0000batchTransactions\u0000%s\u0000", testEncodedTxID))
	assert.NotNil(t, state)
	assert.NoError(t, err)

	pending := new(proto.PendingTx)
	err = pb.Unmarshal(state, pending)
	assert.NoError(t, err)

	assert.Equal(t, pending.Args, testArgs)

	pending, err = chainCode.loadFromBatch(mockStub, ser.testID, batchTimestamp.Seconds)
	if err != nil {
		assert.Equal(t, ser.errorMsg, err.Error())
	} else {
		assert.NoError(t, err)
		assert.Equal(t, pending.Method, ser.FnName)
		assert.Equal(t, pending.CreatorSKI, creatorSKI)
		assert.Equal(t, pending.Args, testArgs)
	}
}

type serieBatcheExecute struct {
	testIDBytes       []byte
	testhexCreatorSKI string
	paramsWrongON     bool
}

// TestBatchExecutehWithRightParams - positive test for SaveBatch, LoadBatch and batchExecute
func TestBatchExecutehWithRightParams(t *testing.T) {
	t.Parallel()

	s := &serieBatcheExecute{
		testIDBytes:       txIDBytes,
		testhexCreatorSKI: hexCreatorSKI,
		paramsWrongON:     false,
	}

	BatchExecuteTest(t, s)
}

// TestBatchExecutehWithWrongParams - negative test with wrong parameters in batchExecute
// Test must be failed, but it is passed
func TestBatchExecutehWithWrongParams(t *testing.T) {
	t.Parallel()

	s := &serieBatcheExecute{
		testIDBytes:       []byte("wonder"),
		testhexCreatorSKI: "wonder",
		paramsWrongON:     true,
	}

	BatchExecuteTest(t, s)
}

// BatchExecuteTest - basic test for SaveBatch, LoadBatch and batchExecute
func BatchExecuteTest(t *testing.T, ser *serieBatcheExecute) {
	chainCode, err := NewChainCode(&testContract{}, allowedMspID, nil)
	assert.NoError(t, err)

	chainCode.init = &proto.InitArgs{}
	mockStub := stub.NewMockStub(testChaincodeName, chainCode)

	mockStub.TxID = testEncodedTxID
	mockStub.MockTransactionStart(testEncodedTxID)

	batchTimestamp, err := mockStub.GetTxTimestamp()
	assert.NoError(t, err)

	err = chainCode.saveToBatch(mockStub, testFunctionBatch, creatorSKI, nil, args, uint64(batchTimestamp.Seconds))
	assert.NoError(t, err)
	mockStub.MockTransactionEnd(testEncodedTxID)
	state, err := mockStub.GetState(fmt.Sprintf("\u0000batchTransactions\u0000%s\u0000", testEncodedTxID))
	assert.NotNil(t, state)
	assert.NoError(t, err)

	pending := new(proto.PendingTx)
	err = pb.Unmarshal(state, pending)
	assert.NoError(t, err)

	assert.Equal(t, pending.Method, testFunctionBatch)
	assert.Equal(t, pending.CreatorSKI, creatorSKI)
	assert.Equal(t, pending.Timestamp, batchTimestamp.Seconds)
	assert.Equal(t, pending.Args, args)

	dataIn, err := pb.Marshal(&proto.Batch{TxIDs: [][]byte{ser.testIDBytes}})
	assert.NoError(t, err)

	resp := chainCode.batchExecute(mockStub, ser.testhexCreatorSKI, string(dataIn))
	assert.NotNil(t, resp)
	assert.Equal(t, resp.GetStatus(), int32(200))

	if !ser.paramsWrongON {
		response := &proto.BatchResponse{}
		err = pb.Unmarshal(resp.GetPayload(), response)
		assert.NoError(t, err)
		assert.Equal(t, len(response.TxResponses), 1)
		assert.Equal(t, response.TxResponses[0].Id, txIDBytes)
		assert.Equal(t, response.TxResponses[0].Method, testFunctionBatch)
	}
}

// TestBatchedTxExecute - positive test for batchedTxExecute
func TestBatchedTxExecute(t *testing.T) {
	chainCode, err := NewChainCode(&testContract{}, allowedMspID, nil)
	assert.NoError(t, err)

	chainCode.init = &proto.InitArgs{}
	mockStub := stub.NewMockStub(testChaincodeName, chainCode)

	mockStub.TxID = testEncodedTxID

	btchStub := newBatchStub(mockStub)

	mockStub.MockTransactionStart(testEncodedTxID)

	batchTimestamp, err := mockStub.GetTxTimestamp()
	assert.NoError(t, err)

	err = chainCode.saveToBatch(mockStub, testFunctionBatch, creatorSKI, nil, args, uint64(batchTimestamp.Seconds))
	assert.NoError(t, err)
	mockStub.MockTransactionEnd(testEncodedTxID)

	resp, event := chainCode.batchedTxExecute(btchStub, txIDBytes, batchTimestamp.Seconds)
	assert.NotNil(t, resp)
	assert.NotNil(t, event)
	assert.Nil(t, resp.Error)
	assert.Nil(t, event.Error)
}

// TestBatchedTxDelete - positive test for batchedTxDelete
func TestBatchedTxDelete(t *testing.T) {
	chainCode, err := NewChainCode(&testContract{}, allowedMspID, nil)
	assert.NoError(t, err)

	chainCode.init = &proto.InitArgs{}
	mockStub := stub.NewMockStub(testChaincodeName, chainCode)
	mockStub.TxID = testEncodedTxID
	mockStub.MockTransactionStart(testEncodedTxID)

	batchTimestamp, err := mockStub.GetTxTimestamp()
	assert.NoError(t, err)

	err = chainCode.saveToBatch(mockStub, testFunctionBatch, creatorSKI, nil, args, uint64(batchTimestamp.Seconds))
	assert.NoError(t, err)
	mockStub.MockTransactionEnd(testEncodedTxID)

	state, err := mockStub.GetState(fmt.Sprintf("\u0000batchTransactions\u0000%s\u0000", testEncodedTxID))
	assert.NotNil(t, state)
	assert.NoError(t, err)

	pending := new(proto.PendingTx)
	err = pb.Unmarshal(state, pending)
	assert.NoError(t, err)

	assert.Equal(t, pending.Method, testFunctionBatch)
	assert.Equal(t, pending.CreatorSKI, creatorSKI)
	assert.Equal(t, pending.Timestamp, batchTimestamp.Seconds)
	assert.Equal(t, pending.Args, args)

	batchedTxDelete(mockStub, batchKey, testEncodedTxID)

	stateAfterDel, err := mockStub.GetState(fmt.Sprintf("\u0000batchTransactions\u0000%s\u0000", testEncodedTxID))
	assert.Nil(t, stateAfterDel)
	assert.NoError(t, err)
}

// TestOkTxExecuteWithTTL - positive test for batchedTxExecute whit ttl
func TestOkTxExecuteWithTTL(t *testing.T) {
	chainCode, err := NewChainCode(&testContract{}, allowedMspID, &ContractOptions{
		TxTTL: 5,
	})
	assert.NoError(t, err)

	chainCode.init = &proto.InitArgs{}
	mockStub := stub.NewMockStub(testChaincodeName, chainCode)
	mockStub.TxID = testEncodedTxID
	btchStub := newBatchStub(mockStub)
	mockStub.MockTransactionStart(testEncodedTxID)

	batchTimestamp, err := mockStub.GetTxTimestamp()
	assert.NoError(t, err)

	err = chainCode.saveToBatch(mockStub, testFunctionBatch, creatorSKI, nil, args, uint64(batchTimestamp.Seconds))
	assert.NoError(t, err)
	mockStub.MockTransactionEnd(testEncodedTxID)

	resp, event := chainCode.batchedTxExecute(btchStub, txIDBytes, batchTimestamp.Seconds)
	assert.NotNil(t, resp)
	assert.NotNil(t, event)
	assert.Nil(t, resp.Error)
	assert.Nil(t, event.Error)

	assert.Equal(t, resp.Id, txIDBytes)
	assert.Equal(t, resp.Method, testFunctionBatch)
	assert.Equal(t, event.Id, txIDBytes)
	assert.Equal(t, event.Method, testFunctionBatch)
}

// TestFalseTxExecuteWithTTL - negative test for batchedTxExecute whit ttl
func TestFailTxExecuteWithTTL(t *testing.T) {
	chainCode, err := NewChainCode(&testContract{}, allowedMspID, &ContractOptions{
		TxTTL: 5,
	})
	assert.NoError(t, err)

	chainCode.init = &proto.InitArgs{}
	mockStub := stub.NewMockStub(testChaincodeName, chainCode)
	mockStub.TxID = testEncodedTxID
	btchStub := newBatchStub(mockStub)
	mockStub.MockTransactionStart(testEncodedTxID)

	batchTimestamp, err := mockStub.GetTxTimestamp()
	assert.NoError(t, err)

	err = chainCode.saveToBatch(mockStub, testFunctionBatch, creatorSKI, nil, args, uint64(batchTimestamp.Seconds))
	assert.NoError(t, err)
	mockStub.MockTransactionEnd(testEncodedTxID)

	resp, event := chainCode.batchedTxExecute(btchStub, txIDBytes, batchTimestamp.Seconds+6)
	assert.NotNil(t, resp)
	assert.NotNil(t, event)
	assert.NotNil(t, resp.Error)
	assert.NotNil(t, event.Error)
	assert.Equal(t, resp.Id, txIDBytes)
	assert.Equal(t, event.Id, txIDBytes)
	assert.Equal(t, resp.Error.Error, "function and args loading error: transaction expired")
	assert.Equal(t, event.Error.Error, "function and args loading error: transaction expired")
}
