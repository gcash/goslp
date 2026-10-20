// Package goslp has helpers for Simple Ledger Protocol (SLP) token
// transactions that work across SLP versions and token types.  Its v1parser
// package parses SLP metadata, and its metadatamaker package builds it.
package goslp

import (
	"errors"
	"fmt"

	"github.com/gcash/bchd/wire"
	"github.com/gcash/goslp/v1parser"
)

// GetSlpVersionType returns the token type of the SLP message in slpPkScript.
// It returns an error for any script ParseSLP rejects, which includes every
// token type other than 0x01, 0x41 and 0x81.
func GetSlpVersionType(slpPkScript []byte) (*uint8, error) {

	// TODO: replace the following with a method to only parse the version number
	slpMsg, err := v1parser.ParseSLP(slpPkScript)
	if err != nil {
		return nil, errors.New("unable to parse slp version")
	}
	tokenType := uint8(slpMsg.TokenType())
	return &tokenType, nil
}

// GetSlpTokenID parses the SLP message in the first output of tx and returns
// its token ID, as GetSlpTokenIDFromMsg does.
func GetSlpTokenID(tx *wire.MsgTx) ([]byte, error) {

	if len(tx.TxOut) < 1 {
		return nil, errors.New("transaction has no outputs")
	}

	slpMsg, err := v1parser.ParseSLP(tx.TxOut[0].PkScript)
	if err != nil {
		return nil, err
	}

	return GetSlpTokenIDFromMsg(tx, slpMsg)
}

// GetSlpTokenIDFromMsg returns the token ID of tx, whose first output
// ParseSLP parsed as slpMsg.  A GENESIS creates a token whose ID is its own
// transaction hash, in the byte order MINT and SEND messages carry it in.
func GetSlpTokenIDFromMsg(tx *wire.MsgTx, slpMsg v1parser.ParseResult) ([]byte, error) {
	switch msg := slpMsg.(type) {
	case *v1parser.SlpGenesis:
		hash := tx.TxHash()
		var tokenID []byte
		// reverse the bytes here since tokenID is coming from txn hash
		for i := len(hash[:]) - 1; i >= 0; i-- {
			tokenID = append(tokenID, hash[i])
		}
		return tokenID, nil
	case *v1parser.SlpMint:
		return msg.TokenID(), nil
	case *v1parser.SlpSend:
		return msg.TokenID(), nil
	default:
		return nil, fmt.Errorf("unknown SLP message type %T", slpMsg)
	}
}

//nolint:unused // GetSlpTokenID no longer needs it.
func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
