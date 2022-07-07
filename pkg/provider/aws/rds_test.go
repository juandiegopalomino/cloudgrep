package aws

import (
	"testing"

	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/testingutil"
	testprovider "github.com/run-x/cloudgrep/pkg/testingutil/provider"
)

func TestFetchRDSInstances(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "rds.DBInstance")

	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		AccountId: ctx.accountId,
		Type:      "rds.DBInstance",
		Region:    defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "rds-instance-0",
			},
		},
		RawData: map[string]any{
			"DBInstanceClass": "db.t3.micro",
			"Engine":          "postgres",
		},
	})
}

func TestFetchRDSClusters(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "rds.DBCluster")

	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		AccountId: ctx.accountId,
		Type:      "rds.DBCluster",
		Region:    defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "rds-cluster-0",
			},
		},
		RawData: map[string]any{
			"Engine": "aurora-postgresql",
		},
	})
}

func TestFetchRDSClusterSnapshots(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "rds.DBClusterSnapshot")

	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		AccountId: ctx.accountId,
		Type:      "rds.DBClusterSnapshot",
		Region:    defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "rds-cluster-snapshot-0",
			},
		},
	})
}

func TestFetchRDSSnapshots(t *testing.T) {
	t.Parallel()

	ctx := setupIntegrationTest(t)

	resources := testprovider.FetchResources(ctx.ctx, t, ctx.p, "rds.DBSnapshot")

	testingutil.AssertResourceFilteredCount(t, resources, 1, testingutil.ResourceFilter{
		AccountId: ctx.accountId,
		Type:      "rds.DBSnapshot",
		Region:    defaultRegion,
		Tags: model.Tags{
			{
				Key:   testingutil.TestTag,
				Value: "rds-snapshot-0",
			},
		},
	})
}
