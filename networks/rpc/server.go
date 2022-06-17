// Modifications Copyright 2018 The klaytn Authors
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
// This file is derived from rpc/server.go (2018/06/04).
// Modified and improved for the klaytn development.

package rpc

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"gopkg.in/fatih/set.v0"
)

const MetadataApi = "rpc"

// CodecOption specifies which type of messages this codec supports
type CodecOption int

const (
	// OptionMethodInvocation is an indication that the codec supports RPC method calls
	OptionMethodInvocation CodecOption = 1 << iota

	// OptionSubscriptions is an indication that the codec supports RPC notifications
	OptionSubscriptions = 1 << iota // support pub sub

	// pendingRequestLimit is a limit for concurrent RPC method calls
	pendingRequestLimit = 200000
)

var (
	// ConcurrencyLimit is a limit for the number of concurrency connection for RPC servers.
	// It can be overwritten by rpc.concurrencylimit flag
	ConcurrencyLimit = 3000

	// pendingRequestCount is a total number of concurrent RPC method calls
	pendingRequestCount int64 = 0

	// TODO-Klaytn: move websocket configurations to Config struct in /network/rpc/server.go
	// MaxSubscriptionPerWSConn is a maximum number of subscription for a websocket connection
	MaxSubscriptionPerWSConn int32 = 3000

	// WebsocketReadDeadline is the read deadline on the underlying network connection in seconds. 0 means read will not timeout
	WebsocketReadDeadline int64 = 0

	// WebsocketWriteDeadline is the write deadline on the underlying network connection in seconds. 0 means write will not timeout
	WebsocketWriteDeadline int64 = 0

	// MaxWebsocketConnections is a maximum number of websocket connections
	MaxWebsocketConnections int32 = 3000

	// NonEthCompatible is a bool value that determines whether to use return formatting of the eth namespace API  provided for compatibility.
	// It can be overwritten by rpc.eth.noncompatible flag
	NonEthCompatible = false
)

// NewServer will create a new server instance with no registered handlers.
func NewServer() *Server {
	server := &Server{
		services:    make(serviceRegistry),
		codecs:      set.New(),
		run:         1,
		wsConnCount: 0,
	}

	// register a default service which will provide meta information about the RPC service such as the services and
	// methods it offers.
	rpcService := &RPCService{server}
	server.RegisterName(MetadataApi, rpcService)

	return server
}

// RPCService gives meta information about the server.
// e.g. gives information about the loaded modules.
type RPCService struct {
	server *Server
}

// Modules returns the list of RPC services with their version number
func (s *RPCService) Modules() map[string]string {
	modules := make(map[string]string)
	for name := range s.server.services {
		modules[name] = "1.0"
	}
	return modules
}

func (s *Server) GetServices() serviceRegistry {
	return s.services
}

// RegisterName will create a service for the given rcvr type under the given name. When no methods on the given rcvr
// match the criteria to be either a RPC method or a subscription an error is returned. Otherwise a new service is
// created and added to the service collection this server instance serves.
func (s *Server) RegisterName(name string, rcvr interface{}) error {
	if s.services == nil {
		s.services = make(serviceRegistry)
	}

	svc := new(service)
	svc.typ = reflect.TypeOf(rcvr)
	rcvrVal := reflect.ValueOf(rcvr)

	if name == "" {
		return fmt.Errorf("no service name for type %s", svc.typ.String())
	}
	if !isExported(reflect.Indirect(rcvrVal).Type().Name()) {
		return fmt.Errorf("%s is not exported", reflect.Indirect(rcvrVal).Type().Name())
	}

	methods, subscriptions := suitableCallbacks(rcvrVal, svc.typ)

	// already a previous service register under given sname, merge methods/subscriptions
	if regsvc, present := s.services[name]; present {
		if len(methods) == 0 && len(subscriptions) == 0 {
			return fmt.Errorf("Service %T doesn't have any suitable methods/subscriptions to expose", rcvr)
		}
		for _, m := range methods {
			regsvc.callbacks[formatName(m.method.Name)] = m
		}
		for _, s := range subscriptions {
			regsvc.subscriptions[formatName(s.method.Name)] = s
		}
		return nil
	}

	svc.name = name
	svc.callbacks, svc.subscriptions = methods, subscriptions

	if len(svc.callbacks) == 0 && len(svc.subscriptions) == 0 {
		return fmt.Errorf("Service %T doesn't have any suitable methods/subscriptions to expose", rcvr)
	}

	s.services[svc.name] = svc
	return nil
}

