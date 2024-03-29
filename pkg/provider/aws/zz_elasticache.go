package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/service/elasticache/types"

	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"github.com/juandiegopalomino/cloudgrep/pkg/resourceconverter"
)

func (p *Provider) registerElasticache(mapping map[string]mapper) {
	mapping["elasticache.CacheCluster"] = mapper{
		ServiceEndpointID: "elasticache",
		FetchFunc:         p.fetchElasticacheCacheCluster,
		IdField:           "ARN",
		IsGlobal:          false,
	}
}

func (p *Provider) fetchElasticacheCacheCluster(ctx context.Context, output chan<- model.Resource) error {
	client := elasticache.NewFromConfig(p.config)
	input := &elasticache.DescribeCacheClustersInput{}

	resourceConverter := p.converterFor("elasticache.CacheCluster")
	var transformers resourceconverter.Transformers[types.CacheCluster]
	transformers.AddNamed("tags", resourceconverter.TagTransformer(p.getTagsElasticacheCacheCluster))
	transformers.AddResource(displayIdArnPrefix("cluster:"))
	paginator := elasticache.NewDescribeCacheClustersPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "elasticache.CacheCluster", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.CacheClusters, transformers); err != nil {
			return err
		}
	}

	return nil
}
func (p *Provider) getTagsElasticacheCacheCluster(ctx context.Context, resource types.CacheCluster) (model.Tags, error) {
	client := elasticache.NewFromConfig(p.config)
	input := &elasticache.ListTagsForResourceInput{}

	input.ResourceName = resource.ARN

	output, err := client.ListTagsForResource(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s tags: %w", "elasticache.CacheCluster", err)
	}
	tagField_0 := output.TagList

	var tags model.Tags

	for _, field := range tagField_0 {
		tags = append(tags, model.Tag{
			Key:   *field.Key,
			Value: *field.Value,
		})
	}

	return tags, nil
}
