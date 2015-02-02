// Copyright (c) 2014-2014 PPCD developers.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"strconv"
	"github.com/mably/ppcd/blockchain"
	"github.com/mably/ppcd/database"
	"github.com/mably/btcjson"
	"github.com/mably/btcwire"
	"github.com/mably/btcws"
)

// getDifficultyRatio returns the latest PoW or PoS difficulty up to block sha.
func ppcGetDifficultyRatio(db database.Db, sha *btcwire.ShaHash, proofOfStake bool) (float64, error) {
	bh, _, err := blockchain.GetLastBlockHeader(db, sha, proofOfStake)
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
	c := cmd.(*btcws.GetKernelStakeModifierCmd)
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

	// When the verbose flag isn't set, simply return a string.
	if !c.Verbose {
		return strconv.FormatUint(kernelStakeModifier, 10), nil
	}

	// The verbose flag is set, so generate the JSON object and return it.
	ksmReply := btcws.KernelStakeModifierResult{
		Hash:                c.Hash,
		KernelStakeModifier: btcws.StakeModifier(kernelStakeModifier),
	}

	return ksmReply, nil
}

// ppcHandleGetNextRequiredTarget implements the getNextRequiredTarget command.
func ppcHandleGetNextRequiredTarget(s *rpcServer, cmd btcjson.Cmd, closeChan <-chan struct{}) (interface{}, error) {
	c := cmd.(*btcws.GetNextRequiredTargetCmd)
	chain := s.server.blockManager.blockChain
	nextRequiredTarget, err := chain.PPCCalcNextRequiredDifficulty(c.ProofOfStake)
	if err != nil {
		rpcsLog.Errorf("Error getting next required target : %v", err)
		return nil, btcjson.ErrDifficulty
	}

	// When the verbose flag isn't set, simply return a string.
	if !c.Verbose {
		return strconv.FormatUint(uint64(nextRequiredTarget), 10), nil
	}

	// The verbose flag is set, so generate the JSON object and return it.
	ksmReply := btcws.NextRequiredTargetResult{
		Target: btcws.RequiredTarget(nextRequiredTarget),
	}

	return ksmReply, nil
}
