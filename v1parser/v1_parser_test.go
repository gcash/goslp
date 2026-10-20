package v1parser

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math"
	"math/big"
	"os"
	"reflect"
	"testing"
)

func TestSlpMessageUnitTests(t *testing.T) {

	// parsing tests are from:
	// https://github.com/simpleledger/slp-unit-test-data#part-a-parsing-of-op_return-scripts
	inputTestsFile, err := os.Open("v1_parser_test_opreturn.json")
	if err != nil {
		t.Fatal(err.Error())
	}
	data, err := io.ReadAll(inputTestsFile)
	if err != nil {
		t.Fatal(err)
	}
	defer inputTestsFile.Close()

	// TestCase is the json object parsed from v1_parser_test_op_return.json
	// Code is for invalid reasons, for meanings see:
	// https://github.com/simpleledger/slp-unit-test-data#invalidation-reason-codes
	type TestCase struct {
		Msg    string
		Script string
		Code   *float64
	}
	var tests []TestCase
	err = json.Unmarshal(data, &tests)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(tests) == 0 {
		t.Fatal("no test vectors")
	}
	for i, test := range tests {
		slpbuf, err := hex.DecodeString(test.Script)
		if err != nil {
			t.Fatalf("Test %d (%s): bad script hex: %v", i, test.Msg, err)
		}
		_, err = ParseSLP(slpbuf)
		switch {
		case test.Code == nil:
			if err != nil {
				t.Errorf("Test %d (%s): got error %v, want a valid parse", i, test.Msg, err)
			}
		case err == nil:
			t.Errorf("Test %d (%s): parsed, want an error with code %v", i, test.Msg, *test.Code)
		case *test.Code == 255 && !errors.Is(err, ErrUnsupportedSlpVersion):
			t.Errorf("Test %d (%s): got error %v, want %v", i, test.Msg, err, ErrUnsupportedSlpVersion)
		case *test.Code != 255 && errors.Is(err, ErrUnsupportedSlpVersion):
			t.Errorf("Test %d (%s): got %v for code %v", i, test.Msg, err, *test.Code)
		}
	}
}

func TestGetOutputAmountSend(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c500001010453454e4420c4b0d62156b3fa5c8f3436079b5394f7edc1bef5dc1cd2f9d0c4d46f82cca47908000000000000000108000000000000000408000000000000005a")
	slpMsg, err := ParseSLP(scriptPubKey)
	if err != nil {
		t.Fatal(err)
	}
	amt, _ := slpMsg.GetVoutValue(3)
	if amt.Cmp(big.NewInt(90)) != 0 {
		t.Error("incorrect amount parsed for index")
	}
}

func TestGetOutputAmountSendNotNegative(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c500001010453454e4420c4b0d62156b3fa5c8f3436079b5394f7edc1bef5dc1cd2f9d0c4d46f82cca47908ffffffffffffffff")
	slpMsg, err := ParseSLP(scriptPubKey)
	if err != nil {
		t.Fatal(err)
	}
	amt, _ := slpMsg.GetVoutValue(1)
	if amt.Cmp(big.NewInt(0)) < 1 {
		t.Error("amount is less than zero")
	}
	if amt.Cmp(new(big.Int).SetUint64(18446744073709551615)) != 0 {
		t.Error("amount is incorrect")
	}
}

func TestGetOutputAmountMint(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c50000101044d494e5420d6876f0fce603be43f15d34348bb1de1a8d688e1152596543da033a060cff798010208000000000bebc200")
	slpMsg, err := ParseSLP(scriptPubKey)
	if err != nil {
		t.Fatal(err)
	}
	amt, _ := slpMsg.GetVoutValue(1)
	if amt.Cmp(big.NewInt(200000000)) != 0 {
		t.Error("incorrect amount parsed for index")
	}
}

func TestGetOutputAmountMintVout2(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c50000101044d494e5420d6876f0fce603be43f15d34348bb1de1a8d688e1152596543da033a060cff798010208000000000bebc200")
	slpMsg, err := ParseSLP(scriptPubKey)
	if err != nil {
		t.Fatal(err)
	}
	amt, isBaton := slpMsg.GetVoutValue(2)
	if isBaton != true {
		t.Error("this is a baton")
	}
	if amt != nil {
		t.Error("incorrect amount parsed for index")
	}
}

