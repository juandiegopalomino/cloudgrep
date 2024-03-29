package datastore

import (
	"context"
	"sync"

	"github.com/juandiegopalomino/cloudgrep/pkg/datastore"
	"github.com/juandiegopalomino/cloudgrep/pkg/model"
)

// Datastore implementation that drops all inserted data, but keeps track of number of resources written (returned via Stats())
// Useful for tests, and can return a provided error on write
type Blackhole struct {
	l          sync.Mutex
	count      int
	writeError error
}

var _ datastore.Datastore = &Blackhole{}

func (s *Blackhole) Ping() error {
	return nil
}

func (s *Blackhole) GetResource(ctx context.Context, id string) (*model.Resource, error) {
	return nil, nil
}

func (s *Blackhole) GetResources(ctx context.Context, query []byte) (model.ResourcesResponse, error) {
	return model.ResourcesResponse{}, nil
}

func (s *Blackhole) WriteResources(ctx context.Context, resources model.Resources) error {
	s.l.Lock()
	defer s.l.Unlock()
	if s.writeError != nil {
		return s.writeError
	}

	s.count += len(resources)
	return nil
}

func (s *Blackhole) Stats(context.Context) (model.Stats, error) {
	return model.Stats{
		ResourcesCount: s.Count(),
	}, nil
}

func (s *Blackhole) GetFields(context.Context) (model.FieldGroups, error) {
	return nil, nil
}

func (s *Blackhole) WriteEvent(ctx context.Context, event model.Event) error {
	return nil
}

func (s *Blackhole) EngineStatus(ctx context.Context) (model.Event, error) {
	return model.Event{}, nil
}
