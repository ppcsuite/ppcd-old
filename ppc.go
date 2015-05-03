// Copyright (c) 2014-2014 PPCD developers.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"sync"

	"github.com/ppcsuite/btcutil"
	"github.com/ppcsuite/ppcd/blockchain"
	"github.com/ppcsuite/ppcd/btcjson"
	"github.com/ppcsuite/ppcd/database"
	"github.com/ppcsuite/ppcd/wire"
)

// mintState houses state that is used in between multiple RPC invocations to
// mintblock.
type mintState struct {
	sync.Mutex
	blockInfo map[wire.ShaHash]*BlockTemplate
}

// newMintState returns a new instance of a mintState with all internal fields
// initialized and ready to use.
func newMintState() *mintState {
	return &mintState{
		blockInfo: make(map[wire.ShaHash]*BlockTemplate),
	}
}

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

// ppcHandleSendCoinStakeTransaction implements the sendCoinStakeTransaction command.
func ppcHandleSendCoinStakeTransaction(s *rpcServer, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	c := cmd.(*btcjson.SendCoinStakeTransactionCmd)
	hexStr := c.HexTx
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	serializedTx, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, rpcDecodeHexError(hexStr)
	}
	msgtx := wire.NewMsgTx()
	err = msgtx.Deserialize(bytes.NewReader(serializedTx))
	if err != nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCDeserialization,
			Message: "TX decode failed: " + err.Error(),
		}
	}

	blockTemplate := s.server.cpuMiner.BuildMintBlock(msgtx)

	scstrReply := btcjson.SendCoinStakeTransactionResult{}

	if blockTemplate != nil {

		block := blockTemplate.block

		// Protect concurrent access from multiple RPC invocations for mint
		// requests and submission.
		s.mintState.Lock()
		defer s.mintState.Unlock()

		s.mintState.blockInfo[msgtx.TxSha()] = blockTemplate

		scstrReply.HexBlockSha = block.BlockSha().String()
	}

	return scstrReply, nil
}

// ppcHandleSendMintBlockSignature implements the sendMintBlockSignature command.
func ppcHandleSendMintBlockSignature(s *rpcServer, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	c := cmd.(*btcjson.SendMintBlockSignatureCmd)

	s.mintState.Lock()
	defer s.mintState.Unlock()

	txSha, _ := wire.NewShaHashFromStr(c.HexTx)

	if blockTemplate, ok := s.mintState.blockInfo[*txSha]; ok {

		block := blockTemplate.block
		block.Signature, _ = hex.DecodeString(c.HexSignature)

		s.server.cpuMiner.SubmitMintBlock(block)

	}

	scstrReply := btcjson.SendMintBlockSignatureResult{}

	return scstrReply, nil
}

// BuildMintBlock
func (m *CPUMiner) BuildMintBlock(coinStakeTx *wire.MsgTx) *BlockTemplate {

	minrLog.Infof("MintBlock: coinstaketx = %v", coinStakeTx.TxSha().String())

	// No point in searching for a solution before the chain is
	// synced.  Also, grab the same lock as used for block
	// submission, since the current block will be changing and
	// this would otherwise end up building a new block template on
	// a block that is in the process of becoming stale.
	m.submitBlockLock.Lock()
	_, curHeight := m.server.blockManager.chainState.Best()
	if curHeight != 0 && !m.server.blockManager.IsCurrent() {
		m.submitBlockLock.Unlock()
		return nil
	}

	// Create a new block template using the available transactions
	// in the memory pool as a source of transactions to potentially
	// include in the block.
	template, err := NewBlockTemplate(
		m.server.txMemPool, nil, btcutil.NewTx(coinStakeTx))
	m.submitBlockLock.Unlock()
	if err != nil {
		errStr := fmt.Sprintf("Failed to create new block "+
			"template: %v", err)
		minrLog.Errorf(errStr)
		return nil
	}

	return template
}

// SubmitMintBlock
func (m *CPUMiner) SubmitMintBlock(mintBlock *wire.MsgBlock) bool {

	block := btcutil.NewBlock(mintBlock)

	return m.submitBlock(block)
}