func TestGetOutputAmountMintVoutNotNegative(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c50000101044d494e5420d6876f0fce603be43f15d34348bb1de1a8d688e1152596543da033a060cff798010208ffffffffffffffff")
	slpMsg, err := ParseSLP(scriptPubKey)
	if err != nil {
		t.Fatal(err)
	}
	amt, _ := slpMsg.GetVoutValue(1)
	if amt.Cmp(big.NewInt(0)) < 1 {
		t.Error("amount is less than zero")
	}
	if amt.Cmp(new(big.Int).SetUint64(18446744073709551615)) != 0 {
		t.Error("amount is incorrect")
	}
}

func TestGetOutputAmountGenesis(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c500001010747454e45534953074f6e65436f696e074f6e65436f696e4c5468747470733a2f2f7468656e6578747765622e636f6d2f68617264666f726b2f323031392f31322f32332f6f6e65636f696e2d63727970746f63757272656e63792d7363616d2d6e6565642d746f2d6b6e6f772f4c00010401020800000061c9f36800")
	slpMsg, err := ParseSLP(scriptPubKey)
	if err != nil {
		t.Fatal(err)
	}
	amt, _ := slpMsg.GetVoutValue(1)
	if amt.Cmp(big.NewInt(420000000000)) != 0 {
		t.Error("incorrect amount parsed for index")
	}
}

func TestGetOutputAmountGenesisVout2(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c500001010747454e45534953074f6e65436f696e074f6e65436f696e4c5468747470733a2f2f7468656e6578747765622e636f6d2f68617264666f726b2f323031392f31322f32332f6f6e65636f696e2d63727970746f63757272656e63792d7363616d2d6e6565642d746f2d6b6e6f772f4c00010401020800000061c9f36800")
	slpMsg, err := ParseSLP(scriptPubKey)
	if err != nil {
		t.Fatal(err)
	}
	amt, isBaton := slpMsg.GetVoutValue(2)
	if isBaton != true {
		t.Error("this is a baton")
	}
	if amt != nil {
		t.Error("incorrect amount parsed for index")
	}
}

func TestGetOutputAmountGenesisVoutNotNegative(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c500001010747454e45534953074f6e65436f696e074f6e65436f696e4c5468747470733a2f2f7468656e6578747765622e636f6d2f68617264666f726b2f323031392f31322f32332f6f6e65636f696e2d63727970746f63757272656e63792d7363616d2d6e6565642d746f2d6b6e6f772f4c000104010208ffffffffffffffff")
	slpMsg, err := ParseSLP(scriptPubKey)
	if err != nil {
		t.Fatal(err)
	}
	amt, _ := slpMsg.GetVoutValue(1)
	if amt.Cmp(big.NewInt(0)) < 1 {
		t.Error("amount is less than zero")
	}
	if amt.Cmp(new(big.Int).SetUint64(18446744073709551615)) != 0 {
		t.Error("amount is incorrect")
	}
}

func TestTotalOutputAmountSend(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c500001010453454e4420c4b0d62156b3fa5c8f3436079b5394f7edc1bef5dc1cd2f9d0c4d46f82cca47908000000000000000108000000000000000408000000000000005a")
	slpMsg, err := ParseSLP(scriptPubKey)
	if err != nil {
		t.Fatal(err)
	}
	amt, _ := slpMsg.TotalSlpMsgOutputValue()
	if amt.Cmp(big.NewInt(95)) != 0 {
		t.Error("incorrect total amount from SEND script")
	}
}

func TestGetTotalOutputAmountSendNotNegative(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c500001010453454e442044e3d05a07091091a63a4074287a784fcd96c26095682e05c22c4bd4e5bf8681080000000000000001080000000000000001080000000000000001080000000000000001080000000000000001080000000000000001080000000000000001080000000000000001088ac7230489e80000")
	slpMsg, err := ParseSLP(scriptPubKey)
	if err != nil {
		t.Fatal(err)
	}
	amt, _ := slpMsg.TotalSlpMsgOutputValue()
	if amt.Cmp(big.NewInt(0)) < 1 {
		t.Error("amount is less than zero")
	}
}

func TestTotalOutputAmountMint(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c50000101044d494e5420d6876f0fce603be43f15d34348bb1de1a8d688e1152596543da033a060cff798010208000000000bebc200")
	slpMsg, err := ParseSLP(scriptPubKey)
	if err != nil {
		t.Fatal(err)
	}
	amt, _ := slpMsg.TotalSlpMsgOutputValue()
	if amt.Cmp(big.NewInt(200000000)) != 0 {
		t.Error("incorrect total amount from script")
	}
}

