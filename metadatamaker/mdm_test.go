package metadatamaker

import (
	"bytes"
	"testing"

	"github.com/gcash/goslp/v1parser"
)

func TestCreateOpReturnGenesis(t *testing.T) {
	versionType := 1
	ticker := []byte("TEST")
	name := []byte("some name")
	documentURL := []byte("")
	documentHash := []byte("")
	decimals := 0
	mintBatonVout := &MintBatonVout{vout: 2}
	quantity := uint64(1)

	slpMsg, err := CreateOpReturnGenesis(
		versionType,
		ticker,
		name,
		documentURL,
		documentHash,
		decimals,
		mintBatonVout,
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

func TestCreateOpReturnGenesisBadMintVout(t *testing.T) {
	vouts := []*MintBatonVout{{vout: 0}, {vout: 1}, {vout: 256}}
	for _, vout := range vouts {
		versionType := 1
		ticker := []byte("TEST")
		name := []byte("some name")
		documentURL := []byte("")
		documentHash := []byte("")
		decimals := 0
		quantity := uint64(1)

		_, err := CreateOpReturnGenesis(
			versionType,
			ticker,
			name,
			documentURL,
			documentHash,
			decimals,
			vout,
			quantity,
		)
		if err == nil || err.Error() != "mintBatonVout out of range (0x02 < > 0xFF)" {
			t.Fatalf("vout %d: got error %v, want out of range", vout.vout, err)
		}
	}
}

func TestCreateOpReturnMintBadMintVout(t *testing.T) {
	vouts := []*MintBatonVout{{vout: 0}, {vout: 1}, {vout: 256}}
	for _, vout := range vouts {
		versionType := 1
		tokenID := make([]byte, 32)
		quantity := uint64(1)

		_, err := CreateOpReturnMint(
			versionType,
			tokenID,
			vout,
			quantity,
		)
		if err == nil || err.Error() != "mintBatonVout out of range (0x02 < > 0xFF)" {
			t.Errorf("vout %d: got error %v, want out of range", vout.vout, err)
		}
	}
}

func TestCreateOpReturnSend(t *testing.T) {
	slpMsg, err := CreateOpReturnSend(
		1,
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

// TestPushSlpDataOpcodes checks that each push uses the smallest opcode that
// fits its length, and that an empty push is OP_PUSHDATA1 with a zero length,
// since SLP does not allow OP_0.
func TestPushSlpDataOpcodes(t *testing.T) {
	tests := []struct {
		len    int
		prefix []byte
	}{
		{0, []byte{0x4c, 0x00}},
		{1, []byte{0x01}},
		{75, []byte{0x4b}},
		{76, []byte{0x4c, 0x4c}},
		{77, []byte{0x4c, 0x4d}},
		{255, []byte{0x4c, 0xff}},
		{256, []byte{0x4d, 0x00, 0x01}},
		{65535, []byte{0x4d, 0xff, 0xff}},
		{65536, []byte{0x4e, 0x00, 0x00, 0x01, 0x00}},
	}
	for _, test := range tests {
		data := bytes.Repeat([]byte{0xaa}, test.len)
		push, err := pushSlpData(data)
		if err != nil {
			t.Fatalf("length %d: %v", test.len, err)
		}
		want := append(append([]byte{}, test.prefix...), data...)
		if !bytes.Equal(push, want) {
			n := min(len(push), len(test.prefix))
			t.Errorf("length %d: push is %d bytes starting %x, want %d bytes starting %x",
				test.len, len(push), push[:n], len(want), test.prefix)
		}
	}
}

// TestCreateOpReturnGenesisFieldSizes checks that a GENESIS ticker, name or
// document URL comes back from ParseSLP byte for byte at the lengths where
// the push opcode changes.
func TestCreateOpReturnGenesisFieldSizes(t *testing.T) {
	names := []string{"ticker", "name", "document URL"}
	for _, n := range []int{75, 76, 77, 255, 256, 65535, 65536} {
		// A field starting with 0x4b ("K") written after a bare length
		// byte of 0x4c reads as OP_PUSHDATA1 of 75 bytes, which takes the
		// rest of a 76-byte field and still parses as a valid GENESIS.
		field := append([]byte("K"), bytes.Repeat([]byte("x"), n-1)...)
		for i := range names {
			fields := [][]byte{[]byte("TEST"), []byte("some name"), []byte("")}
			fields[i] = field
			script, err := CreateOpReturnGenesis(0x01, fields[0], fields[1], fields[2], nil, 0, nil, 1)
			if err != nil {
				t.Fatalf("%d-byte %s: %v", n, names[i], err)
			}
			slpMsg, err := v1parser.ParseSLP(script)
			if err != nil {
				t.Errorf("%d-byte %s: %v", n, names[i], err)
				continue
			}
			genesis := slpMsg.(*v1parser.SlpGenesis)
			got := [][]byte{genesis.Ticker, genesis.Name, genesis.DocumentURI}
			for j := range fields {
				if !bytes.Equal(got[j], fields[j]) {
					t.Errorf("%d-byte %s: %s came back as %d bytes, want %d",
						n, names[i], names[j], len(got[j]), len(fields[j]))
				}
			}
		}
	}
}
