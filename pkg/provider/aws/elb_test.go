package aws

import (
	"testing"

	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"github.com/juandiegopalomino/cloudgrep/pkg/testingutil"
	testprovider "github.com/juandiegopalomino/cloudgrep/pkg/testingutil/provider"
)

func TestFetchLoadBalancers(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "elb.LoadBalancer")

	testingutil.AssertResourceCount(t, resources, "", 2)
	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "elb.LoadBalancer",
		Region: defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "elb-alb-1",
			},
		},
		RawData: map[string]any{
			"Scheme": "internal",
		},
	})
}
