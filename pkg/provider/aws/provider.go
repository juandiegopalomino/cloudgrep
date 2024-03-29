package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cfg "github.com/juandiegopalomino/cloudgrep/pkg/config"
	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	regionutil "github.com/juandiegopalomino/cloudgrep/pkg/provider/aws/regions"
	awsutil "github.com/juandiegopalomino/cloudgrep/pkg/provider/aws/util"
	"github.com/juandiegopalomino/cloudgrep/pkg/provider/types"
	"github.com/juandiegopalomino/cloudgrep/pkg/resourceconverter"
	_ "github.com/juandiegopalomino/cloudgrep/pkg/util/rlimit"
	"go.uber.org/zap"
)

type Provider struct {
	config    aws.Config
	accountId string
	region    regionutil.Region
}

func (p Provider) String() string {
	return fmt.Sprintf("AWS Provider for account %v, region %v", p.accountId, p.region.ID())
}

func (p Provider) AccountId() string {
	return p.accountId
}

func (p Provider) FetchFunctions() map[string]types.FetchFunc {
	funcMap := make(map[string]types.FetchFunc)
	for resourceType, mapping := range p.getTypeMapping() {
		if p.region.IsGlobal() != mapping.IsGlobal {
			continue
		}

		if mapping.ServiceEndpointID != "" && !p.region.IsServiceSupported(mapping.ServiceEndpointID) {
			continue
		}

		funcMap[resourceType] = mapping.FetchFunc
	}
	return funcMap
}

func (p *Provider) converterFor(resourceType string) resourceconverter.ResourceConverter {
	mapping, ok := p.getTypeMapping()[resourceType]
	if !ok {
		panic(fmt.Sprintf("Could not find mapping for resource type %v", resourceType))
	}

	region := p.region.ID()
	factory := func() model.Resource {
		return model.Resource{
			AccountId: p.accountId,
			Region:    region,
			Type:      resourceType,
		}
	}

	if mapping.UseMapConverter {
		return &resourceconverter.MapConverter{
			ResourceFactory: factory,
			TagField:        mapping.TagField,
			IdField:         mapping.IdField,
			DisplayIdField:  mapping.DisplayIDField,
		}
	}
	return &resourceconverter.ReflectionConverter{
		ResourceFactory: factory,
		TagField:        mapping.TagField,
		IdField:         mapping.IdField,
		DisplayIdField:  mapping.DisplayIDField,
	}
}

func NewProviders(ctx context.Context, cfg cfg.Provider, logger *zap.Logger) ([]types.Provider, error) {
	logger.Info("Connecting to AWS account")
	if cfg.Profile != "" {
		logger.Sugar().Infof("Using AWS profile '%v'", cfg.Profile)
	}
	defaultConfig, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(cfg.Profile),
		config.WithDefaultsMode(aws.DefaultsModeCrossRegion),
	)
	if err != nil {
		return nil, err
	}

	identity, err := awsutil.VerifyCreds(ctx, defaultConfig)
	if err != nil {
		if err.Error() == "no AWS credentials found" {
			err = fmt.Errorf("%w\nPlease set your AWS credentials using this guide: https://docs.aws.amazon.com/sdk-for-java/v1/developer-guide/setup-credentials.html", err)
		}
		return nil, err
	}

	regions, err := regionutil.SelectRegions(ctx, cfg.Regions, defaultConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot select regions for AWS provider: %w", err)
	}

	regionutil.SetConfigRegion(&defaultConfig, regions)
	logger.Sugar().Infof("Using the following identity: %v", *identity.Arn)
	logger.Sugar().Infof("Will look in regions %v", regions)
	var providers []types.Provider

	for _, region := range regions {
		newConfig := defaultConfig.Copy()
		if !region.IsGlobal() {
			newConfig.Region = region.ID()
		}

		logger.Sugar().Infof("Creating provider for AWS region %v", region)
		newProvider := Provider{
			config:    newConfig,
			accountId: *identity.Account,
			region:    region,
		}
		providers = append(providers, newProvider)
	}
	return providers, nil
}
