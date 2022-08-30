package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"net/http"
	"os"
	"strings"
)

type DatadogTracer struct {
	Tags           []string
	Service        string
	KlaytnResponse bool
}

func newDatadogTracer() *DatadogTracer {
	ddTraceEnabled := os.Getenv("DD_TRACE_ENABLED")
	if strings.ToLower(ddTraceEnabled) != "true" {
		return nil
	}

	tracer.Start()

	tags := strings.Split(os.Getenv("DD_TRACE_HEADER_TAGS"), ",")
	service := os.Getenv("DD_SERVICE")
	klaytnResponse := strings.ToLower(os.Getenv("DD_KLAYTN_RPC_RESPONSE")) == "true"

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

		for _, header := range ddTracer.Tags {
			header = strings.TrimSpace(header)
			header = strings.ToLower(header)
			key := fmt.Sprintf("http.%s", header)

			var tag tracer.StartSpanOption
			if header == "remote_addr" {
				tag = tracer.Tag(key, r.RemoteAddr)
			} else {
				tag = tracer.Tag(key, r.Header.Get(header))
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
				var span, _ = tracer.SpanFromContext(r.Context())
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
	var span, _ = tracer.StartSpanFromContext(r.Context(), "response.batch")
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
