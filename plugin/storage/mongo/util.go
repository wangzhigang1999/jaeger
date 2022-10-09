package mongo

import (
	"encoding/json"
	"fmt"
	"github.com/jaegertracing/jaeger/model"
	"strings"
)

var unNecessaryUserAgent = []string{"Prometheus"}
var unNecessaryURL = []string{"metrics"}

func convertProcess(process *model.Process) Process {
	return Process{
		ServiceName: process.ServiceName,
		Tags:        convertKeyValues(process.Tags),
	}
}

func convertReferences(span *model.Span) []Reference {
	out := make([]Reference, 0, len(span.References))
	for _, ref := range span.References {
		out = append(out, Reference{
			RefType: convertRefType(ref.RefType),
			TraceID: TraceID(ref.TraceID.String()),
			SpanID:  SpanID(ref.SpanID.String()),
		})
	}
	return out
}

func convertRefType(refType model.SpanRefType) ReferenceType {
	if refType == model.FollowsFrom {
		return FollowsFrom
	}
	return ChildOf
}

func convertKeyValues(keyValues model.KeyValues) []KeyValue {
	kvs := make([]KeyValue, 0)
	for _, kv := range keyValues {
		if kv.GetVType() != model.BinaryType {
			kvs = append(kvs, convertKeyValue(kv))
		}
	}
	return kvs
}

func convertKeyValue(kv model.KeyValue) KeyValue {
	return KeyValue{
		Key:   kv.Key,
		Type:  ValueType(strings.ToLower(kv.VType.String())),
		Value: kv.AsString(),
	}
}

func toJson(v interface{}) string {
	marshal, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(marshal)
}

func Extractor(span *model.Span) SpanParsed {
	kind, _ := span.GetSpanKind()
	code, _ := span.GetSpanStatus()
	service, _ := span.GetSpanService()
	url, _ := span.GetSpanHttpUrl()
	respSize, _ := span.GetSpanRespSize()
	reqSize, _ := span.GetSpanReqSize()
	nodeId, _ := span.GetSpanNodeId()
	method, _ := span.GetSpanMethod()

	realOperationName := fmt.Sprintf("%s~%s", method, parseUrl(url))

	parsedSpan := SpanParsed{
		TraceID:           span.TraceID.String(),
		SpanID:            span.SpanID.String(),
		OperationName:     span.OperationName,
		StartTime:         span.StartTime,
		Duration:          span.Duration,
		Service:           service,
		Type:              kind,
		StatusCode:        code,
		ParentSpan:        span.ParentSpanID().String(),
		URL:               url,
		ResponseSize:      respSize,
		RequestSize:       reqSize,
		CMDB:              nodeId,
		Method:            method,
		RealOperationName: realOperationName,
	}
	return parsedSpan
}

func parseUrl(url string) string {
	if strings.Contains(url, "http:/") {
		var tmp []string
		split := strings.Split(url, "/")
		for _, str := range split {
			if len(str) < 20 {
				if !strings.Contains(str, "30001") {
					tmp = append(tmp, str)
				} else {
					tmp = append(tmp, "front")
				}
			}
		}
		return strings.Join(tmp, "/")
	}
	return url
}

func SpanNeedFilter(span *model.Span) bool {
	// 如果没有URL，可以认为它不是来自Istio
	url, fromIstio := span.GetSpanHttpUrl()
	if !fromIstio || StringContainsAny(url, unNecessaryURL) {
		return true
	}
	component, b := span.GetSpanComponent()
	if !b || component != "proxy" {
		return true
	}
	agent, b2 := span.GetSpanUserAgent()
	if !b2 || InSlice(agent, unNecessaryUserAgent) {
		return true
	}

	name := span.OperationName
	if strings.Contains(name, "zipkin-mongo") {
		return true
	}

	return false
}

func InSlice(val interface{}, slice []string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func StringContainsAny(str string, slice []string) bool {
	for index := range slice {
		if strings.Contains(str, slice[index]) {
			return true
		}
	}
	return false
}
