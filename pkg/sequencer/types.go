package sequencer

import (
	"context"

	"github.com/juandiegopalomino/cloudgrep/pkg/datastore"
	"github.com/juandiegopalomino/cloudgrep/pkg/provider"
)

type Sequencer interface {
	Run(ctx context.Context, ds datastore.Datastore, providers []provider.Provider) error
}
