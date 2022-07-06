package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

func (p *Provider) register_autoscaling(mapping map[string]mapper) {
	mapping["autoscaling.AutoScalingGroup"] = mapper{
		ServiceEndpointID: "autoscaling",
		FetchFunc:         p.fetch_autoscaling_AutoScalingGroup,
		IdField:           "AutoScalingGroupName",
		IsGlobal:          false,
		TagField: resourceconverter.TagField{
			Name:  "Tags",
			Key:   "Key",
			Value: "Value",
		},
	}
}

func (p *Provider) fetch_autoscaling_AutoScalingGroup(ctx context.Context, output chan<- model.Resource) error {
	client := autoscaling.NewFromConfig(p.config)
	input := &autoscaling.DescribeAutoScalingGroupsInput{}

	resourceConverter := p.converterFor("autoscaling.AutoScalingGroup")
	commonTransformers := p.baseTransformers("autoscaling.AutoScalingGroup")
	transformers := append(
		resourceconverter.AllToGeneric[types.AutoScalingGroup](commonTransformers...),
		resourceconverter.WithConverter[types.AutoScalingGroup](resourceConverter),
	)
	paginator := autoscaling.NewDescribeAutoScalingGroupsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "autoscaling.AutoScalingGroup", err)
		}

		if err := resourceconverter.SendAll(ctx, output, page.AutoScalingGroups, transformers...); err != nil {
			return err
		}
	}

	return nil
}
