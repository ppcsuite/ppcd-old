wire
====

[![Build Status](https://travis-ci.org/ppcsuite/ppcd/wire.png?branch=master)]
(https://travis-ci.org/ppcsuite/ppcd/wire) [![Coverage Status]
(https://coveralls.io/repos/ppcsuite/ppcd/wire/badge.png?branch=master)]
(https://coveralls.io/r/ppcsuite/ppcd/wire?branch=master)
[![tip for next commit](http://peer4commit.com/projects/130.svg)](http://peer4commit.com/projects/130)

Package wire implements the bitcoin wire protocol.  A comprehensive suite of
tests with 100% test coverage is provided to ensure proper functionality.
Package wire is licensed under the liberal ISC license.

There is an associated blog post about the release of this package
[here](https://blog.conformal.com/wire-the-bitcoin-wire-protocol-package-from-btcd/).

This package is one of the core packages from btcd, an alternative full-node
implementation of bitcoin which is under active development by Conformal.
Although it was primarily written for btcd, this package has intentionally been
designed so it can be used as a standalone package for any projects needing to
interface with bitcoin peers at the wire protocol level.

## Documentation

[![GoDoc](https://godoc.org/github.com/ppcsuite/ppcd/wire?status.png)]
(http://godoc.org/github.com/ppcsuite/ppcd/wire)

Full `go doc` style documentation for the project can be viewed online without
installing this package by using the GoDoc site here:
http://godoc.org/github.com/ppcsuite/ppcd/wire

You can also view the documentation locally once the package is installed with
the `godoc` tool by running `godoc -http=":6060"` and pointing your browser to
http://localhost:6060/pkg/github.com/ppcsuite/ppcd/wire

## Installation

```bash
$ go get github.com/ppcsuite/ppcd/wire
```

## Bitcoin Message Overview

The bitcoin protocol consists of exchanging messages between peers. Each message
is preceded by a header which identifies information about it such as which
bitcoin network it is a part of, its type, how big it is, and a checksum to
verify validity. All encoding and decoding of message headers is handled by this
package.

To accomplish this, there is a generic interface for bitcoin messages named
`Message` which allows messages of any type to be read, written, or passed
around through channels, functions, etc. In addition, concrete implementations
of most of the currently supported bitcoin messages are provided. For these
supported messages, all of the details of marshalling and unmarshalling to and
from the wire using bitcoin encoding are handled so the caller doesn't have to
concern themselves with the specifics.

## Reading Messages Example

In order to unmarshal bitcoin messages from the wire, use the `ReadMessage`
function. It accepts any `io.Reader`, but typically this will be a `net.Conn`
to a remote node running a bitcoin peer.  Example syntax is:

```Go
	// Use the most recent protocol version supported by the package and the
	// main bitcoin network.
	pver := wire.ProtocolVersion
	chaincfg := wire.MainNet

	// Reads and validates the next bitcoin message from conn using the
	// protocol version pver and the bitcoin network chaincfg.  The returns
	// are a wire.Message, a []byte which contains the unmarshalled
	// raw payload, and a possible error.
	msg, rawPayload, err := wire.ReadMessage(conn, pver, chaincfg)
	if err != nil {
		// Log and handle the error
	}
```

See the package documentation for details on determining the message type.

## Writing Messages Example

In order to marshal bitcoin messages to the wire, use the `WriteMessage`
function. It accepts any `io.Writer`, but typically this will be a `net.Conn`
to a remote node running a bitcoin peer. Example syntax to request addresses
from a remote peer is:

```Go
	// Use the most recent protocol version supported by the package and the
	// main bitcoin network.
	pver := wire.ProtocolVersion
	chaincfg := wire.MainNet

	// Create a new getaddr bitcoin message.
	msg := wire.NewMsgGetAddr()

	// Writes a bitcoin message msg to conn using the protocol version
	// pver, and the bitcoin network chaincfg.  The return is a possible
	// error.
	err := wire.WriteMessage(conn, msg, pver, chaincfg)
	if err != nil {
		// Log and handle the error
	}
```

## GPG Verification Key

All official release tags are signed by Conformal so users can ensure the code
has not been tampered with and is coming from Conformal.  To verify the
signature perform the following:

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

Package wire is licensed under the [copyfree](http://copyfree.org) ISC
License.
