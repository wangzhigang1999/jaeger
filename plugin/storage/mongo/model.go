package mongo

import (
	"time"
)

// ReferenceType is the reference type of one span to another
type ReferenceType string

// TraceID is the shared trace ID of all spans in the trace.
type TraceID string

// SpanID is the id of a span
type SpanID string

// ValueType is the type of value stored in KeyValue struct.
type ValueType string

const (
	// ChildOf means a span is the child of another span
	ChildOf ReferenceType = "CHILD_OF"
	// FollowsFrom means a span follows from another span
	FollowsFrom ReferenceType = "FOLLOWS_FROM"

	// StringType indicates a string value stored in KeyValue
	StringType ValueType = "string"
	// BoolType indicates a Boolean value stored in KeyValue
	BoolType ValueType = "bool"
	// Int64Type indicates a 64bit signed integer value stored in KeyValue
	Int64Type ValueType = "int64"
	// Float64Type indicates a 64bit float value stored in KeyValue
	Float64Type ValueType = "float64"
)

// MongoSpan is MongoDB representation of the domain span.
type MongoSpan struct {
	TraceID       string      `bson:"traceID" json:"traceID,omitempty"`
	SpanID        string      `bson:"spanID" json:"spanID,omitempty"`
	OperationName string      `bson:"operationName" json:"operationName,omitempty"`
	StartTime     time.Time   `bson:"startTime" json:"startTime"`         // microseconds since Unix epoch
	Duration      int64       `bson:"duration" json:"duration,omitempty"` // microseconds
	References    []Reference `bson:"references" json:"references,omitempty"`
	ProcessID     string      `bson:"processID" json:"processID,omitempty"`
	Process       Process     `bson:"process,omitempty" json:"process"`
	Tags          []KeyValue  `bson:"tags" json:"tags,omitempty"`
	Logs          []Log       `bson:"logs" json:"logs,omitempty"`
	Warnings      []string    `bson:"warnings" json:"warnings,omitempty"`
}

// SpanParsed only for our project,you may remove it if not needed.
type SpanParsed struct {
	TraceID           string        `json:"traceID,omitempty" bson:"traceID"`
	SpanID            string        `json:"spanID,omitempty" bson:"spanID"`
	OperationName     string        `json:"operationName,omitempty" bson:"operationName"`
	RealOperationName string        `json:"realOperationName,omitempty" bson:"realOperationName"`
	StartTime         time.Time     `json:"startTime" bson:"startTime"`
	Duration          time.Duration `json:"duration,omitempty" bson:"duration"`
	Service           string        `json:"service,omitempty" bson:"service"`
	CMDB              string        `json:"CMDB,omitempty" bson:"CMDB"`
	Type              string        `json:"type,omitempty" bson:"type"`
	StatusCode        string        `json:"statusCode,omitempty" bson:"statusCode"`
	ParentSpan        string        `json:"parentSpan,omitempty" bson:"parentSpan"`
	URL               string        `json:"url,omitempty" bson:"url"`
	RequestSize       string        `json:"requestSize,omitempty" bson:"requestSize"`
	ResponseSize      string        `json:"responseSize,omitempty" bson:"responseSize"`
	Method            string        `json:"method,omitempty" bson:"method"`
}

// Reference is a reference from one span to another
type Reference struct {
	RefType ReferenceType `bson:"refType" json:"refType,omitempty"`
	TraceID TraceID       `bson:"traceID" json:"traceID,omitempty"`
	SpanID  SpanID        `bson:"spanID" json:"spanID,omitempty"`
}

// Process is the process emitting a set of spans
type Process struct {
	ServiceName string     `bson:"serviceName" json:"serviceName,omitempty"`
	Tags        []KeyValue `bson:"tags" json:"tags,omitempty"`
}

// Log is a log emitted in a span
type Log struct {
	Timestamp uint64     `bson:"timestamp" json:"timestamp,omitempty"`
	Fields    []KeyValue `bson:"fields" json:"fields,omitempty"`
}

// KeyValue is a key-value pair with typed value.
type KeyValue struct {
	Key   string      `bson:"key" json:"key,omitempty"`
	Type  ValueType   `bson:"type,omitempty" json:"type,omitempty"`
	Value interface{} `bson:"value" json:"value,omitempty"`
}
