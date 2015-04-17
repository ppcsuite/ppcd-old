// Copyright (c) 2014-2014 PPCD developers.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ppcsuite/btcutil"
	"github.com/ppcsuite/ppcd/chaincfg"
)

const (
	// Protocol switch time of v0.3 kernel protocol
	nProtocolV03SwitchTime     int64 = 1363800000
	nProtocolV03TestSwitchTime int64 = 1359781000
	// Protocol switch time of v0.4 kernel protocol
	nProtocolV04SwitchTime     int64 = 1399300000
	nProtocolV04TestSwitchTime int64 = 1395700000
	// TxDB upgrade time for v0.4 protocol
	// Note: v0.4 upgrade does not require block chain re-download. However,
	//       user must upgrade before the protocol switch deadline, otherwise
	//       re-download of blockchain is required. The timestamp of upgrade
	//       is recorded in transaction database to alert user of the requirement.
	nProtocolV04UpgradeTime int64 = 0
)

func getBlockTrust(block *btcutil.Block) *big.Int {
	return calcTrust(block.MsgBlock().Header.Bits, block.MsgBlock().IsProofOfStake())
}

// ppcoin: entropy bit for stake modifier if chosen by modifier
func getStakeEntropyBit(b *BlockChain, block *btcutil.Block) (uint32, error) {

	defer timeTrack(now(), fmt.Sprintf("getStakeEntropyBit(%v)", slice(block.Sha())[0]))

	nEntropyBit := uint32(0)
	hash := block.Sha()

	if isProtocolV04(b, int64(block.MsgBlock().Header.Timestamp.Unix())) {

		nEntropyBit = uint32((ShaHashToBig(hash).Int64()) & 1) // last bit of block hash

		//if (fDebug && GetBoolArg("-printstakemodifier"))
		//    printf("GetStakeEntropyBit(v0.4+): nTime=%d hashBlock=%s entropybit=%d\n", nTime, GetHash().ToString().c_str(), nEntropyBit);

	} else {

		// old protocol for entropy bit pre v0.4
		hashSigBytes := btcutil.Hash160(block.MsgBlock().Signature)
		// to big-endian
		blen := len(hashSigBytes)
		for i := 0; i < blen/2; i++ {
			hashSigBytes[i], hashSigBytes[blen-1-i] = hashSigBytes[blen-1-i], hashSigBytes[i]
		}
		//if (fDebug && GetBoolArg("-printstakemodifier"))
		//    printf("GetStakeEntropyBit(v0.3): nTime=%d hashSig=%s", nTime, hashSig.ToString().c_str());
		hashSig := new(big.Int).SetBytes(hashSigBytes)
		hashSig.Rsh(hashSig, 159) // take the first bit of the hash
		nEntropyBit = uint32(hashSig.Int64())

		//if (fDebug && GetBoolArg("-printstakemodifier"))
		//    printf(" entropybit=%d\n", nEntropyBit)
	}

	log.Tracef("Entropy bit = %d for block %v", nEntropyBit, hash)

	return nEntropyBit, nil
}

func getStakeModifierHexString(stakeModifier uint64) string {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, stakeModifier)
	return hex.EncodeToString(bytes)
}

func getStakeModifierCSHexString(stakeModifierCS uint32) string {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, stakeModifierCS)
	return hex.EncodeToString(bytes)
}

// isProtocolV03
func isProtocolV03(b *BlockChain, nTime int64) bool {
	var switchTime int64
	if b.chainParams.Name == "testnet3" {
		switchTime = nProtocolV03TestSwitchTime
	} else {
		switchTime = nProtocolV03SwitchTime
	}
	return nTime >= switchTime
}

// isProtocolV03FromParams
func isProtocolV03FromParams(params *chaincfg.Params, nTime int64) bool {
	var switchTime int64
	if params.Name == "testnet3" {
		switchTime = nProtocolV03TestSwitchTime
	} else {
		switchTime = nProtocolV03SwitchTime
	}
	return nTime >= switchTime
}

// isProtocolV04
func isProtocolV04(b *BlockChain, nTime int64) bool {
	var v04SwitchTime int64
	if b.chainParams.Name == "testnet3" {
		v04SwitchTime = nProtocolV04TestSwitchTime
	} else {
		v04SwitchTime = nProtocolV04SwitchTime
	}
	return nTime >= v04SwitchTime
}

// dateTimeStrFormat displays time in RFC3339 format
func dateTimeStrFormat(t int64) string {
	return time.Unix(t, 0).UTC().Format(time.RFC3339)
}

func minInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func minInt64(a int64, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func now() time.Time {
	return btcutil.Now()
}

func timeTrack(start time.Time, name string) {
	btcutil.TimeTrack(log, start, name)
}

func slice(args ...interface{}) []interface{} {
	return btcutil.Slice(args)
}
