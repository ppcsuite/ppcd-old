// Copyright (c) 2014-2014 PPCD developers.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"github.com/mably/btcchain"
	"github.com/mably/btcdb"
	"github.com/mably/btcjson"
	"github.com/mably/btcwire"
)

// getDifficultyRatio returns the latest PoW or PoS difficulty up to block sha.
func ppcGetDifficultyRatio(db btcdb.Db, sha *btcwire.ShaHash, proofOfStake bool) (float64, error) {
	bh, _, err := btcchain.GetLastBlockHeader(db, sha, proofOfStake)
	if err != nil {
		return 0, err
	}
	return getDifficultyRatio(bh.Bits), nil
}

// ppcHandleGetDifficulty implements the getdifficulty command.
func ppcHandleGetDifficulty(s *rpcServer, cmd btcjson.Cmd, closeChan <-chan struct{}) (interface{}, error) {
	sha, _, err := s.server.db.NewestSha()
	if err != nil {
		rpcsLog.Errorf("Error getting sha: %v", err)
		return nil, btcjson.ErrDifficulty
	}
	powDifficulty, err := ppcGetDifficultyRatio(s.server.db, sha, false) // ppc: PoW
	if err != nil {
		rpcsLog.Errorf("Error getting difficulty: %v", err)
		return nil, btcjson.ErrDifficulty
	}
	posDifficulty, err := ppcGetDifficultyRatio(s.server.db, sha, true) // ppc: PoS
	if err != nil {
		rpcsLog.Errorf("Error getting difficulty: %v", err)
		return nil, btcjson.ErrDifficulty
	}

	ret := &btcjson.GetDifficultyResult{
		ProofOfWork:    powDifficulty,
		ProofOfStake:   posDifficulty,
		SearchInterval: int32(0),
	}

	return ret, nil
}

// ppcHandleGetKernelStakeModifier implements the getkernelstakeModifier command.
func ppcHandleGetKernelStakeModifier(s *rpcServer, cmd btcjson.Cmd, closeChan <-chan struct{}) (interface{}, error) {
	c := cmd.(*btcjson.GetKernelStakeModifierCmd)
	sha, err := btcwire.NewShaHashFromStr(c.Hash)
	if err != nil {
		rpcsLog.Errorf("Error generating sha: %v", err)
		return nil, btcjson.ErrBlockNotFound
	}

	chain := s.server.blockManager.blockChain
	kernelStakeModifier, err := chain.GetKernelStakeModifier(sha, s.server.timeSource)
	if err != nil {
		rpcsLog.Errorf("Error getting kernel stake modifier for block %v : %v", sha, err)
		return nil, btcjson.ErrBlockNotFound
	}

	// When the verbose flag isn't set, simply return the network-serialized
	// block as a hex-encoded string.
	if !c.Verbose {
		return kernelStakeModifier, nil
	}

	// The verbose flag is set, so generate the JSON object and return it.
	ksmReply := btcjson.KernelStakeModifierResult{
		Hash:                c.Hash,
		KernelStakeModifier: kernelStakeModifier,
	}

	return ksmReply, nil
}
