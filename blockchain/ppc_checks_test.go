package blockchain_test

import (
	"bytes"
	"github.com/ppcsuite/ppcd/chaincfg"
	"github.com/ppcsuite/ppcd/blockchain"
	"github.com/ppcsuite/ppcd/wire"
	"testing"
)

func TestCheckBlockSignature(t *testing.T) {
	if !blockchain.CheckBlockSignature(&Block100000, &chaincfg.MainNetParams) {
		t.Error("bad block signature, valid expected")
	}
	var buf bytes.Buffer
	err := Block100000.Serialize(&buf)
	if err != nil {
		t.Error(err)
		return
	}
	rbuf := bytes.NewReader(buf.Bytes())
	block := new(wire.MsgBlock)
	err = block.Deserialize(rbuf)
	if err != nil {
		t.Error(err)
		return
	}
	block.Signature[5] ^= 0xff
	if blockchain.CheckBlockSignature(block, &chaincfg.MainNetParams) {
		t.Error("good block signature, invalid expected")
	}
}
