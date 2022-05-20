package provider

import (
	"context"
	"testing"

	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func NewTestEngine(t *testing.T) Engine {
	ctx := context.Background()
	logger := zaptest.NewLogger(t)
	cfg := config.Config{
		// Datastore: datastoreConfig,
		Logging: config.Logging{
			Logger: logger,
			Mode:   "dev",
		},
	}
	datastore := datastore.NewMemoryStore(ctx, cfg)

	engine, err := NewEngine(ctx, cfg, datastore)
	assert.NoError(t, err)
	testProvider, err := NewTestProvider(ctx, config.Provider{}, logger)
	assert.NoError(t, err)
	engine.Providers = []Provider{testProvider}
	assert.NoError(t, err)
	return engine
}

func TestEngineRun(t *testing.T) {
	ctx := context.Background()
	//set some resources to return
	tr1 := TestResource{
		InstanceId:   "i-121",
		Architecture: nil,
		SomeTags:     []TestTag{},
	}
	tr2 := TestResource{
		InstanceId:   "i-122",
		Architecture: nil,
		SomeTags:     []TestTag{},
	}
	trMissingId := TestResource{
		Architecture: nil,
		SomeTags:     []TestTag{},
	}

	//run an engine that fetch 2 resources - no error
	engine := NewTestEngine(t)
	err := engine.Run(
		context.WithValue(ctx, Return("FetchTestResources"), []TestResource{tr1, tr2}),
	)
	assert.NoError(t, err)
	//check that the resources were stored
	resources, err := engine.GetResources(ctx, model.EmptyFilter())
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resources))

	//run an engine that fetch 2 resources - one with an error
	engine = NewTestEngine(t)
	ctx = context.WithValue(ctx, Return("FetchTestResources"), []TestResource{tr1, trMissingId})
	err = engine.Run(ctx)
	assert.ErrorContains(t, err, "could not find id field")
	//check that 1 resource was stored
	resources, err = engine.GetResources(ctx, model.EmptyFilter())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resources))
}