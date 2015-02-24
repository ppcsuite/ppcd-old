package main

import (
	"strconv"

	"github.com/ppcsuite/ppcd/btcjson"
	"github.com/ppcsuite/ppcd/btcjson/btcws"
)

// makeGetKernelStakeModifier generates the cmd structure for getkernelstakemodifier commands.
func makeGetKernelStakeModifier(args []interface{}) (btcjson.Cmd, error) {
	// Create the getblock command with defaults for the optional
	// parameters.
	getKernelStakeModifierCmd, err :=
		btcws.NewGetKernelStakeModifierCmd("btcctl", args[0].(string))
	if err != nil {
		return nil, err
	}

	// Override the optional parameters if they were specified.
	if len(args) > 1 {
		getKernelStakeModifierCmd.Verbose = args[1].(bool)
	}

	return getKernelStakeModifierCmd, nil
}

// makeGetNextRequiredTarget generates the cmd structure for getNextRequiredTarget commands.
func makeGetNextRequiredTarget(args []interface{}) (btcjson.Cmd, error) {
	// Create the getblock command with defaults for the optional
	// parameters.
	getNextRequiredTargetCmd, err :=
		btcws.NewGetNextRequiredTargetCmd("btcctl", args[0].(bool))
	if err != nil {
		return nil, err
	}

	// Override the optional parameters if they were specified.
	if len(args) > 1 {
		getNextRequiredTargetCmd.Verbose = args[1].(bool)
	}

	return getNextRequiredTargetCmd, nil
}

// makeFindStake generates the cmd structure for findstake commands.
func makeFindStake(args []interface{}) (btcjson.Cmd, error) {
	// Create the getblock command with defaults for the optional
	// parameters.
	findStakeCmd, err :=
		btcws.NewFindStakeCmd("btcctl", args[0].(int64), args[1].(float64))
	if err != nil {
		return nil, err
	}

	// Override the optional parameters if they were specified.
	if len(args) > 2 {
		findStakeCmd.Verbose = args[2].(bool)
	}

	return findStakeCmd, nil
}

// toFloat64 attempts to convert the passed string to a float64.  It returns the
// float packed into an interface so it can be used in the calls which expect
// interfaces.  An error will be returned if the string can't be converted to a
// float.
func toFloat64(val string) (interface{}, error) {
	idx, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return nil, err
	}

	return idx, nil
}
