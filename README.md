# goslp

[![Build Status](https://github.com/gcash/goslp/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/gcash/goslp/actions/workflows/ci.yml)

Golang packages for the Simple Ledger Protocol (SLP).

This is a copy of `github.com/simpleledgerinc/goslp` at commit 3c2e5f2ef33f,
the version bchd used, taken from the Go module cache after the original
repository was deleted. Its earlier history is not included.

    go get github.com/gcash/goslp

### v1parser - for parsing transaction metadata

This package is used for parsing SLP metadata from the scriptPubKey of the SLP transaction's first output.

```go
// here is an example marshaled SLP message extracted from transaction's scriptPubKey.
scriptPubKey, err := hex.DecodeString("6a04534c500001010453454e4420c4b0d62156b3fa5c8f3436079b5394f7edc1bef5dc1cd2f9d0c4d46f82cca47908000000000000000108000000000000000408000000000000005a")
if err != nil {
    // handle the error ...
}

// get the unmarshalled slp message
slpMsg, err := v1parser.ParseSLP(scriptPubKey)

// do something ...

```

The gRPC server in [bchd](https://github.com/gcash/bchd/blob/master/bchrpc/server.go) uses the unmarshalled SLP metadata object throughout, which makes it a good example of how to work with it.

Differential fuzzer testing has been performed with the [slp-validate.js](https://github.com/simpleledger/slp-validate) npm package, and can be reproduced following the instructions in the `./fuzz` directory.

### metadatamaker - for creating new transaction metadata

This package is used for creating marshaled SLP metadata for adding to a transaction's first output scriptPubKey.  Helper methods are provided for Type 1, NFT Group, and NFT children transactions.

**Genesis** - use CreateOpReturnGenesis, NFT1GroupGenesis, or NFT1ChildGenesis

```go
// put the mint baton on output 2, or pass nil for no baton
mintBatonVout := NewMintBatonVout(2)

scriptPubKey, err := CreateOpReturnGenesis(
      versionType,
      ticker,
      name,
      documentURL,
      documentHash,
      decimals,
      mintBatonVout,
      quantity,
)

```



**Mint** - use CreateOpReturnMint or NFT1GroupMint

```go
// pass the baton on to output 2, or pass nil to end the minting
mintBatonVout := NewMintBatonVout(2)

scriptPubKey, err := CreateOpReturnMint(
    versionType,
    tokenID,
    mintBatonVout,
    quantity,
)

```



**Send** - use CreateOpReturnSend, NFT1GroupSend, or NFT1ChildSend

```go
scriptPubKey, err := CreateOpReturnSend(
		1,
		tokenID,
		[]uint64{1, 2},
)

```

