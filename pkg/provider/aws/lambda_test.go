package aws

import (
	"testing"

	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"github.com/juandiegopalomino/cloudgrep/pkg/testingutil"
	testprovider "github.com/juandiegopalomino/cloudgrep/pkg/testingutil/provider"
)

func TestFetchFunctions(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "lambda.Function")

	testingutil.AssertResourceCount(t, resources, "", 2)
	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:            "lambda.Function",
		Region:          defaultRegion,
		DisplayIdPrefix: "testing-",
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "lambda-function-0",
			},
		},
		RawData: map[string]any{
			"Runtime": "nodejs16.x",
		},
	})
}
