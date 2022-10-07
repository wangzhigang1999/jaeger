package mongo

import (
	"context"
	"flag"
	"fmt"
	"github.com/jaegertracing/jaeger/pkg/metrics"
	"github.com/jaegertracing/jaeger/storage/dependencystore"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Factory struct {
	operator         *SpanWriter
	options          Options
	mongoCli         *mongo.Client
	collection       *mongo.Collection
	collectionParsed *mongo.Collection
}

func (f *Factory) Initialize(metricsFactory metrics.Factory, logger *zap.Logger) error {
	cfg := f.options.Configuration
	var uri = fmt.Sprintf("mongodb://%s:%s@%s:%d", cfg.UserName, cfg.PassWord, cfg.Host, cfg.Port)
	newClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	collection := newClient.Database(cfg.Database).Collection(cfg.Collection)
	collectionParsed := newClient.Database(cfg.Database).Collection(cfg.Collection + "-parsed")

	f.collection = collection
	f.mongoCli = newClient
	f.collectionParsed = collectionParsed

	logger.Info("mongodb connected.")

	return nil
}

func NewFactory() *Factory {
	return &Factory{}
}

// InitFromViper implements plugin.Configurable
func (f *Factory) InitFromViper(v *viper.Viper, logger *zap.Logger) {
	f.options.InitFromViper(v)
}

func (f Factory) AddFlags(flagSet *flag.FlagSet) {
	f.options.AddFlags(flagSet)
}

// InitFromOptions initializes factory from the supplied options
func (f *Factory) InitFromOptions(opts Options) {
	f.options = opts
}

func (f Factory) CreateSpanWriter() (spanstore.Writer, error) {
	return SpanWriter{mongoClient: f.mongoCli, collection: f.collection, collectionParsed: f.collectionParsed, output: f.options.Configuration.Output}, nil
}

func (f *Factory) CreateSpanReader() (spanstore.Reader, error) {
	return &SpanReader{mongoClient: f.mongoCli, collection: f.collection}, nil
}

func (f *Factory) CreateDependencyReader() (dependencystore.Reader, error) {
	return &SpanReader{mongoClient: f.mongoCli, collection: f.collection}, nil
}
