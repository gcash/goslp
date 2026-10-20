package metadatamaker

import (
	"testing"

	"github.com/gcash/goslp/v1parser"
)

func TestNft1ChildCreateOpReturnGenesis(t *testing.T) {
	ticker := []byte("TEST")
	name := []byte("some name")
	documentURL := []byte("")
	documentHash := []byte("")
	decimals := 0
	quantity := uint64(1)

	slpMsg, err := NFT1ChildGenesis(
		ticker,
		name,
		documentURL,
		documentHash,
		decimals,
		quantity,
	)
	if err != nil {
		t.Fatal(err.Error())
	}
	if slpMsg == nil {
		t.Fatal("parameters were not marshalled")
	}
	_, err = v1parser.ParseSLP(slpMsg)
	if err != nil {
		t.Error(err.Error())
	}
}

// TestCreateOpReturnMintTokenTypes checks that CreateOpReturnMint builds a
// MINT that ParseSLP accepts for each token type that can mint, and refuses
// one for an NFT1 child, which cannot.
func TestCreateOpReturnMintTokenTypes(t *testing.T) {
	for _, versionType := range []int{0x01, 0x81} {
		script, err := CreateOpReturnMint(versionType, make([]byte, 32), NewMintBatonVout(2), 1)
		if err != nil {
			t.Fatalf("token type %#x: %v", versionType, err)
		}
		if _, err := v1parser.ParseSLP(script); err != nil {
			t.Errorf("token type %#x: ParseSLP: %v", versionType, err)
		}
	}
	if _, err := CreateOpReturnMint(0x41, make([]byte, 32), nil, 1); err == nil {
		t.Error("token type 0x41: built a MINT, want an error")
	}
}

func TestNft1ChildCreateOpReturnSend(t *testing.T) {
	slpMsg, err := NFT1ChildSend(
		make([]byte, 32),
		[]uint64{0, 1},
	)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = v1parser.ParseSLP(slpMsg)
	if err != nil {
		t.Error(err.Error())
	}
}
