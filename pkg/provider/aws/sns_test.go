package aws

import (
	"testing"

	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"github.com/juandiegopalomino/cloudgrep/pkg/testingutil"
	testprovider "github.com/juandiegopalomino/cloudgrep/pkg/testingutil/provider"
)

func TestFetchSnsTopic(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "sns.Topic")

	testingutil.AssertResourceCount(t, resources, "", 1)
	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:            "sns.Topic",
		Region:          defaultRegion,
		DisplayIdPrefix: "testing-",
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "sns-topic-0",
			},
		},
	})
}
