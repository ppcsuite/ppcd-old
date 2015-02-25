// Copyright (c) 2014-2014 PPCD developers.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"strconv"

	"github.com/ppcsuite/ppcd/blockchain"
	"github.com/ppcsuite/ppcd/btcjson/v2/btcjson"
	"github.com/ppcsuite/ppcd/database"
	"github.com/ppcsuite/ppcd/wire"
)

// getDifficultyRatio returns the latest PoW or PoS difficulty up to block sha.
func ppcGetDifficultyRatio(db database.Db, sha *wire.ShaHash, proofOfStake bool) (float64, error) {
	bh, _, err := blockchain.GetLastBlockHeader(db, sha, proofOfStake)
	if err != nil {
		return 0, err
	}
	return getDifficultyRatio(bh.Bits), nil
}

// ppcHandleGetDifficulty implements the getdifficulty command.
func ppcHandleGetDifficulty(s *rpcServer, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	sha, _, err := s.server.db.NewestSha()
	if err != nil {
		context := "Error getting sha"
		return nil, internalRPCError(err.Error(), context)
	}
	powDifficulty, err := ppcGetDifficultyRatio(s.server.db, sha, false) // ppc: PoW
	if err != nil {
		context := "Error getting difficulty"
		return nil, internalRPCError(err.Error(), context)
	}
	posDifficulty, err := ppcGetDifficultyRatio(s.server.db, sha, true) // ppc: PoS
	if err != nil {
		context := "Error getting difficulty"
		return nil, internalRPCError(err.Error(), context)
	}

	ret := &btcjson.GetDifficultyResult{
		ProofOfWork:    powDifficulty,
		ProofOfStake:   posDifficulty,
		SearchInterval: int32(0),
	}

	return ret, nil
}

// ppcHandleGetKernelStakeModifier implements the getkernelstakeModifier command.
func ppcHandleGetKernelStakeModifier(s *rpcServer, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	c := cmd.(*btcjson.GetKernelStakeModifierCmd)
	sha, err := wire.NewShaHashFromStr(c.Hash)
	if err != nil {
		context := "Error generating sha"
		return nil, internalRPCError(err.Error(), context)
	}

	chain := s.server.blockManager.blockChain
	kernelStakeModifier, err := chain.GetKernelStakeModifier(sha, s.server.timeSource)
	if err != nil {
		context := "Error getting kernel stake modifier for block " + sha.String()
		return nil, internalRPCError(err.Error(), context)
	}

	// When the verbose flag isn't set, simply return a string.
	if !*c.Verbose {
		return strconv.FormatUint(kernelStakeModifier, 10), nil
	}

	// The verbose flag is set, so generate the JSON object and return it.
	ksmReply := btcjson.KernelStakeModifierResult{
		Hash:                c.Hash,
		KernelStakeModifier: btcjson.StakeModifier(kernelStakeModifier),
	}

	return ksmReply, nil
}

// ppcHandleGetNextRequiredTarget implements the getNextRequiredTarget command.
func ppcHandleGetNextRequiredTarget(s *rpcServer, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	c := cmd.(*btcjson.GetNextRequiredTargetCmd)
	chain := s.server.blockManager.blockChain
	nextRequiredTarget, err := chain.PPCCalcNextRequiredDifficulty(*c.ProofOfStake)
	if err != nil {
		context := "Error getting next required target"
		return nil, internalRPCError(err.Error(), context)
	}

	// When the verbose flag isn't set, simply return a string.
	if !*c.Verbose {
		return strconv.FormatUint(uint64(nextRequiredTarget), 10), nil
	}

	// The verbose flag is set, so generate the JSON object and return it.
	ksmReply := btcjson.NextRequiredTargetResult{
		Target: btcjson.RequiredTarget(nextRequiredTarget),
	}

	return ksmReply, nil
}
