package main

import "github.com/mably/btcjson"

// makeGetKernelStakeModifier generates the cmd structure for getkernelstakemodifier commands.
func makeGetKernelStakeModifier(args []interface{}) (btcjson.Cmd, error) {
	// Create the getblock command with defaults for the optional
	// parameters.
	getKernelStakeModifierCmd, err :=
		btcjson.NewGetKernelStakeModifierCmd("btcctl", args[0].(string))
	if err != nil {
		return nil, err
	}

	// Override the optional parameters if they were specified.
	if len(args) > 1 {
		getKernelStakeModifierCmd.Verbose = args[1].(bool)
	}

	return getKernelStakeModifierCmd, nil
}
