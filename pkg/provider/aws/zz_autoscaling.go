package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

func (p *Provider) register_autoscaling(mapping map[string]mapper) {
	mapping["autoscaling.AutoScalingGroup"] = mapper{
		FetchFunc: p.fetch_autoscaling_AutoScalingGroup,
		IdField:   "AutoScalingGroupName",
		IsGlobal:  false,
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
	paginator := autoscaling.NewDescribeAutoScalingGroupsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "autoscaling.AutoScalingGroup", err)
		}

		if err := resourceconverter.SendAllConverted(ctx, output, resourceConverter, page.AutoScalingGroups); err != nil {
			return err
		}
	}

	return nil
}
