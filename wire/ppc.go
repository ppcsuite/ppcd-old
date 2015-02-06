// Copyright (c) 2014-2015 PPCD developers.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"encoding/binary"
	"io"
	"math/big"
)

var ZeroSha = ShaHash{}

type Meta struct {
	StakeModifier         uint64
	StakeModifierChecksum uint32 // checksum of index; in-memeory only (main.h)
	HashProofOfStake      ShaHash
	Flags                 uint32
	ChainTrust            big.Int
	Mint                  int64
	MoneySupply           int64
	TxOffsets             []uint32
}

func (m *Meta) Serialize(w io.Writer) error {
	e := binary.Write(w, binary.LittleEndian, &m.StakeModifier)
	if e != nil {
		return e
	}
	binary.Write(w, binary.LittleEndian, &m.StakeModifierChecksum)
	if e != nil {
		return e
	}
	binary.Write(w, binary.LittleEndian, &m.Flags)
	if e != nil {
		return e
	}
	binary.Write(w, binary.LittleEndian, &m.HashProofOfStake)
	if e != nil {
		return e
	}
	bytes := m.ChainTrust.Bytes()
	var blen byte
	blen = byte(len(bytes))
	binary.Write(w, binary.LittleEndian, &blen)
	if e != nil {
		return e
	}
	binary.Write(w, binary.LittleEndian, &bytes)
	if e != nil {
		return e
	}
	binary.Write(w, binary.LittleEndian, &m.Mint)
	if e != nil {
		return e
	}
	binary.Write(w, binary.LittleEndian, &m.MoneySupply)
	if e != nil {
		return e
	}

	binary.Write(w, binary.LittleEndian, uint32(len(m.TxOffsets)))
	if e != nil {
		return e
	}

	for _, txOffset := range m.TxOffsets {
		binary.Write(w, binary.LittleEndian, txOffset)
		if e != nil {
			return e
		}
	}

	return nil
}

func (m *Meta) Deserialize(r io.Reader) error {
	e := binary.Read(r, binary.LittleEndian, &m.StakeModifier)
	if e != nil {
		return e
	}
	e = binary.Read(r, binary.LittleEndian, &m.StakeModifierChecksum)
	if e != nil {
		return e
	}
	e = binary.Read(r, binary.LittleEndian, &m.Flags)
	if e != nil {
		return e
	}
	e = binary.Read(r, binary.LittleEndian, &m.HashProofOfStake)
	if e != nil {
		return e
	}

	var blen byte
	e = binary.Read(r, binary.LittleEndian, &blen)
	if e != nil {
		return e
	}
	var arr = make([]byte, blen)
	e = binary.Read(r, binary.LittleEndian, &arr)
	if e != nil {
		return e
	}
	m.ChainTrust.SetBytes(arr)

	e = binary.Read(r, binary.LittleEndian, &m.Mint)
	if e != nil {
		return e
	}
	e = binary.Read(r, binary.LittleEndian, &m.MoneySupply)
	if e != nil {
		return e
	}

	var txOffsetCount uint32
	e = binary.Read(r, binary.LittleEndian, &txOffsetCount)
	if e != nil {
		return e
	}

	m.TxOffsets = make([]uint32, txOffsetCount)
	for i := uint32(0); i < txOffsetCount; i++ {
		e := binary.Read(r, binary.LittleEndian, &m.TxOffsets[i])
		if e != nil {
			return e
		}
	}

	return nil
}

// https://github.com/ppcoin/ppcoin/blob/v0.4.0ppc/src/main.h#L528
// ppc: the coin stake transaction is marked with the first output empty
func (msg *MsgTx) IsCoinStake() bool {
	return len(msg.TxIn) > 0 &&
		(!(msg.TxIn[0].PreviousOutPoint.Hash.IsEqual(&ShaHash{}) &&
			msg.TxIn[0].PreviousOutPoint.Index == MaxPrevOutIndex)) &&
		len(msg.TxOut) >= 2 &&
		(msg.TxOut[0].Value == 0 && len(msg.TxOut[0].PkScript) == 0)
}

// ppc:
func (t *TxOut) IsEmpty() bool {
	return t.Value == 0 && len(t.PkScript) == 0
}

// ppc: https://github.com/ppcoin/ppcoin/blob/master/src/main.h#L217
func (o *OutPoint) IsNull() bool {
	return o.Hash.IsEqual(&ZeroSha) && o.Index == MaxPrevOutIndex
}

// https://github.com/ppcoin/ppcoin/blob/v0.4.0ppc/src/main.h#L962
// ppc: two types of block: proof-of-work or proof-of-stake
func (msg *MsgBlock) IsProofOfStake() bool {
	return len(msg.Transactions) > 1 &&
		msg.Transactions[1].IsCoinStake()
}

func (m *Meta) GetSerializedSize() int {
	return 8 + // StakeModifier uint64
		4 + // StakeModifierChecksum uint32
		32 + // HashProofOfStake ShaHash
		4 + // Flags uint32
		1 + len(m.ChainTrust.Bytes()) + //ChainTrust big.Int
		8 + // Mint int64
		8 + // MoneySupply int64
		4 + // TxOffsets array size uint32
		4*len(m.TxOffsets) // TxOffsets uint32 array
}
