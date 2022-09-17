<img align="centor" width="290" height="290" src="https://www.jaegertracing.io/img/jaeger-vector.svg">

# Jaeger

This is an unofficial fork of Jaeger.

In this fork, we implement a extremely simple mongodb-based storge backend.

## Usage

**Please read the official doc firstly**: https://github.com/jaegertracing/jaeger/blob/main/README.md

To set mongodb as storge backend, simply add these commands :

```bash
--span-storage.type="mongo"
--mongo_url="xxx" 			# default as localhost
--mogo_port=xxx 			# default as 27017
--mongo_database="xxx" 		# default as sock-shop
--mongo_colection="xxx" 	# default as span
--mongo_user="xxx" 			# default as root
--mongo_pass="xxx" 			# default as root
```

## Detail

>  The span will be saved in mongodb with the format below.

```go
// plugin/storage/mongo/model.go
type Span struct {
	TraceID       string      `bson:"traceID"`
	SpanID        string      `bson:"spanID"`
	OperationName string      `bson:"operationName"`
	StartTime     time.Time   `bson:"startTime"` // microseconds since Unix epoch
	Duration      int64       `bson:"duration"`  // microseconds
	References    []Reference `bson:"references"`
	ProcessID     string      `bson:"processID"`
	Process       Process     `bson:"process,omitempty"`
	Tags          []KeyValue  `bson:"tags"`
	Logs          []Log       `bson:"logs"`
	Warnings      []string    `bson:"warnings"`
}
```

> To meet our needs, we will write a parsed span to mongo too. You could remove it if don’t needed.

```go
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
    // !!!!!!!!!!!!!!!!!!!!!!!!  remove this  !!!!!!!!!!!!!!!!!!!!!!!!
	go receiver.WriteCustomSpan(ctx, span)
	return err
}
```

