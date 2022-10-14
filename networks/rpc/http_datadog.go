package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type tagInfo struct {
	header string
	key    string
}

type DatadogTracer struct {
	Tags           []tagInfo
	Service        string
	KlaytnResponse bool
}

func newDatadogTracer() *DatadogTracer {
	if v := os.Getenv("DD_TRACE_ENABLED"); v != "" {
		ddTraceEnabled, err := strconv.ParseBool(v)
		if err != nil || ddTraceEnabled == false {
			return nil
		}
	}

	headers := strings.Split(os.Getenv("DD_TRACE_HEADER_TAGS"), ",")
	tags := make([]tagInfo, len(headers))
	for i, header := range headers {
		header = strings.TrimSpace(header)
		header = strings.ToLower(header)
		key := fmt.Sprintf("http.%s", header)

		tags[i].header = header
		tags[i].key = key
	}
	service := os.Getenv("DD_SERVICE")

	klaytnResponse := false
	if v := os.Getenv("DD_KLAYTN_RPC_RESPONSE"); v != "" {
		var err error
		klaytnResponse, err = strconv.ParseBool(v)
		if err != nil {
			return nil
		}
	}

	tracer.Start()

	return &DatadogTracer{tags, service, klaytnResponse}
}

func newDatadogHTTPHandler(ddTracer *DatadogTracer, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.ErrorWithStack("Datadog http handler panic", "err", err)
			}
		}()

		reqMethod := ""
		reqParam := ""

		// parse RPC requests
		reqs, isBatch, err := getRPCRequests(r)
		if err != nil || len(reqs) < 1 {
			// The error will be handled in `handler.ServeHTTP()` and printed with `printRPCErrorLog()`
			logger.Debug("failed to parse RPC request", "err", err, "len(reqs)", len(reqs))
		} else {
			reqMethod = reqs[0].method
			if isBatch {
				reqMethod += "_batch"
			}
			encoded, _ := json.Marshal(reqs[0].params)
			reqParam = string(encoded)
		}

		// datadog transaction name contains the first API method of the request
		resource := fmt.Sprintf("%s %s %s", r.Method, r.URL.String(), reqMethod)

		// duplicate writer
		dupW := &dupWriter{
			ResponseWriter: w,
			body:           bytes.NewBufferString(""),
		}

		spanOpts := []ddtrace.StartSpanOption{
			tracer.Tag("request.method", reqMethod),
			tracer.Tag("request.params", reqParam),
		}

		for _, ti := range ddTracer.Tags {
			var tag tracer.StartSpanOption
			if ti.header == "remote_addr" {
				tag = tracer.Tag(ti.key, r.RemoteAddr)
			} else {
				tag = tracer.Tag(ti.key, r.Header.Get(ti.header))
			}
			spanOpts = append(spanOpts, tag)
		}

		responseHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler.ServeHTTP(w, r)
			if isBatch {
				var rpcReturns []interface{}
				if err := json.Unmarshal(dupW.body.Bytes(), &rpcReturns); err == nil {
					for i, rpcReturn := range rpcReturns {
						ddTracer.traceBatchRpcResponse(r, rpcReturn, reqs[i], i)
					}
				}
			} else {
				span, _ := tracer.SpanFromContext(r.Context())
				ddTracer.traceRpcResponse(dupW.body.Bytes(), reqMethod, span)
			}
		})

		httptrace.TraceAndServe(responseHandler, dupW, r, &httptrace.ServeConfig{
			Service:     ddTracer.Service,
			Resource:    resource,
			QueryParams: true,
			SpanOpts:    spanOpts,
		})
	})
}

func (dt *DatadogTracer) traceBatchRpcResponse(r *http.Request, rpcReturn interface{}, req rpcRequest, offset int) {
	span, _ := tracer.StartSpanFromContext(r.Context(), "response.batch")
	defer span.Finish()
	span.SetTag("offset", offset)
	if data, err := json.Marshal(rpcReturn); err == nil {
		dt.traceRpcResponse(data, req.method, span)
	}
}

func (dt *DatadogTracer) traceRpcResponse(response []byte, method string, span tracer.Span) {
	var rpcError jsonErrResponse
	if err := json.Unmarshal(response, &rpcError); err == nil && rpcError.Error.Code != 0 {
		span.SetTag("response.code", rpcError.Error.Code)

		errJson, _ := json.Marshal(rpcError.Error)
		span.SetTag("response.error", string(errJson))
		span.SetTag("error", string(errJson))

		message := fmt.Sprintf("RPC error response %v", span)
		logger.Error(message, "rpcErr", rpcError.Error.Message, "method", method)
	} else {
		span.SetTag("response.code", 0)
	}

	if dt.KlaytnResponse {
		var rpcSuccess jsonSuccessResponse
		if err := json.Unmarshal(response, &rpcSuccess); err == nil && rpcSuccess.Result != nil {
			successJson, _ := json.Marshal(rpcSuccess.Result)
			span.SetTag("response.success", string(successJson))
		}
	}
}
