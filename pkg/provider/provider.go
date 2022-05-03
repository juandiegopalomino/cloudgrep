package provider

import (
	"context"
	"fmt"
	"reflect"

	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/provider/aws"
	"github.com/run-x/cloudgrep/pkg/provider/mapper"
	"go.uber.org/zap"
)

//Provider is an interface to be implemented for a cloud provider to fetch resources
//The provider must provide a mapping configuration which references the methods to fetch the resources.
//These methods need to be implemented and they will be called by a Mapper using reflection.
type Provider interface {
	//Init intializes the provider and return the mapping configuration
	Init(context.Context, config.Provider, *zap.Logger) (mapper.Config, error)
	Region() string
}

//Run the providers: fetches data about cloud resources and save them to store
func Run(ctx context.Context, cfg config.Config, datastore datastore.Datastore) error {

	//TODO use go routine to start the provider, review error handling to continue on error
	for _, config := range cfg.Providers {
		// create a provider
		provider, err := NewProvider(ctx, config)
		if err != nil {
			return err
		}
		mapperConfig, err := provider.Init(ctx, config, cfg.Logging.Logger)
		if err != nil {
			return err
		}

		//create a mapper
		mapper, err := mapper.New(mapperConfig, *cfg.Logging.Logger, reflect.ValueOf(provider))
		if err != nil {
			return err
		}

		// fetch the resources
		resources, err := FetchResources(ctx, provider, mapper)
		if err != nil {
			return err
		}
		// save to store
		err = datastore.WriteResources(ctx, resources)
		if err != nil {
			return err
		}
	}

	return nil
}

func FetchResources(ctx context.Context, provider Provider, mapper mapper.Mapper) ([]*model.Resource, error) {

	//load resources for each mapping
	var resources []*model.Resource
	for _, mapping := range mapper.Mappings {
		newResources, err := mapper.FetchResources(ctx, mapping, reflect.ValueOf(provider), provider.Region())
		if err != nil {
			return []*model.Resource{}, err
		}
		resources = append(resources, newResources...)
	}

	return resources, nil
}

func NewProvider(ctx context.Context, config config.Provider) (Provider, error) {

	var provider Provider
	if config.Cloud == "aws" {
		provider = &aws.AWSProvider{}
	}
	if provider == nil {
		return nil, fmt.Errorf("unknown provider cloud '%v'", config.Cloud)
	}
	return provider, nil
}
