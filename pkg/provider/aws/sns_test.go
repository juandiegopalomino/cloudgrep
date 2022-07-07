package aws

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/testingutil"
	testprovider "github.com/run-x/cloudgrep/pkg/testingutil/provider"
)

func TestFetchSnsTopic(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "sns.SNS")

	testingutil.AssertResourceCount(t, resources, "", 1)
	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		AccountId: ctx.accountId,
		Type:      "sns.SNS",
		Region:    defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "sns-topic-0",
			},
		},
	})
}