var exeCount int64

// serveRequest will reads requests from the codec, calls the RPC callback and
// writes the response to the given codec.
//
// If singleShot is true it will process a single request, otherwise it will handle
// requests until the codec returns an error when reading a request (in most cases
// an EOF). It executes requests in parallel when singleShot is false.
func (s *Server) serveRequest(ctx context.Context, codec ServerCodec, singleShot bool, options CodecOption) error {
	var pend sync.WaitGroup

	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			logger.Error(string(buf))
		}
		s.codecsMu.Lock()
		s.codecs.Remove(codec)
		s.codecsMu.Unlock()
	}()

	//	ctx, cancel := context.WithCancel(context.Background())
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	defer func() {
		notifier, supported := NotifierFromContext(ctx)
		if supported { // interface doesn't support subscriptions (e.g. http)
			notifier.unsubscribeAll()
		}
	}()
	// if the codec supports notification include a notifier that callbacks can use
	// to send notification to clients. It is thight to the codec/connection. If the
	// connection is closed the notifier will stop and cancels all active subscriptions.
	if options&OptionSubscriptions == OptionSubscriptions {
		ctx = context.WithValue(ctx, notifierKey{}, newNotifier(codec))
	}
	s.codecsMu.Lock()
	if atomic.LoadInt32(&s.run) != 1 { // server stopped
		s.codecsMu.Unlock()
		return &shutdownError{}
	}
	s.codecs.Add(codec)
	s.codecsMu.Unlock()

	// subscriptionCount counts and limits active subscriptions to avoid resource exhaustion
	subscriptionCount := int32(0)

	// test if the server is ordered to stop
	for atomic.LoadInt32(&s.run) == 1 {
		reqs, batch, err := s.readRequest(codec)
		rpcTotalRequestsCounter.Inc(int64(len(reqs)))
		if err != nil {
			// If a parsing error occurred, send an error
			if err.Error() != "EOF" {
				rpcErrorResponsesCounter.Inc(int64(len(reqs)))
				logger.Debug(fmt.Sprintf("read error %v\n", err))
				codec.Write(codec.CreateErrorResponse(nil, err))
			}
			// Error or end of stream, wait for requests and tear down
			pend.Wait()
			return nil
		}

		if atomic.LoadInt64(&pendingRequestCount) > pendingRequestLimit {
			rpcErrorResponsesCounter.Inc(int64(len(reqs)))
			err := &invalidRequestError{"server requests exceed the limit"}
			logger.Debug(fmt.Sprintf("request error %v\n", err))
			codec.Write(codec.CreateErrorResponse(nil, err))
			// Error or end of stream, wait for requests and tear down
			pend.Wait()
			return nil
		}

		// check if server is ordered to shutdown and return an error
		// telling the client that his request failed.
		if atomic.LoadInt32(&s.run) != 1 {
			rpcErrorResponsesCounter.Inc(int64(len(reqs)))
			err = &shutdownError{}
			if batch {
				resps := make([]interface{}, len(reqs))
				for i, r := range reqs {
					resps[i] = codec.CreateErrorResponse(&r.id, err)
				}
				codec.Write(resps)
			} else {
				codec.Write(codec.CreateErrorResponse(&reqs[0].id, err))
			}
			return nil
		}

		//if reqs[0].callb.method.Name == "SendRawTransaction" {
		//	atomic.AddInt64(&exeCount,1)
		//	logger.Error("## request", "method", reqs[0].callb.method.Name,"#call",atomic.LoadInt64(&exeCount))
		//}

		// If a single shot request is executing, run and return immediately
		if singleShot {
			if batch {
				s.execBatch(ctx, codec, reqs, &subscriptionCount)
			} else {
				s.exec(ctx, codec, reqs[0], &subscriptionCount)
			}
			return nil
		}
		// For multi-shot connections, start a goroutine to serve and loop back
		pend.Add(1)
		atomic.AddInt64(&pendingRequestCount, 1)
		rpcPendingRequestsCount.Inc(int64(len(reqs)))
		go func(reqs []*serverRequest, batch bool) {
			defer func() {
				atomic.AddInt64(&pendingRequestCount, -1)
				if err := recover(); err != nil {
					const size = 64 << 10
					buf := make([]byte, size)
					buf = buf[:runtime.Stack(buf, false)]
					logger.Error(string(buf))
				}
			}()

			defer pend.Done()
			if batch {
				s.execBatch(ctx, codec, reqs, &subscriptionCount)
			} else {
				s.exec(ctx, codec, reqs[0], &subscriptionCount)
			}
		}(reqs, batch)
	}
	return nil
}

