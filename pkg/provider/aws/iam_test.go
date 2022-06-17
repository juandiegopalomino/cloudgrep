package aws

import (
	"regexp"
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/testingutil"
	testprovider "github.com/run-x/cloudgrep/pkg/testingutil/provider"
	"github.com/stretchr/testify/assert"
)

func TestFetchOpenIDConnectProviders(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "iam.OpenIDConnectProvider")

	resources = testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "iam.OpenIDConnectProvider",
		Region: globalRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "iam-oidc-provider-eks-main",
			},
		},
	})

	if len(resources) < 1 {
		return
	}

	resource := resources[0]
	regxp := regexp.MustCompile(`oidc\.eks\.[a-z0-9-]+?\.amazonaws\.com`)
	assert.Regexp(t, regxp, resource.Id)
}

func TestFetchPolicies(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "iam.Policy")

	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "iam.Policy",
		Region: globalRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "iam-policy-0",
			},
		},
		RawData: map[string]any{
			"Path": "/test/",
		},
	})
}

func TestFetchRoles(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "iam.Role")

	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "iam.Role",
		Region: globalRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "iam-role-0",
			},
		},
		RawData: map[string]any{
			"Path": "/test/",
		},
	})
}

func TestFetchUsers(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "iam.User")

	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		Type:   "iam.User",
		Region: globalRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "iam-user-0",
			},
		},
		RawData: map[string]any{
			"Path": "/test/",
		},
	})
}
