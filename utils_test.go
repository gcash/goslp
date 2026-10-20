package goslp_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/gcash/bchd/wire"
	"github.com/gcash/goslp"
	"github.com/gcash/goslp/metadatamaker"
	"github.com/gcash/goslp/v1parser"
)

func TestGetSlpTokenIDGenesis(t *testing.T) {
	txnHex := "0100000001fb6080a6ca752a808e4f86a6225164b3348b10572610d85001d00d1a9b151629030000006441e9e0035a15773bf8f86b65415c4827a9cec018abe15bc162b9974f3001a0a5ff751d35500837b9e8e77d322a8289fa92ff718be4701e5fa38f3e2726c729fb7b412102afef5c197947afa712fd6094935531935d24834cb2f0fefd811691e7230eb82bfeffffff040000000000000000416a04534c500001810747454e45534953034244441b426974636f696e20446f6e6174696f6e73204469726563746f72794c004c000100010208000000000000000122020000000000001976a914294e1c12d3f976f2dd5bd10467c4c605d6996b8e88ac22020000000000001976a914294e1c12d3f976f2dd5bd10467c4c605d6996b8e88ac9d350200000000001976a914e7abe33c8b9d58366b3114a8979509fc80420ad288ac52050a00"

	tx := wire.NewMsgTx(1)
	serializedTx, err := hex.DecodeString(txnHex)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = tx.BchDecode(bytes.NewReader(serializedTx), wire.ProtocolVersion, wire.LatestEncoding)
	if err != nil {
		t.Fatal(err.Error())
	}

	tokenID, err := goslp.GetSlpTokenID(tx)
	if err != nil {
		t.Fatal(err)
	}
	// MINT and SEND messages carry a token ID in the byte order a
	// transaction hash is displayed in.
	if got, want := hex.EncodeToString(tokenID), tx.TxHash().String(); got != want {
		t.Errorf("token ID %s, want the transaction hash %s", got, want)
	}

	slpMsg, err := v1parser.ParseSLP(tx.TxOut[0].PkScript)
	if err != nil {
		t.Fatal(err)
	}
	fromMsg, err := goslp.GetSlpTokenIDFromMsg(tx, slpMsg)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(fromMsg, tokenID) {
		t.Errorf("GetSlpTokenIDFromMsg returned %x, want %x", fromMsg, tokenID)
	}
}

// TestGetSlpTokenIDMintSend checks that the token ID of a MINT or SEND is the
// one its message carries.
func TestGetSlpTokenIDMintSend(t *testing.T) {
	tokenID := bytes.Repeat([]byte{0x01, 0x23, 0x45, 0x67}, 8)
	mint, err := metadatamaker.TokenType1Mint(tokenID, nil, 1)
	if err != nil {
		t.Fatal(err)
	}
	send, err := metadatamaker.TokenType1Send(tokenID, []uint64{1})
	if err != nil {
		t.Fatal(err)
	}

	for name, script := range map[string][]byte{"MINT": mint, "SEND": send} {
		tx := wire.NewMsgTx(1)
		tx.AddTxOut(&wire.TxOut{PkScript: script})

		got, err := goslp.GetSlpTokenID(tx)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if !bytes.Equal(got, tokenID) {
			t.Errorf("%s: GetSlpTokenID returned %x, want %x", name, got, tokenID)
		}

		slpMsg, err := v1parser.ParseSLP(script)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		got, err = goslp.GetSlpTokenIDFromMsg(tx, slpMsg)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if !bytes.Equal(got, tokenID) {
			t.Errorf("%s: GetSlpTokenIDFromMsg returned %x, want %x", name, got, tokenID)
		}
	}
}

// otherMsg is a ParseResult that ParseSLP never returns.
type otherMsg struct{ v1parser.ParseResult }

// TestGetSlpTokenIDErrors checks the transactions and messages that have no
// token ID.
func TestGetSlpTokenIDErrors(t *testing.T) {
	if _, err := goslp.GetSlpTokenID(wire.NewMsgTx(1)); err == nil {
		t.Error("transaction with no outputs: got a token ID, want an error")
	}

	tx := wire.NewMsgTx(1)
	tx.AddTxOut(&wire.TxOut{Value: 1000, PkScript: []byte{0x76, 0xa9}})
	if _, err := goslp.GetSlpTokenID(tx); err == nil {
		t.Error("non-SLP first output: got a token ID, want an error")
	}

	if _, err := goslp.GetSlpTokenIDFromMsg(tx, nil); err == nil {
		t.Error("nil message: got a token ID, want an error")
	}
	if _, err := goslp.GetSlpTokenIDFromMsg(tx, otherMsg{}); err == nil {
		t.Error("unknown message type: got a token ID, want an error")
	}
}
