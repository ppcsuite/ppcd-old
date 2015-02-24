// Copyright (c) 2014-2014 PPCD developers.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package btcjson

import "fmt"

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
