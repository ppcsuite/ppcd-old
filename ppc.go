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

	kernelStakeModifier, err := s.server.blockManager.GetKernelStakeModifier(sha)
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
	nextRequiredTarget, err := s.server.blockManager.PPCCalcNextRequiredDifficulty(*c.ProofOfStake)
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

// ppcHandleGetLastProofOfWorkReward implements the getLastProofOfWorkReward command.
func ppcHandleGetLastProofOfWorkReward(s *rpcServer, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	//c := cmd.(*btcjson.GetLastProofOfWorkRewardCmd)
	reward, err := s.server.blockManager.PPCGetLastProofOfWorkReward()
	if err != nil {
		context := "Error getting last proof of work reward"
		return nil, internalRPCError(err.Error(), context)
	}

	// When the verbose flag isn't set, simply return a string.
	//return strconv.FormatInt(int64(reward), 10), nil

	// The verbose flag is set, so generate the JSON object and return it.
	lpowrReply := btcjson.LastProofOfWorkRewardResult{
		Subsidy: btcjson.BlockReward(reward),
	}

	return lpowrReply, nil
}

// ppcCalcNextReqDifficultyResponse is a response sent to the reply channel of a
// ppcCalcNextReqDifficultyMsg query.
type ppcCalcNextReqDifficultyResponse struct {
	difficulty uint32
	err        error
}

// ppcCalcNextReqDifficultyMsg is a message type to be sent across the message
// channel for requesting the required difficulty of the next block.
type ppcCalcNextReqDifficultyMsg struct {
	proofOfStake bool
	reply        chan ppcCalcNextReqDifficultyResponse
}

// PPCCalcNextRequiredDifficulty calculates the required difficulty for the block
// after the end of the current best chain based on the difficulty retarget
// rules.
func (b *blockManager) PPCCalcNextRequiredDifficulty(proofOfStake bool) (uint32, error) {
	reply := make(chan ppcCalcNextReqDifficultyResponse, 1)
	b.msgChan <- ppcCalcNextReqDifficultyMsg{proofOfStake: proofOfStake, reply: reply}
	response := <-reply
	return response.difficulty, response.err
}

// getKernelStakeModifierResponse is a response sent to the reply channel of a
// getKernelStakeModifierMsg.
type getKernelStakeModifierResponse struct {
	StakeModifier uint64
	err           error
}

// getKernelStakeModifierMsg is a message type to be sent across the message
// channel getting the stake modifier for some block.
type getKernelStakeModifierMsg struct {
	hash       *wire.ShaHash
	timeSource blockchain.MedianTimeSource
	reply      chan getKernelStakeModifierResponse
}

// GetKernelStakeModifier TODO(mably)
func (b *blockManager) GetKernelStakeModifier(hash *wire.ShaHash) (uint64, error) {
	reply := make(chan getKernelStakeModifierResponse, 1)
	b.msgChan <- getKernelStakeModifierMsg{hash: hash, timeSource: b.server.timeSource, reply: reply}
	response := <-reply
	return response.StakeModifier, response.err
}

// ppcGetLastProofOfWorkRewardResponse is a response sent to the reply channel of a
// ppcGetLastProofOfWorkRewardMsg query.
type ppcGetLastProofOfWorkRewardResponse struct {
	subsidy int64
	err     error
}

// ppcGetLastProofOfWorkRewardMsg is a message type to be sent across the message
// channel for requesting the last proof of work reward.
type ppcGetLastProofOfWorkRewardMsg struct {
	reply chan ppcGetLastProofOfWorkRewardResponse
}

// PPCGetLastProofOfWorkReward calculates the last proof of work reward.
func (b *blockManager) PPCGetLastProofOfWorkReward() (int64, error) {
	reply := make(chan ppcGetLastProofOfWorkRewardResponse, 1)
	b.msgChan <- ppcGetLastProofOfWorkRewardMsg{reply: reply}
	response := <-reply
	return response.subsidy, response.err
}