// ServeCodec reads incoming requests from codec, calls the appropriate callback and writes the
// response back using the given codec. It will block until the codec is closed or the server is
// stopped. In either case the codec is closed.
func (s *Server) ServeCodec(codec ServerCodec, options CodecOption) {
	defer codec.Close()
	s.serveRequest(context.Background(), codec, false, options)
}

// ServeSingleRequest reads and processes a single RPC request from the given codec. It will not
// close the codec unless a non-recoverable error has occurred. Note, this method will return after
// a single request has been processed!
func (s *Server) ServeSingleRequest(ctx context.Context, codec ServerCodec, options CodecOption) {
	s.serveRequest(ctx, codec, true, options)
}

// Stop will stop reading new requests, wait for stopPendingRequestTimeout to allow pending requests to finish,
// close all codecs which will cancel pending requests/subscriptions.
func (s *Server) Stop() {
	if atomic.CompareAndSwapInt32(&s.run, 1, 0) {
		logger.Debug("RPC Server shutdown initiatied")
		s.codecsMu.Lock()
		defer s.codecsMu.Unlock()
		s.codecs.Each(func(c interface{}) bool {
			c.(ServerCodec).Close()
			return true
		})
	}
}

// createSubscription will call the subscription callback and returns the subscription id or error.
func (s *Server) createSubscription(ctx context.Context, c ServerCodec, req *serverRequest) (ID, error) {
	// subscription have as first argument the context following optional arguments
	args := []reflect.Value{req.callb.rcvr, reflect.ValueOf(ctx)}
	args = append(args, req.args...)
	reply := req.callb.method.Func.Call(args)

	if !reply[1].IsNil() { // subscription creation failed
		return "", reply[1].Interface().(error)
	}

	return reply[0].Interface().(*Subscription).ID, nil
}

var callCount = 0
var callSendTx = 0

