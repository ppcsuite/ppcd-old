package main

import (
	"github.com/mably/btcjson"
	"github.com/mably/btcws"
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
		btcws.NewFindStakeCmd("btcctl", args[0].(int64))
	if err != nil {
		return nil, err
	}

	// Override the optional parameters if they were specified.
	if len(args) > 1 {
		findStakeCmd.Verbose = args[1].(bool)
	}

	return findStakeCmd, nil
}
