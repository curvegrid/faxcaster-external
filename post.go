package main

import (
	// copied from: https://github.com/farcasterxyz/hub-monorepo/tree/6ef0638492ea2aada241f5b0c55696737e152f37/packages/hub-web/examples/golang-submitmessage/protobufs
	"github.com/curvegrid/faxcaster/protobufs"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"google.golang.org/protobuf/proto"
)

type CastID struct {
	Fid  int    `json:"fid"`
	Hash string `json:"hash"`
}

type UntrustedData struct {
	Fid         int    `json:"fid"`
	URL         string `json:"url"`
	MessageHash string `json:"messageHash"`
	Timestamp   int64  `json:"timestamp"`
	Network     int    `json:"network"`
	ButtonIndex int    `json:"buttonIndex"`
	CastID      CastID `json:"castId"`
}

type TrustedData struct {
	MessageBytes string `json:"messageBytes"`
}

type DataRepresentation struct {
	UntrustedData UntrustedData `json:"untrustedData"`
	TrustedData   TrustedData   `json:"trustedData"`
}

func decodeMessageBytes(hexData string) (*protobufs.Message, error) {
	// change the hex string to bytes
	data, err := hexutil.Decode("0x" + hexData)
	if err != nil {
		return nil, err
	}

	// protobuf unmarshal the bytes
	var message protobufs.Message
	err = proto.Unmarshal([]byte(data), &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