func TestTotalOutputAmountGenesis(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c500001010747454e45534953074f6e65436f696e074f6e65436f696e4c5468747470733a2f2f7468656e6578747765622e636f6d2f68617264666f726b2f323031392f31322f32332f6f6e65636f696e2d63727970746f63757272656e63792d7363616d2d6e6565642d746f2d6b6e6f772f4c00010401020800000061c9f36800")
	slpMsg, err := ParseSLP(scriptPubKey)
	if err != nil {
		t.Fatal(err)
	}
	amt, _ := slpMsg.TotalSlpMsgOutputValue()
	if amt.Cmp(big.NewInt(420000000000)) != 0 {
		t.Error("incorrect total amount from script")
	}
}

func TestGenesisParseSlp(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c500001010747454e455349534c004c004c004c0001004c0008ffffffffffffffff")
	slpMsg, err := ParseSLP(scriptPubKey)
	if err != nil {
		t.Fatal(err.Error())
	}
	if slpMsg.(*SlpGenesis).Qty != 18446744073709551615 {
		t.Error("incorrect mint qty")
	}
}

func TestSendParseSlp(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c500001010453454e4420d6876f0fce603be43f15d34348bb1de1a8d688e1152596543da033a060cff798080000000165a0bc00")
	slpMsg, err := ParseSLP(scriptPubKey)
	if err != nil {
		t.Fatal(err.Error())
	}
	if slpMsg.(*SlpSend).Amounts[0] != 6000000000 {
		t.Error("incorrect send qty")
	}
	tokenID, err := hex.DecodeString("d6876f0fce603be43f15d34348bb1de1a8d688e1152596543da033a060cff798")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(slpMsg.(*SlpSend).TokenID(), tokenID) {
		t.Error("incorrect tokenID")
	}
}

func TestMintParseSlp(t *testing.T) {
	scriptPubKey, _ := hex.DecodeString("6a04534c50000101044d494e5420d6876f0fce603be43f15d34348bb1de1a8d688e1152596543da033a060cff7980102080000000017d78400")
	slpMsg, err := ParseSLP(scriptPubKey)
	if err != nil {
		t.Fatal(err.Error())
	}
	if slpMsg.(*SlpMint).Qty != 400000000 {
		t.Error("incorrect mint qty")
	}
	tokenID, err := hex.DecodeString("d6876f0fce603be43f15d34348bb1de1a8d688e1152596543da033a060cff798")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(slpMsg.(*SlpMint).TokenID(), tokenID) {
		t.Error("incorrect tokenID")
	}
}

// TestParseResultPointerReceivers checks that only pointers to the parse
// results implement ParseResult, as ParseSLP returns them, so that a type
// assertion to a value type fails to compile instead of never matching.
func TestParseResultPointerReceivers(t *testing.T) {
	iface := reflect.TypeOf((*ParseResult)(nil)).Elem()
	for _, v := range []any{SlpGenesis{}, SlpMint{}, SlpSend{}} {
		typ := reflect.TypeOf(v)
		if typ.Implements(iface) {
			t.Errorf("%v implements ParseResult, want only %v to", typ, reflect.PointerTo(typ))
		}
		if !reflect.PointerTo(typ).Implements(iface) {
			t.Errorf("%v does not implement ParseResult", reflect.PointerTo(typ))
		}
	}
}

