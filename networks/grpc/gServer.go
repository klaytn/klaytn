// Modifications Copyright 2019 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from rpc/http.go (2018/06/04).
// Modified and improved for the klaytn development.

package grpc

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var logger = log.NewModuleLogger(log.NetworksGRPC)

type Listener struct {
	Addr       string
	handler    *rpc.Server
	grpcServer *grpc.Server
}

// grpcReadWriteNopCloser wraps an io.Reader and io.Writer with a NOP Close method.
type grpcReadWriteNopCloser struct {
	io.Reader
	io.Writer
}

// Close does nothing and returns always nil.
func (t *grpcReadWriteNopCloser) Close() error {
	return nil
}

// klaytnServer is an implementation of KlaytnNodeServer.
type klaytnServer struct {
	handler *rpc.Server
}

type grpcWriter struct {
	writer   KlaytnNode_SubscribeServer
	writeErr chan error
}

func (gw *grpcWriter) Write(p []byte) (n int, err error) {
	resp := &RPCResponse{Payload: p}
	if err := gw.writer.Send(resp); err != nil {
		if gw.writeErr != nil {
			gw.writeErr <- err
		}
		return 0, err
	}
	return len(p), nil
}

type bufWriter struct {
	writer   io.Writer
	writeErr chan error
	writeOk  chan []byte
}

func (gw *bufWriter) Write(p []byte) (n int, err error) {
	if _, err := gw.writer.Write(p); err != nil {
		gw.writeErr <- err
		return 0, err
	}
	gw.writeOk <- p
	return len(p), nil
}

// BiCall handles bidirectional communication between client and server.
func (kns *klaytnServer) BiCall(stream KlaytnNode_BiCallServer) error {
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			logger.Error("error reading the request", "err", err)
			return err
		}

		preader := bytes.NewReader(request.Params)

		// Create a custom encode/decode pair to enforce payload size and number encoding
		encoder := func(v interface{}) error {
			msg, err := json.Marshal(v)
			if err != nil {
				return err
			}
			resp := &RPCResponse{Payload: msg}
			err = stream.Send(resp)
			if err != nil {
				return err
			}
			return err
		}
		decoder := func(v interface{}) error {
			dec := json.NewDecoder(preader)
			dec.UseNumber()
			return dec.Decode(v)
		}

		ctx := context.Background()

		reader := bufio.NewReaderSize(preader, common.MaxRequestContentLength)
		kns.handler.ServeSingleRequest(ctx, rpc.NewCodec(&grpcReadWriteNopCloser{reader, &grpcWriter{stream, nil}}, encoder, decoder), rpc.OptionMethodInvocation|rpc.OptionSubscriptions)
	}
}

// only server can send message to client repeatedly
func (kns *klaytnServer) Subscribe(request *RPCRequest, stream KlaytnNode_SubscribeServer) error {
	var (
		writeErr = make(chan error, 1)
		readErr  = make(chan error, 1)
	)

	preader := bytes.NewReader(request.Params)

	// Create a custom encode/decode pair to enforce payload size and number encoding
	encoder := func(v interface{}) error {
		msg, err := json.Marshal(v)
		if err != nil {
			return err
		}
		resp := &RPCResponse{Payload: msg}
		err = stream.Send(resp)
		if err != nil {
			writeErr <- err
			return err
		}
		return err
	}
	decoder := func(v interface{}) error {
		dec := json.NewDecoder(preader)
		dec.UseNumber()
		err := dec.Decode(v)
		if err != nil {
			readErr <- err
		}
		return err
	}

	ctx := context.Background()

	reader := bufio.NewReaderSize(preader, common.MaxRequestContentLength)
	kns.handler.ServeSingleRequest(ctx, rpc.NewCodec(&grpcReadWriteNopCloser{reader, &grpcWriter{stream, writeErr}}, encoder, decoder), rpc.OptionMethodInvocation|rpc.OptionSubscriptions)

	var err error
loop:
	for {
		select {
		case err = <-writeErr:
			logger.Warn("fail to write", "err", err)
			break loop
		case err = <-readErr:
			logger.Warn("fail to read", "err", err)
			break loop
		}
	}

	return nil
}

// general RPC call, such as one-to-one communication
func (kns *klaytnServer) Call(ctx context.Context, request *RPCRequest) (*RPCResponse, error) {
	var (
		err      error
		writeErr = make(chan error, 1)
		readErr  = make(chan error, 1)
		writeOk  = make(chan []byte, 1)
	)

	preader := bytes.NewReader(request.Params)

	var res bytes.Buffer
	writer := &bufWriter{&res, writeErr, writeOk}

	// Create a custom encode/decode pair to enforce payload size and number encoding
	encoder := func(v interface{}) error {
		msg, err := json.Marshal(v)
		if err != nil {
			return err
		}
		_, err = writer.Write(msg)
		if err != nil {
			writeErr <- err
			return err
		}
		return err
	}
	decoder := func(v interface{}) error {
		dec := json.NewDecoder(preader)
		dec.UseNumber()
		err := dec.Decode(v)
		if err != nil {
			readErr <- err
		}
		return err
	}

	reader := bufio.NewReaderSize(preader, common.MaxRequestContentLength)
	kns.handler.ServeSingleRequest(ctx, rpc.NewCodec(&grpcReadWriteNopCloser{reader, writer}, encoder, decoder), rpc.OptionMethodInvocation)

loop:
	for {
		select {
		case _ = <-writeOk:
			break loop
		case err = <-writeErr:
			logger.Error("fail to write", "err", err)
			break loop
		case err = <-readErr:
			logger.Error("fail to read", "err", err)
			break loop
		}
	}

	if err == nil {
		return &RPCResponse{Payload: res.Bytes()}, nil
	} else {
		return &RPCResponse{Payload: []byte("")}, err
	}
}

// SetRPCServer sets the RPC server.
func (gs *Listener) SetRPCServer(handler *rpc.Server) {
	gs.handler = handler
}

func (gs *Listener) Start() {
	lis, err := net.Listen("tcp", gs.Addr)
	if err != nil {
		// TODO-Klaytn-gRPC Need to handle err
		logger.Error("failed to listen", "err", err)
	}
	gs.grpcServer = grpc.NewServer()

	RegisterKlaytnNodeServer(gs.grpcServer, &klaytnServer{handler: gs.handler})

	// Register reflection service on gRPC server.
	reflection.Register(gs.grpcServer)
	if err := gs.grpcServer.Serve(lis); err != nil {
		// TODO-Klaytn-gRPC Need to handle err
		logger.Error("failed to serve", "err", err)
	}
}

func (gs *Listener) Stop() {
	if gs.grpcServer != nil {
		gs.grpcServer.Stop()
	}
}
