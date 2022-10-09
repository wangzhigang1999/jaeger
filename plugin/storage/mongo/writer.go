package mongo

import (
	"context"
	"fmt"
	"github.com/jaegertracing/jaeger/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SpanWriter struct {
	mongoClient      *mongo.Client
	collection       *mongo.Collection
	collectionParsed *mongo.Collection
	output           bool
}

func (receiver SpanWriter) WriteSpan(ctx context.Context, span *model.Span) error {
	spanNeedFilter := SpanNeedFilter(span)
	if spanNeedFilter {
		return nil
	}
	go func() {
		err := receiver.WriteDefault(ctx, span)
		if err != nil {

		}
	}()
	return receiver.WriteCustomSpan(ctx, span)
}

// WriteDefault This will write everything to mongo
func (receiver SpanWriter) WriteDefault(ctx context.Context, span *model.Span) error {
	mSpan := MongoSpan{
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
	if receiver.output {
		fmt.Println(toJson(mSpan))
	}
	return err
}

// WriteCustomSpan Only for our own project,you could remove it or change it.
func (receiver SpanWriter) WriteCustomSpan(ctx context.Context, span *model.Span) error {
	parsedSpan := Extractor(span)
	parsedSpanBson, _ := bson.Marshal(parsedSpan)
	_, err := receiver.collectionParsed.InsertOne(ctx, parsedSpanBson)
	return err
}
