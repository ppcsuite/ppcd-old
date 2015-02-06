// Copyright (c) 2014 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

import (
	"math/big"
	"time"

	"github.com/ppcsuite/ppcd/wire"
)

// genesisCoinbaseTx is the coinbase transaction for the genesis blocks for
// the main network, regression test network, and test network (version 3).
var genesisCoinbaseTx = wire.MsgTx{
	Version: 1,
	Time:    time.Unix(0x502c59a2, 0),
	TxIn: []*wire.TxIn{
		{
			PreviousOutPoint: wire.OutPoint{
				Hash:  wire.ShaHash{},
				Index: 0xffffffff,
			},
			SignatureScript: []byte{
				0x04, 0xff, 0xff, 0x00, 0x1d, 0x02, 0x0f, 0x27,
				0x4b, 0x4d, 0x61, 0x74, 0x6f, 0x6e, 0x69, 0x73,
				0x20, 0x30, 0x37, 0x2d, 0x41, 0x55, 0x47, 0x2d,
				0x32, 0x30, 0x31, 0x32, 0x20, 0x50, 0x61, 0x72,
				0x61, 0x6c, 0x6c, 0x65, 0x6c, 0x20, 0x43, 0x75,
				0x72, 0x72, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73,
				0x20, 0x41, 0x6e, 0x64, 0x20, 0x54, 0x68, 0x65,
				0x20, 0x52, 0x6f, 0x61, 0x64, 0x6d, 0x61, 0x70,
				0x20, 0x54, 0x6f, 0x20, 0x4d, 0x6f, 0x6e, 0x65,
				0x74, 0x61, 0x72, 0x79, 0x20, 0x46, 0x72, 0x65,
				0x65, 0x64, 0x6f, 0x6d,
			},
			Sequence: 0xffffffff,
		},
	},
	TxOut: []*wire.TxOut{
		{
			Value:    0x00,
			PkScript: []byte{},
		},
	},
	LockTime: 0,
}

// genesisHash is the hash of the first block in the block chain for the main
// network (genesis block).
var genesisHash = wire.ShaHash([wire.HashSize]byte{ // Make go vet happy.
	0xe3, 0x27, 0xcd, 0x80, 0xc8, 0xb1, 0x7e, 0xfd,
	0xa4, 0xea, 0x08, 0xc5, 0x87, 0x7e, 0x95, 0xd8,
	0x77, 0x46, 0x2a, 0xb6, 0x63, 0x49, 0xd5, 0x66,
	0x71, 0x67, 0xfe, 0x32, 0x00, 0x00, 0x00, 0x00,
})

// genesisMerkleRoot is the hash of the first transaction in the genesis block
// for the main network.
var genesisMerkleRoot = wire.ShaHash([wire.HashSize]byte{ // Make go vet happy.
	0xc2, 0x93, 0x59, 0x2c, 0x05, 0x90, 0x56, 0x98,
	0x29, 0x0c, 0x89, 0xeb, 0x6d, 0xde, 0xf0, 0xcf,
	0x8a, 0xa5, 0xa1, 0x48, 0xc6, 0x8c, 0x55, 0xac,
	0x7a, 0xd1, 0xb4, 0xfa, 0x85, 0x8f, 0x2d, 0x3c,
})

// genesisBlock defines the genesis block of the block chain which serves as the
// public transaction ledger for the main network.
var genesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  wire.ShaHash{},           // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: genesisMerkleRoot,        //
		Timestamp:  time.Unix(0x502c5b7f, 0), //
		Bits:       0x1d00ffff,               //
		Nonce:      0x81e582ab,               //
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
	Signature:    []byte{},
}

var genesisMeta = wire.Meta{
	StakeModifier:         0,
	StakeModifierChecksum: 0x0e00670b,
	Flags:      6,
	ChainTrust: *big.NewInt(1),
}