// handle executes a request and returns the response from the callback.
func (s *Server) handle(ctx context.Context, codec ServerCodec, req *serverRequest, subCnt *int32) (interface{}, func()) {
	method := ""
	if req.callb != nil {
		method = fmt.Sprintf("%s%s%s", req.svcname, serviceMethodSeparator, req.callb.method.Name)
	}
	logger.Trace("Request info", "reqId", fmt.Sprintf("%s", req.id), "reqErr", req.err, "isUnsubscribe", req.isUnsubscribe, "reqMethod", method)

	if req.err != nil {
		rpcErrorResponsesCounter.Inc(1)
		return codec.CreateErrorResponse(&req.id, req.err), nil
	}

	if req.isUnsubscribe { // cancel subscription, first param must be the subscription id
		if len(req.args) >= 1 && req.args[0].Kind() == reflect.String {
			notifier, supported := NotifierFromContext(ctx)
			if !supported { // interface doesn't support subscriptions (e.g. http)
				rpcErrorResponsesCounter.Inc(1)
				return codec.CreateErrorResponse(&req.id, &callbackError{ErrNotificationsUnsupported.Error()}), nil
			}

			subid := ID(req.args[0].String())
			if err := notifier.unsubscribe(subid); err != nil {
				rpcErrorResponsesCounter.Inc(1)
				return codec.CreateErrorResponse(&req.id, &callbackError{err.Error()}), nil
			}

			atomic.AddInt32(subCnt, -1)
			rpcSuccessResponsesCounter.Inc(1)
			return codec.CreateResponse(req.id, true), nil
		}
		rpcErrorResponsesCounter.Inc(1)
		wsUnsubscriptionReqCounter.Inc(1)
		return codec.CreateErrorResponse(&req.id, &invalidParamsError{"Expected subscription id as first argument"}), nil
	}

	if req.callb.isSubscribe {
		if atomic.LoadInt32(subCnt) >= MaxSubscriptionPerWSConn {
			return codec.CreateErrorResponse(&req.id, &callbackError{
				fmt.Sprintf("Maximum %d subscriptions are allowed for a websocket connection. "+
					"The limit can be updated with 'admin_setMaxSubscriptionPerWSConn' API", MaxSubscriptionPerWSConn),
			}), nil
		}

		subid, err := s.createSubscription(ctx, codec, req)
		if err != nil {
			rpcErrorResponsesCounter.Inc(1)
			return codec.CreateErrorResponse(&req.id, &callbackError{err.Error()}), nil
		}

		// active the subscription after the sub id was successfully sent to the client
		activateSub := func() {
			notifier, _ := NotifierFromContext(ctx)
			notifier.activate(subid, req.svcname)
		}
		atomic.AddInt32(subCnt, 1)
		rpcSuccessResponsesCounter.Inc(1)
		wsSubscriptionReqCounter.Inc(1)
		return codec.CreateResponse(req.id, subid), activateSub
	}

	// regular RPC call, prepare arguments
	if len(req.args) != len(req.callb.argTypes) {
		rpcErr := &invalidParamsError{fmt.Sprintf("%s%s%s expects %d parameters, got %d",
			req.svcname, serviceMethodSeparator, req.callb.method.Name,
			len(req.callb.argTypes), len(req.args))}
		rpcErrorResponsesCounter.Inc(1)
		return codec.CreateErrorResponse(&req.id, rpcErr), nil
	}

	arguments := []reflect.Value{req.callb.rcvr}
	if req.callb.hasCtx {
		arguments = append(arguments, reflect.ValueOf(ctx))
	}
	if len(req.args) > 0 {
		arguments = append(arguments, req.args...)
	}

	//if req.callb.method.Name == "SendRawTransaction" {
	//	callSendTx++
	//}
	//if req.callb.method.Name == "GetTransactionReceipt" {
	//	callCount++
	//}
	//logger.Error("### rpc.server", "#tx", callSendTx, "#receipt", callCount)

	// execute RPC method and return result
	reply := req.callb.method.Func.Call(arguments)
	if len(reply) == 0 {
		rpcSuccessResponsesCounter.Inc(1)
		return codec.CreateResponse(req.id, nil), nil
	}
	if req.callb.errPos >= 0 { // test if method returned an error
		if !reply[req.callb.errPos].IsNil() {
			e := reply[req.callb.errPos].Interface().(error)
			rpcErrorResponsesCounter.Inc(1)
			res := codec.CreateErrorResponse(&req.id, e)
			logger.Trace("RPCError", "reqId", fmt.Sprintf("%s", req.id), "err", e, "method", fmt.Sprintf("%s%s%s", req.svcname, serviceMethodSeparator, req.callb.method.Name))
			return res, nil
		}
	}

	rpcSuccessResponsesCounter.Inc(1)
	return codec.CreateResponse(req.id, reply[0].Interface()), nil
}

// exec executes the given request and writes the result back using the codec.
func (s *Server) exec(ctx context.Context, codec ServerCodec, req *serverRequest, subCnt *int32) {
	var response interface{}
	var callback func()
	if req.err != nil {
		rpcErrorResponsesCounter.Inc(1)
		response = codec.CreateErrorResponse(&req.id, req.err)
	} else {
		response, callback = s.handle(ctx, codec, req, subCnt)
	}

	if err := codec.Write(response); err != nil {
		logger.Error(fmt.Sprintf("%v\n", err))
		codec.Close()
	}

	// when request was a subscribe request this allows these subscriptions to be actived
	if callback != nil {
		callback()
	}
}

