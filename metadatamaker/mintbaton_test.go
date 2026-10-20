package metadatamaker_test

import (
	"testing"

	"github.com/gcash/goslp/metadatamaker"
	"github.com/gcash/goslp/v1parser"
)

// TestNewMintBatonVout checks that a caller outside the package can put the
// mint baton of a GENESIS or MINT on the output it chooses.
func TestNewMintBatonVout(t *testing.T) {
	genesis, err := metadatamaker.TokenType1Genesis([]byte("TEST"), []byte("some name"), nil, nil, 0,
		metadatamaker.NewMintBatonVout(3), 1)
	if err != nil {
		t.Fatal(err)
	}
	slpMsg, err := v1parser.ParseSLP(genesis)
	if err != nil {
		t.Fatal(err)
	}
	if vout := slpMsg.(*v1parser.SlpGenesis).MintBatonVout; vout != 3 {
		t.Errorf("GENESIS puts the mint baton on output %d, want 3", vout)
	}

	mint, err := metadatamaker.TokenType1Mint(make([]byte, 32), metadatamaker.NewMintBatonVout(255), 1)
	if err != nil {
		t.Fatal(err)
	}
	slpMsg, err = v1parser.ParseSLP(mint)
	if err != nil {
		t.Fatal(err)
	}
	if vout := slpMsg.(*v1parser.SlpMint).MintBatonVout; vout != 255 {
		t.Errorf("MINT puts the mint baton on output %d, want 255", vout)
	}
}