// regTestGenesisHash is the hash of the first block in the block chain for the
// regression test network (genesis block).
var regTestGenesisHash = wire.ShaHash([wire.HashSize]byte{ // Make go vet happy.
	0x06, 0x22, 0x6e, 0x46, 0x11, 0x1a, 0x0b, 0x59,
	0xca, 0xaf, 0x12, 0x60, 0x43, 0xeb, 0x5b, 0xbf,
	0x28, 0xc3, 0x4f, 0x3a, 0x5e, 0x33, 0x2a, 0x1f,
	0xc7, 0xb2, 0xb7, 0x3c, 0xf1, 0x88, 0x91, 0x0f,
})

// regTestGenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the regression test network.  It is the same as the merkle root for
// the main network.
var regTestGenesisMerkleRoot = genesisMerkleRoot

// regTestGenesisBlock defines the genesis block of the block chain which serves
// as the public transaction ledger for the regression test network.
var regTestGenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  wire.ShaHash{},           // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: regTestGenesisMerkleRoot, // 4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
		Timestamp:  time.Unix(1296688602, 0), // 2011-02-02 23:16:42 +0000 UTC
		Bits:       0x207fffff,               // 545259519 [7fffff0000000000000000000000000000000000000000000000000000000000]
		Nonce:      2,
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}

// testNet3GenesisHash is the hash of the first block in the block chain for the
// test network (version 3).
var testNet3GenesisHash = wire.ShaHash([wire.HashSize]byte{ // Make go vet happy.
	0x06, 0x9f, 0x7c, 0xc4, 0xae, 0x81, 0xca, 0x0c,
	0x7c, 0x72, 0xcc, 0x30, 0xe6, 0x8c, 0x65, 0xb0,
	0x17, 0xcd, 0x17, 0x3e, 0x50, 0x96, 0x65, 0x7f,
	0x73, 0xbb, 0x57, 0xf7, 0x01, 0x00, 0x00, 0x00,
})

// testNet3GenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the test network (version 3).  It is the same as the merkle root
// for the main network.
var testNet3GenesisMerkleRoot = genesisMerkleRoot

// testNet3GenesisBlock defines the genesis block of the block chain which
// serves as the public transaction ledger for the test network (version 3).
var testNet3GenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  wire.ShaHash{},           // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: genesisMerkleRoot,        //
		Timestamp:  time.Unix(0x502C71D0, 0), //
		Bits:       0x1d0fffff,               //
		Nonce:      0x07533A5A,               //
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
	Signature:    []byte{},
}

// simNetGenesisHash is the hash of the first block in the block chain for the
// simulation test network.
var simNetGenesisHash = wire.ShaHash([wire.HashSize]byte{ // Make go vet happy.
	0xf6, 0x7a, 0xd7, 0x69, 0x5d, 0x9b, 0x66, 0x2a,
	0x72, 0xff, 0x3d, 0x8e, 0xdb, 0xbb, 0x2d, 0xe0,
	0xbf, 0xa6, 0x7b, 0x13, 0x97, 0x4b, 0xb9, 0x91,
	0x0d, 0x11, 0x6d, 0x5c, 0xbd, 0x86, 0x3e, 0x68,
})

// simNetGenesisMerkleRoot is the hash of the first transaction in the genesis
// block for the simulation test network.  It is the same as the merkle root for
// the main network.
var simNetGenesisMerkleRoot = genesisMerkleRoot

// simNetGenesisBlock defines the genesis block of the block chain which serves
// as the public transaction ledger for the simulation test network.
var simNetGenesisBlock = wire.MsgBlock{
	Header: wire.BlockHeader{
		Version:    1,
		PrevBlock:  wire.ShaHash{},           // 0000000000000000000000000000000000000000000000000000000000000000
		MerkleRoot: simNetGenesisMerkleRoot,  // 4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
		Timestamp:  time.Unix(1401292357, 0), // 2014-05-28 15:52:37 +0000 UTC
		Bits:       0x207fffff,               // 545259519 [7fffff0000000000000000000000000000000000000000000000000000000000]
		Nonce:      2,
	},
	Transactions: []*wire.MsgTx{&genesisCoinbaseTx},
}
