chaincfg
========

[![Build Status](https://travis-ci.org/ppcsuite/ppcd.png?branch=master)]
(https://travis-ci.org/ppcsuite/ppcd) [![Coverage Status]
(https://coveralls.io/repos/ppcsuite/ppcd/badge.png?branch=master)]
(https://coveralls.io/r/ppcsuite/ppcd?branch=master)
[![tip for next commit](http://peer4commit.com/projects/130.svg)](http://peer4commit.com/projects/130)

Package chaincfg defines chain configuration parameters for the three standard
Bitcoin networks and provides the ability for callers to define their own custom
Bitcoin networks.

Although this package was primarily written for btcd, it has intentionally been
designed so it can be used as a standalone package for any projects needing to
use parameters for the standard Bitcoin networks or for projects needing to
define their own network.

## Sample Use

```Go
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ppcsuite/ppcd/chaincfg"
	"github.com/mably/btcutil"
)

var testnet = flag.Bool("testnet", false, "operate on the testnet Bitcoin network")

// By default (without -testnet), use mainnet.
var netParams = &chaincfg.MainNetParams

func main() {
	flag.Parse()

	// Modify active network parameters if operating on testnet.
	if *testnet {
		netParams = &chaincfg.TestNet3Params
	}

	// later...

	// Create and print new payment address, specific to the active network.
	pubKeyHash := make([]byte, 20)
	addr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, netParams)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(addr)
}
```

## Documentation

[![GoDoc](https://godoc.org/github.com/ppcsuite/ppcd/chaincfg?status.png)]
(http://godoc.org/github.com/ppcsuite/ppcd/chaincfg)

Full `go doc` style documentation for the project can be viewed online without
installing this package by using the GoDoc site
[here](http://godoc.org/github.com/ppcsuite/ppcd/chaincfg).

You can also view the documentation locally once the package is installed with
the `godoc` tool by running `godoc -http=":6060"` and pointing your browser to
http://localhost:6060/pkg/github.com/ppcsuite/ppcd/chaincfg

## Installation

```bash
$ go get github.com/ppcsuite/ppcd/chaincfg
```

## GPG Verification Key

All official release tags are signed by Conformal so users can ensure the code
has not been tampered with and is coming from the btcsuite developers.  To
verify the signature perform the following:

- Download the public key from the Conformal website at
  https://opensource.conformal.com/GIT-GPG-KEY-conformal.txt

- Import the public key into your GPG keyring:
  ```bash
  gpg --import GIT-GPG-KEY-conformal.txt
  ```

- Verify the release tag with the following command where `TAG_NAME` is a
  placeholder for the specific tag:
  ```bash
  git tag -v TAG_NAME
  ```

## License

Package chaincfg is licensed under the [copyfree](http://copyfree.org) ISC
License.