// TestGetVoutValueRange checks GetVoutValue at and past the ends of each
// message's outputs.  Output 0 and negative outputs have no amount and are
// never a baton, and only the output a GENESIS or MINT names is its baton.
func TestGetVoutValueRange(t *testing.T) {
	const tokenID = "d6876f0fce603be43f15d34348bb1de1a8d688e1152596543da033a060cff798"
	genesis := func(baton string) string {
		return "6a04534c50000101" + "0747454e45534953" + "4c004c004c004c00" + "0100" + baton + "080000000000000064"
	}
	mint := func(baton string) string {
		return "6a04534c50000101" + "044d494e54" + "20" + tokenID + baton + "080000000000000064"
	}

	type want struct {
		amount *big.Int
		baton  bool
	}
	var none want
	tests := []struct {
		name   string
		script string
		wants  map[int]want
	}{
		{"GENESIS, baton on 2", genesis("0102"), map[int]want{
			math.MinInt: none, -1: none, 0: none, 1: {big.NewInt(100), false}, 2: {nil, true}, 3: none}},
		{"GENESIS, no baton", genesis("4c00"), map[int]want{
			math.MinInt: none, -1: none, 0: none, 1: {big.NewInt(100), false}, 2: none}},
		{"MINT, baton on 2", mint("0102"), map[int]want{
			math.MinInt: none, -1: none, 0: none, 1: {big.NewInt(100), false}, 2: {nil, true}, 3: none}},
		{"MINT, no baton", mint("4c00"), map[int]want{
			math.MinInt: none, -1: none, 0: none, 1: {big.NewInt(100), false}, 2: none}},
		{"SEND", validSendScript, map[int]want{
			math.MinInt: none, -1: none, 0: none, 1: {big.NewInt(6000000000), false}, 2: none}},
	}
	for _, test := range tests {
		script, err := hex.DecodeString(test.script)
		if err != nil {
			t.Fatalf("%s: %v", test.name, err)
		}
		slpMsg, err := ParseSLP(script)
		if err != nil {
			t.Fatalf("%s: %v", test.name, err)
		}
		for vout, w := range test.wants {
			amount, baton := slpMsg.GetVoutValue(vout)
			amountOK := amount == nil && w.amount == nil ||
				amount != nil && w.amount != nil && amount.Cmp(w.amount) == 0
			if !amountOK || baton != w.baton {
				t.Errorf("%s: GetVoutValue(%d) = %v, %v, want %v, %v",
					test.name, vout, amount, baton, w.amount, w.baton)
			}
		}
	}
}

// validSendScript is a type 1 SEND of 6000000000 base units to output 1.
const validSendScript = "6a04534c500001010453454e4420d6876f0fce603be43f15d34348bb1de1a8d688e1152596543da033a060cff798080000000165a0bc00"

// TestParseSLPPushPastEnd checks that a push running past the end of the
// script is rejected.  The PUSHDATA4 lengths near and above 2^31 are for
// 32-bit platforms, where they overflow an int.
func TestParseSLPPushPastEnd(t *testing.T) {
	send, err := hex.DecodeString(validSendScript)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseSLP(send); err != nil {
		t.Fatalf("valid SEND: %v", err)
	}

	tests := []struct {
		name   string
		script string
	}{
		{"PUSHDATA4 of 2^31-1 bytes", "6a4effffff7f00000000"},
		{"lokad ID, PUSHDATA4 of 2^31-1 bytes", "6a04534c50004effffff7f00"},
		{"lokad ID, PUSHDATA4 of 2^32-1 bytes", "6a04534c50004effffffff00"},
		{"SEND, PUSHDATA4 length 2^31 and no data", validSendScript + "4e00000080"},
		{"SEND, PUSHDATA4 length 2^32-1 and no data", validSendScript + "4effffffff"},
		{"SEND, truncated PUSHDATA1 data", validSendScript + "4c05aabb"},
		{"SEND, truncated PUSHDATA2 data", validSendScript + "4d0500aabb"},
		{"SEND, truncated PUSHDATA4 data", validSendScript + "4e05000000aabb"},
	}
	for _, test := range tests {
		script, err := hex.DecodeString(test.script)
		if err != nil {
			t.Fatalf("%s: %v", test.name, err)
		}
		if _, err := ParseSLP(script); err == nil {
			t.Errorf("%s: parsed, want an error", test.name)
		}
	}
}

// TestParseSLPEmptyFinalPush checks that an empty PUSHDATA1, 2 or 4 at the
// end of the script is read as an empty chunk, as slp-validate reads it.  So
// an unsupported token type is still reported as one, and a SEND whose last
// amount is empty fails.
func TestParseSLPEmptyFinalPush(t *testing.T) {
	for _, push := range []string{"4c00", "4d0000", "4e00000000"} {
		script, err := hex.DecodeString("6a04534c50000102" + push)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := ParseSLP(script); !errors.Is(err, ErrUnsupportedSlpVersion) {
			t.Errorf("token type 2, empty push %s: got error %v, want %v", push, err, ErrUnsupportedSlpVersion)
		}

		script, err = hex.DecodeString(validSendScript + push)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := ParseSLP(script); err == nil {
			t.Errorf("SEND, empty push %s: parsed, want an error", push)
		}
	}
}
