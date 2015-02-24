// Copyright (c) 2014-2014 PPCD developers.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package btcws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ppcsuite/ppcd/btcjson"
)

// GetKernelStakeModifierCmd is a type handling custom marshaling and
// unmarshaling of getkernelstakemodifier JSON RPC commands.
type GetKernelStakeModifierCmd struct {
	id      interface{}
	Hash    string
	Verbose bool
}

// Enforce that GetKernelStakeModifierCmd satisifies the Cmd interface.
var _ btcjson.Cmd = &GetKernelStakeModifierCmd{}

// NewGetKernelStakeModifierCmd creates a new GetKernelStakeModifierCmd.
func NewGetKernelStakeModifierCmd(id interface{}, hash string, optArgs ...bool) (*GetKernelStakeModifierCmd, error) {
	// default verbose is set to true to match old behavior
	verbose, verboseTx := true, false

	optArgsLen := len(optArgs)
	if optArgsLen > 0 {
		if optArgsLen > 2 {
			return nil, btcjson.ErrTooManyOptArgs
		}
		verbose = optArgs[0]
		if optArgsLen > 1 {
			verboseTx = optArgs[1]

			if !verbose && verboseTx {
				return nil, btcjson.ErrInvalidParams
			}
		}
	}

	return &GetKernelStakeModifierCmd{
		id:      id,
		Hash:    hash,
		Verbose: verbose,
	}, nil
}

// Id satisfies the Cmd interface by returning the id of the command.
func (cmd *GetKernelStakeModifierCmd) Id() interface{} {
	return cmd.id
}

// Method satisfies the Cmd interface by returning the json method.
func (cmd *GetKernelStakeModifierCmd) Method() string {
	return "getkernelstakemodifier"
}

// parseGetKernelStakeModifierCmd parses a RawCmd into a concrete type satisifying
// the btcjson.Cmd interface.  This is used when registering the custom
// command with the btcjson parser.
func parseGetKernelStakeModifierCmd(r *btcjson.RawCmd) (btcjson.Cmd, error) {

	if len(r.Params) > 2 || len(r.Params) < 1 {
		return nil, btcjson.ErrWrongNumberOfParams
	}

	var hash string
	if err := json.Unmarshal(r.Params[0], &hash); err != nil {
		return nil, fmt.Errorf("parameter 'hash' must be a string: %v", err)
	}

	optArgs := make([]bool, 0, 2)
	if len(r.Params) > 1 {
		var verbose bool
		if err := json.Unmarshal(r.Params[1], &verbose); err != nil {
			return nil, fmt.Errorf("second optional parameter 'verbose' must be a bool: %v", err)
		}
		optArgs = append(optArgs, verbose)
	}

	newCmd, err := NewGetKernelStakeModifierCmd(r.Id, hash, optArgs...)
	if err != nil {
		return nil, err
	}

	return newCmd, nil
}

// parseGetKernelStakeModifierCmdReply parses a the reply to a
// GetKernelStakeModifierCmd into a concrete type and returns it packed into
// an interface.  This is used when registering the custom command with btcjson.
func parseGetKernelStakeModifierCmdReply(message json.RawMessage) (interface{}, error) {
	if bytes.IndexByte(message, '{') > -1 {
		var res KernelStakeModifierResult
		if err := json.Unmarshal(message, &res); err != nil {
			return nil, err
		}
		return res, nil
	} else {
		var res string
		if err := json.Unmarshal(message, &res); err != nil {
			return nil, err
		}
		return res, nil
	}
}

// MarshalJSON returns the JSON encoding of cmd.  Part of the Cmd interface.
func (cmd *GetKernelStakeModifierCmd) MarshalJSON() ([]byte, error) {
	params := make([]interface{}, 1, 3)
	params[0] = cmd.Hash
	if !cmd.Verbose {
		// set optional verbose argument to false
		params = append(params, false)
	}
	// Fill and marshal a RawCmd.
	raw, err := btcjson.NewRawCmd(cmd.id, cmd.Method(), params)
	if err != nil {
		return nil, err
	}
	return json.Marshal(raw)
}

// UnmarshalJSON unmarshals the JSON encoding of cmd into cmd.  Part of
// the Cmd interface.
func (cmd *GetKernelStakeModifierCmd) UnmarshalJSON(b []byte) (err error) {
	// Unmashal into a RawCmd
	var r btcjson.RawCmd
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}

	var parsedCmd btcjson.Cmd
	parsedCmd, err = parseGetKernelStakeModifierCmd(&r)

	*cmd = *parsedCmd.(*GetKernelStakeModifierCmd)

	return nil
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
	id           interface{}
	ProofOfStake bool
	Verbose      bool
}

