// Copyright (c) 2014-2015 PPCD developers.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"io"
)

// MsgCheckPoint implements the Message interface and represents a bitcoin reject
// message.
//
// This message was not added until protocol version RejectVersion.
type MsgCheckPoint struct {
	// Cmd is the command for the message which was rejected such as
	// as CmdBlock or CmdTx.  This can be obtained from the Command function
	// of a Message.
	Cmd string

	// Hash identifies a specific block or transaction that was rejected
	// and therefore only applies the MsgBlock and MsgTx messages.
	Hash ShaHash
}

// BtcDecode decodes r using the bitcoin protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgCheckPoint) BtcDecode(r io.Reader, pver uint32) error {

	// Command that was rejected.
	cmd, err := readVarString(r, pver)
	if err != nil {
		return err
	}
	msg.Cmd = cmd

	// CmdBlock and CmdTx messages have an additional hash field that
	// identifies the specific block or transaction.
	err = readElement(r, &msg.Hash)
	if err != nil {
		return err
	}

	return nil
}

// BtcEncode encodes the receiver to w using the bitcoin protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgCheckPoint) BtcEncode(w io.Writer, pver uint32) error {

	// Command that was rejected.
	err := writeVarString(w, pver, msg.Cmd)
	if err != nil {
		return err
	}

	// CmdBlock and CmdTx messages have an additional hash field that
	// identifies the specific block or transaction.
	err = writeElement(w, &msg.Hash)
	if err != nil {
		return err
	}

	return nil
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgCheckPoint) Command() string {
	return CmdCheckPoint
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgCheckPoint) MaxPayloadLength(pver uint32) uint32 {
	plen := MaxMessagePayload
	return uint32(plen)
}

// NewMsgCheckPoint returns a new bitcoin reject message that conforms to the
// Message interface.  See MsgCheckPoint for details.
func NewMsgCheckPoint(command string) *MsgCheckPoint {
	return &MsgCheckPoint{
		Cmd: command,
	}
}
