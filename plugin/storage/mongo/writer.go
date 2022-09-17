package mongo

import (
	"context"
	"github.com/jaegertracing/jaeger/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

type SpanWriter struct {
	mongoClient      *mongo.Client
	collection       *mongo.Collection
	collectionParsed *mongo.Collection
}

func (receiver SpanWriter) WriteSpan(ctx context.Context, span *model.Span) error {
	mSpan := Span{
		TraceID:       span.TraceID.String(),
		SpanID:        span.SpanID.String(),
		OperationName: span.OperationName,
		StartTime:     span.StartTime,
		Duration:      span.Duration.Microseconds(),
		References:    convertReferences(span),
		ProcessID:     span.ProcessID,
		Process:       convertProcess(span.Process),
		Tags:          convertKeyValues(span.Tags),
		Warnings:      span.Warnings,
	}
	b, _ := bson.Marshal(mSpan)

	_, err := receiver.collection.InsertOne(ctx, b)
	go receiver.WriteCustomSpan(ctx, span)
	return err

}

// WriteCustomSpan Only for our own project,you could remove it.
func (receiver SpanWriter) WriteCustomSpan(ctx context.Context, span *model.Span) error {

	kind, _ := span.GetSpanKind()
	code, _ := span.GetSpanStatus()
	service, _ := span.GetSpanService()

	parsedSpan := SpanParsed{
		TraceID:       span.TraceID.String(),
		SpanID:        span.SpanID.String(),
		OperationName: span.OperationName,
		StartTime:     span.StartTime,
		Duration:      span.Duration.Microseconds(),
		Type:          kind,
		StatusCode:    code,
		ParentSpan:    span.ParentSpanID().String(),
		Service:       service,
	}
	parsedSpanBson, _ := bson.Marshal(parsedSpan)
	_, err := receiver.collectionParsed.InsertOne(ctx, parsedSpanBson)
	return err
}

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