// Enforce that GetNextRequiredTargetCmd satisifies the Cmd interface.
var _ btcjson.Cmd = &GetNextRequiredTargetCmd{}

// NewGetNextRequiredTargetCmd creates a new GetNextRequiredTargetCmd.
func NewGetNextRequiredTargetCmd(id interface{}, proofOfStake bool, optArgs ...bool) (*GetNextRequiredTargetCmd, error) {
	// default verbose is set to true to match old behavior
	verbose := true

	optArgsLen := len(optArgs)
	if optArgsLen > 0 {
		if optArgsLen > 1 {
			return nil, btcjson.ErrTooManyOptArgs
		}
		verbose = optArgs[0]
	}

	return &GetNextRequiredTargetCmd{
		id:           id,
		ProofOfStake: proofOfStake,
		Verbose:      verbose,
	}, nil
}

// Id satisfies the Cmd interface by returning the id of the command.
func (cmd *GetNextRequiredTargetCmd) Id() interface{} {
	return cmd.id
}

// Method satisfies the Cmd interface by returning the json method.
func (cmd *GetNextRequiredTargetCmd) Method() string {
	return "getnextrequiredtarget"
}

// parseGetNextRequiredTargetCmd parses a RawCmd into a concrete type satisifying
// the btcjson.Cmd interface.  This is used when registering the custom
// command with the btcjson parser.
func parseGetNextRequiredTargetCmd(r *btcjson.RawCmd) (btcjson.Cmd, error) {
	if len(r.Params) > 2 || len(r.Params) < 1 {
		return nil, btcjson.ErrWrongNumberOfParams
	}

	var proofOfStake bool
	if err := json.Unmarshal(r.Params[0], &proofOfStake); err != nil {
		return nil, fmt.Errorf("parameter 'proofofstake' must be a bool: %v", err)
	}

	optArgs := make([]bool, 0, 2)
	if len(r.Params) > 1 {
		var verbose bool
		if err := json.Unmarshal(r.Params[1], &verbose); err != nil {
			return nil, fmt.Errorf("second optional parameter 'verbose' must be a bool: %v", err)
		}
		optArgs = append(optArgs, verbose)
	}

	newCmd, err := NewGetNextRequiredTargetCmd(r.Id, proofOfStake, optArgs...)
	if err != nil {
		return nil, err
	}

	return newCmd, nil
}

// parseGetNextRequiredTargetCmdReply parses a the reply to a
// GetNextRequiredTargetCmd into a concrete type and returns it packed into
// an interface.  This is used when registering the custom command with btcjson.
func parseGetNextRequiredTargetCmdReply(message json.RawMessage) (interface{}, error) {
	if bytes.IndexByte(message, '{') > -1 {
		var res NextRequiredTargetResult
		if err := json.Unmarshal(message, &res); err != nil {
			return nil, err
		}
		return res, nil
	} else {
		var res string
		if err := json.Unmarshal(message, &res); err != nil {
			return nil, err
		}
		return res, nil
	}
}

// MarshalJSON returns the JSON encoding of cmd.  Part of the Cmd interface.
func (cmd *GetNextRequiredTargetCmd) MarshalJSON() ([]byte, error) {
	params := make([]interface{}, 1, 3)
	params[0] = cmd.ProofOfStake
	if !cmd.Verbose {
		// set optional verbose argument to false
		params = append(params, false)
	}
	// Fill and marshal a RawCmd.
	raw, err := btcjson.NewRawCmd(cmd.id, cmd.Method(), params)
	if err != nil {
		return nil, err
	}
	return json.Marshal(raw)
}

// UnmarshalJSON unmarshals the JSON encoding of cmd into cmd.  Part of
// the Cmd interface.
func (cmd *GetNextRequiredTargetCmd) UnmarshalJSON(b []byte) (err error) {
	// Unmashal into a RawCmd
	var r btcjson.RawCmd
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}

	var parsedCmd btcjson.Cmd
	parsedCmd, err = parseGetNextRequiredTargetCmd(&r)

	*cmd = *parsedCmd.(*GetNextRequiredTargetCmd)

	return nil
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

// FindStakeCmd is a type handling custom marshaling and
// unmarshaling of FindStake JSON RPC commands.
type FindStakeCmd struct {
	id         interface{}
	MaxTime    int64
	Difficulty float64
	Verbose    bool
}

