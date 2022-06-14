package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/resourceconverter"
)

func (p *Provider) register_lambda(mapping map[string]mapper) {
	mapping["lambda.Function"] = mapper{
		FetchFunc: p.fetch_lambda_Function,
		IdField:   "FunctionArn",
		IsGlobal:  false,
	}
}

func (p *Provider) fetch_lambda_Function(ctx context.Context, output chan<- model.Resource) error {
	client := lambda.NewFromConfig(p.config)
	input := &lambda.ListFunctionsInput{}

	resourceConverter := p.converterFor("lambda.Function")
	paginator := lambda.NewListFunctionsPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)

		if err != nil {
			return fmt.Errorf("failed to fetch %s: %w", "lambda.Function", err)
		}

		if err := resourceconverter.SendAllConvertedTags(ctx, output, resourceConverter, page.Functions, p.getTags_lambda_Function); err != nil {
			return err
		}
	}

	return nil
}
func (p *Provider) getTags_lambda_Function(ctx context.Context, resource types.FunctionConfiguration) (model.Tags, error) {
	client := lambda.NewFromConfig(p.config)
	input := &lambda.GetFunctionInput{}

	input.FunctionName = resource.FunctionArn

	output, err := client.GetFunction(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s tags: %w", "lambda.Function", err)
	}
	tagField_0 := output.Tags

	var tags model.Tags

	for key, value := range tagField_0 {
		tags = append(tags, model.Tag{
			Key:   key,
			Value: value,
		})
	}

	return tags, nil
}