// execBatch executes the given requests and writes the result back using the codec.
// It will only write the response back when the last request is processed.
func (s *Server) execBatch(ctx context.Context, codec ServerCodec, requests []*serverRequest, subCnt *int32) {
	responses := make([]interface{}, len(requests))
	var callbacks []func()
	for i, req := range requests {
		if req.err != nil {
			rpcErrorResponsesCounter.Inc(1)
			responses[i] = codec.CreateErrorResponse(&req.id, req.err)
		} else {
			var callback func()
			if responses[i], callback = s.handle(ctx, codec, req, subCnt); callback != nil {
				callbacks = append(callbacks, callback)
			}
		}
	}

	if err := codec.Write(responses); err != nil {
		logger.Error(fmt.Sprintf("%v\n", err))
		codec.Close()
	}

	// when request holds one of more subscribe requests this allows these subscriptions to be activated
	for _, c := range callbacks {
		c()
	}
}

// readRequest requests the next (batch) request from the codec. It will return the collection
// of requests, an indication if the request was a batch, the invalid request identifier and an
// error when the request could not be read/parsed.
func (s *Server) readRequest(codec ServerCodec) ([]*serverRequest, bool, Error) {
	reqs, batch, err := codec.ReadRequestHeaders()
	if err != nil {
		return nil, batch, err
	}

	requests := make([]*serverRequest, len(reqs))

	// verify requests
	for i, r := range reqs {
		var ok bool
		var svc *service

		if r.err != nil {
			requests[i] = &serverRequest{id: r.id, err: r.err}
			continue
		}

		if r.isPubSub && strings.HasSuffix(r.method, unsubscribeMethodSuffix) {
			requests[i] = &serverRequest{id: r.id, isUnsubscribe: true}
			argTypes := []reflect.Type{reflect.TypeOf("")} // expect subscription id as first arg
			if args, err := codec.ParseRequestArguments(argTypes, r.params); err == nil {
				requests[i].args = args
			} else {
				requests[i].err = &invalidParamsError{err.Error()}
			}
			continue
		}

		if NonEthCompatible && r.service == "eth" {
			// when NonEthCompatible is true, the return formatting for the eth namespace API provided for Ethereum compatibility is disabled.
			// convert ethereum namespace to klay namespace.
			r.service = "klay"
		}

		if svc, ok = s.services[r.service]; !ok { // rpc method isn't available
			logger.Trace("rpc: got request from unsupported API namespace", "requestedNamespaces", r.service)
			requests[i] = &serverRequest{id: r.id, err: &methodNotFoundError{r.service, r.method}}
			continue
		}

		if r.isPubSub { // eth_subscribe, r.method contains the subscription method name
			if callb, ok := svc.subscriptions[r.method]; ok {
				requests[i] = &serverRequest{id: r.id, svcname: svc.name, callb: callb}
				if r.params != nil && len(callb.argTypes) > 0 {
					argTypes := []reflect.Type{reflect.TypeOf("")}
					argTypes = append(argTypes, callb.argTypes...)
					if args, err := codec.ParseRequestArguments(argTypes, r.params); err == nil {
						requests[i].args = args[1:] // first one is service.method name which isn't an actual argument
					} else {
						requests[i].err = &invalidParamsError{err.Error()}
					}
				}
			} else {
				requests[i] = &serverRequest{id: r.id, err: &methodNotFoundError{r.service, r.method}}
			}
			continue
		}

		if callb, ok := svc.callbacks[r.method]; ok { // lookup RPC method
			requests[i] = &serverRequest{id: r.id, svcname: svc.name, callb: callb}
			if r.params != nil && len(callb.argTypes) > 0 {
				if args, err := codec.ParseRequestArguments(callb.argTypes, r.params); err == nil {
					requests[i].args = args
				} else {
					requests[i].err = &invalidParamsError{err.Error()}
				}
			}
			continue
		}

		requests[i] = &serverRequest{id: r.id, err: &methodNotFoundError{r.service, r.method}}
	}

	return requests, batch, nil
}