// Enforce that FindStakeCmd satisifies the Cmd interface.
var _ btcjson.Cmd = &FindStakeCmd{}

// NewFindStakeCmd creates a new FindStakeCmd.
func NewFindStakeCmd(id interface{}, maxTime int64, difficulty float64, optArgs ...bool) (*FindStakeCmd, error) {
	// default verbose is set to true to match old behavior
	verbose := true

	optArgsLen := len(optArgs)
	if optArgsLen > 0 {
		if optArgsLen > 1 {
			return nil, btcjson.ErrTooManyOptArgs
		}
		verbose = optArgs[0]
	}

	return &FindStakeCmd{
		id:         id,
		MaxTime:    maxTime,
		Difficulty: difficulty,
		Verbose:    verbose,
	}, nil
}

// Id satisfies the Cmd interface by returning the id of the command.
func (cmd *FindStakeCmd) Id() interface{} {
	return cmd.id
}

// Method satisfies the Cmd interface by returning the json method.
func (cmd *FindStakeCmd) Method() string {
	return "findstake"
}

// parseFindStakeCmd parses a RawCmd into a concrete type satisifying
// the btcjson.Cmd interface.  This is used when registering the custom
// command with the btcjson parser.
func parseFindStakeCmd(r *btcjson.RawCmd) (btcjson.Cmd, error) {

	if len(r.Params) > 3 || len(r.Params) < 1 {
		return nil, btcjson.ErrWrongNumberOfParams
	}

	var maxTime int64
	if err := json.Unmarshal(r.Params[0], &maxTime); err != nil {
		return nil, fmt.Errorf("parameter 'maxtime' must be a integer: %v", err)
	}

	var difficulty float64
	if err := json.Unmarshal(r.Params[1], &difficulty); err != nil {
		return nil, fmt.Errorf("parameter 'difficulty' must be a float: %v", err)
	}

	optArgs := make([]bool, 0, 2)
	if len(r.Params) > 2 {
		var verbose bool
		if err := json.Unmarshal(r.Params[2], &verbose); err != nil {
			return nil, fmt.Errorf("second optional parameter 'verbose' must be a bool: %v", err)
		}
		optArgs = append(optArgs, verbose)
	}

	newCmd, err := NewFindStakeCmd(r.Id, maxTime, difficulty, optArgs...)
	if err != nil {
		return nil, err
	}

	return newCmd, nil
}

// parseFindStakeCmdReply parses a the reply to a
// FindStakeCmd into a concrete type and returns it packed into
// an interface.  This is used when registering the custom command with btcjson.
func parseFindStakeCmdReply(message json.RawMessage) (interface{}, error) {
	if bytes.IndexByte(message, '[') > -1 {
		res := make([]FindStakeResult, 10)
		if err := json.Unmarshal(message, &res); err != nil {
			return nil, err
		}
		return res, nil
	} else if bytes.IndexByte(message, '{') > -1 {
		var res FindStakeResult
		if err := json.Unmarshal(message, &res); err != nil {
			return nil, err
		}
		return res, nil
	} else {
		var res string
		if err := json.Unmarshal(message, &res); err != nil {
			return nil, err
		}
		return res, nil
	}
}

// MarshalJSON returns the JSON encoding of cmd.  Part of the Cmd interface.
func (cmd *FindStakeCmd) MarshalJSON() ([]byte, error) {
	params := make([]interface{}, 2, 3)
	params[0] = cmd.MaxTime
	params[1] = cmd.Difficulty
	if !cmd.Verbose {
		// set optional verbose argument to false
		params = append(params, false)
	}
	// Fill and marshal a RawCmd.
	raw, err := btcjson.NewRawCmd(cmd.id, cmd.Method(), params)
	if err != nil {
		return nil, err
	}
	return json.Marshal(raw)
}

// UnmarshalJSON unmarshals the JSON encoding of cmd into cmd.  Part of
// the Cmd interface.
func (cmd *FindStakeCmd) UnmarshalJSON(b []byte) (err error) {
	// Unmashal into a RawCmd
	var r btcjson.RawCmd
	if err := json.Unmarshal(b, &r); err != nil {
		return err
	}

	var parsedCmd btcjson.Cmd
	parsedCmd, err = parseFindStakeCmd(&r)

	*cmd = *parsedCmd.(*FindStakeCmd)

	return nil
}

// FindStakeResult models the data from the FindStake
// command when the verbose flag is set.  When the verbose flag is not set,
// FindStake return a simple string.
type FindStakeResult struct {
	Difficulty float64 `json:"difficulty"`
	Time       int64   `json:"time"`
}
