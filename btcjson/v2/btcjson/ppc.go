// Copyright (c) 2014-2014 PPCD developers.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package btcjson

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// FloatAmount specific type with custom marshalling
type FloatAmount float64

// MarshalJSON provides a custom Marshal method for FloatAmount.
func (v *FloatAmount) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.6f", *v)), nil
}

// GetDifficultyResult models the data of getdifficulty command.
type GetDifficultyResult struct {
	ProofOfWork    float64 `json:"proof-of-work"`
	ProofOfStake   float64 `json:"proof-of-stake"`
	SearchInterval int32   `json:"search-interval"`
}

// GetKernelStakeModifierCmd is a type handling custom marshaling and
// unmarshaling of getkernelstakemodifier JSON RPC commands.
type GetKernelStakeModifierCmd struct {
	Hash    string
	Verbose *bool `jsonrpcdefault:"true"`
}

// NewGetKernelStakeModifierCmd creates a new GetKernelStakeModifierCmd.
func NewGetKernelStakeModifierCmd(hash string, verbose *bool) *GetKernelStakeModifierCmd {
	return &GetKernelStakeModifierCmd{
		Hash:    hash,
		Verbose: verbose,
	}
}

// StakeModifier specific type with custom marshalling
type StakeModifier uint64

// MarshalJSON provides a custom Marshal method for StakeModifier.
func (v StakeModifier) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%d\"", uint64(v))), nil
}

// UnmarshalJSON provides a custom Unmarshal method for StakeModifier.
func (v *StakeModifier) UnmarshalJSON(b []byte) (err error) {
	var s string
	json.Unmarshal(b, &s)
	var u uint64
	u, err = strconv.ParseUint(s, 0, 64)
	*v = StakeModifier(u)
	return
}

// KernelStakeModifierResult models the data from the getkernelstakemodifier
// command when the verbose flag is set.  When the verbose flag is not set,
// getkernelstakemodifier return a hex-encoded string.
type KernelStakeModifierResult struct {
	Hash                string        `json:"hash"`
	KernelStakeModifier StakeModifier `json:"kernelstakemodifier"`
}

// GetNextRequiredTargetCmd is a type handling custom marshaling and
// unmarshaling of getNextRequiredTarget JSON RPC commands.
type GetNextRequiredTargetCmd struct {
	ProofOfStake *bool `jsonrpcdefault:"false"`
	Verbose      *bool `jsonrpcdefault:"true"`
}

// NewGetNextRequiredTargetCmd creates a new GetNextRequiredTargetCmd.
func NewGetNextRequiredTargetCmd(proofOfStake *bool, verbose *bool) *GetNextRequiredTargetCmd {
	return &GetNextRequiredTargetCmd{
		ProofOfStake: proofOfStake,
		Verbose:      verbose,
	}
}

// StakeModifier specific type with custom marshalling
type RequiredTarget uint32

// MarshalJSON provides a custom Marshal method for StakeModifier.
func (v RequiredTarget) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%d\"", uint32(v))), nil
}

// UnmarshalJSON provides a custom Unmarshal method for StakeModifier.
func (v *RequiredTarget) UnmarshalJSON(b []byte) (err error) {
	var s string
	json.Unmarshal(b, &s)
	var u uint64
	u, err = strconv.ParseUint(s, 0, 32)
	*v = RequiredTarget(u)
	return
}

// NextRequiredTargetResult models the data from the getnextrequiredtarget
// command when the verbose flag is set.  When the verbose flag is not set,
// getnextrequiredtarget return a simple string.
type NextRequiredTargetResult struct {
	Target RequiredTarget `json:"target"`
}

// GetLastProofOfWorkRewardCmd is a type handling custom marshaling and
// unmarshaling of getLastProofOfWorkReward JSON RPC commands.
type GetLastProofOfWorkRewardCmd struct {
}

// NewGetNextRequiredTargetCmd creates a new GetNextRequiredTargetCmd.
func NewGetLastProofOfWorkRewardCmd() *GetLastProofOfWorkRewardCmd {
	return &GetLastProofOfWorkRewardCmd{}
}

// StakeModifier specific type with custom marshalling
type BlockReward int64

// MarshalJSON provides a custom Marshal method for BlockReward.
func (v BlockReward) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%d\"", int64(v))), nil
}

// UnmarshalJSON provides a custom Unmarshal method for BlockReward.
func (v *BlockReward) UnmarshalJSON(b []byte) (err error) {
	var s string
	json.Unmarshal(b, &s)
	var u int64
	u, err = strconv.ParseInt(s, 0, 64)
	*v = BlockReward(u)
	return
}

// LastProofOfWorkRewardResult models the data from the getlastproofofworkreward
// command.
type LastProofOfWorkRewardResult struct {
	Subsidy BlockReward `json:"subsidy"`
}

// FindStakeCmd is a type handling custom marshaling and
// unmarshaling of FindStake JSON RPC commands.
type FindStakeCmd struct {
	MaxTime    int64
	Difficulty float64
	Verbose    *bool `jsonrpcdefault:"true"`
}

// NewFindStakeCmd creates a new FindStakeCmd.
func NewFindStakeCmd(maxTime int64, difficulty float64, verbose *bool) *FindStakeCmd {
	return &FindStakeCmd{
		MaxTime:    maxTime,
		Difficulty: difficulty,
		Verbose:    verbose,
	}
}

// FindStakeResult models the data from the FindStake
// command when the verbose flag is set.  When the verbose flag is not set,
// FindStake return a simple string.
type FindStakeResult struct {
	Difficulty float64 `json:"difficulty"`
	Time       int64   `json:"time"`
}

func init() {
	// No special flags for commands in this file.
	flags := UsageFlag(0)

	MustRegisterCmd("getkernelstakemodifier", (*GetKernelStakeModifierCmd)(nil), flags)
	MustRegisterCmd("getnextrequiredtarget", (*GetNextRequiredTargetCmd)(nil), flags)
	MustRegisterCmd("getlastproofofworkreward", (*GetLastProofOfWorkRewardCmd)(nil), flags)

	MustRegisterCmd("findstake", (*FindStakeCmd)(nil), UFWalletOnly)
}
