package rpc

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/klaytn/klaytn/common"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// dupWriter writes data to the buffer as well as http response
type dupWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w dupWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// KASAttrs contains identifications for a KAS request
type KASAttrs struct {
	ChainID      string `json:"x-chain-id"`
	AccountID    string `json:"x-account-id"`
	RequestID    string `json:"x-request-id"`
	ParentSpanID string `json:"x-b3-parentspanid,omitempty"`
	SpanID       string `json:"x-b3-spanid,omitempty"`
	TraceID      string `json:"x-b3-traceid,omitempty"`
}

func parseKASHeader(r *http.Request) KASAttrs {
	return KASAttrs{
		ChainID:      r.Header.Get("x-chain-id"),
		AccountID:    r.Header.Get("x-account-id"),
		RequestID:    r.Header.Get("x-request-id"),
		ParentSpanID: r.Header.Get("x-b3-parentspanid"),
		SpanID:       r.Header.Get("x-b3-spanid"),
		TraceID:      r.Header.Get("x-b3-traceid"),
	}
}

func newNewRelicApp() *newrelic.Application {
	appName := os.Getenv("NEWRELIC_APP_NAME")
	license := os.Getenv("NEWRELIC_LICENSE")
	if appName == "" && license == "" {
		return nil
	}

	nrApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(appName),
		newrelic.ConfigLicense(license),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		logger.Crit("failed to create NewRelic application. If you want to register a NewRelic HTTP handler," +
			" specify NEWRELIC_APP_NAME and NEWRELIC_LICENSE os environment variables with valid values. " +
			"If you don't want to register the handler, specify them with an empty string.")
	}

	logger.Info("NewRelic APM is enabled", "appName", appName)
	return nrApp
}

// newNewRelicHTTPHandler enables NewRelic web transaction monitor.
// It also prints error logs when RPC returns contains error messages.
func newNewRelicHTTPHandler(nrApp *newrelic.Application, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.ErrorWithStack("NewRelic http handler panic", "err", err)
			}
		}()

		reqMethod := ""

		// parse RPC requests
		reqs, isBatch, err := getRPCRequests(r)
		if err != nil || len(reqs) < 1 {
			// The error will be handled in `handler.ServeHTTP()` and printed with `printRPCErrorLog()`
			logger.Debug("failed to parse RPC request", "err", err, "len(reqs)", len(reqs))
		} else {
			reqMethod = reqs[0].Method
			if isBatch {
				reqMethod += "_batch"
			}
		}

		// new relic transaction name contains the first API method of the request
		txn := nrApp.StartTransaction(r.Method + " " + r.URL.String() + " " + reqMethod)
		defer txn.End()

		w = txn.SetWebResponse(w)
		txn.SetWebRequestHTTP(r)
		r = newrelic.RequestWithTransactionContext(r, txn)

		// duplicate writer
		dupW := &dupWriter{
			ResponseWriter: w,
			body:           bytes.NewBufferString(""),
		}

		// serve HTTP
		handler.ServeHTTP(dupW, r)

		// print RPC error logs if errors exist
		if isBatch {
			var rpcReturns []interface{}
			if err := json.Unmarshal(dupW.body.Bytes(), &rpcReturns); err == nil {
				for i, rpcReturn := range rpcReturns {
					if data, err := json.Marshal(rpcReturn); err == nil {
						// TODO-Klaytn: make the log level configurable or separate module name of the logger
						printRPCErrorLog(data, reqs[i].Method, r)
					}
				}
			}
		} else {
			// TODO-Klaytn: make the log level configurable or separate module name of the logger
			printRPCErrorLog(dupW.body.Bytes(), reqMethod, r)
		}
	})
}

// getRPCRequests copies a http request body data and parses RPC requests from the data.
// It returns a slice of RPC request, an indication if these requests are in batch, and an error.
// Ethereum returns []*jsonrpcMessage, which replaces []rpcRequest
func getRPCRequests(r *http.Request) ([]*jsonrpcMessage, bool, error) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error("cannot read a request body", "err", err)
		return nil, false, err
	}

	r.Body = ioutil.NopCloser(bytes.NewReader(reqBody))
	body := io.LimitReader(r.Body, int64(common.MaxRequestContentLength))
	conn := &httpServerConn{Reader: body, Writer: bytes.NewBufferString(""), r: r}

	codec := NewCodec(conn)

	defer codec.close()

	return codec.readBatch()
}

// printRPCErrorLog prints an error log if responseBody contains RPC error message.
// It does nothing if responseBody doesn't contain RPC error message.
func printRPCErrorLog(responseBody []byte, method string, r *http.Request) {
	// check whether the responseBody contains json error
	var rpcError jsonErrResponse
	if err := json.Unmarshal(responseBody, &rpcError); err != nil || rpcError.Error.Code == 0 {
		// do nothing if the responseBody didn't contain json error data
		return
	}

	// parse KAS HTTP header
	kasHeader := parseKASHeader(r)
	kasHeaderJson, err := json.Marshal(kasHeader)
	if err != nil {
		logger.Error("failed to marshal a KAS HTTP header", "err", err, "kasHeader", kasHeader)
	}

	// print RPC error log
	logger.Error("RPC error response", "rpcErr", rpcError.Error.Message, "kasHeader", string(kasHeaderJson),
		"method", method)
}
