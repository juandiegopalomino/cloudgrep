package provider

import (
	"context"
	_ "embed"
	"reflect"
	"testing"

	"gorm.io/datatypes"

	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/provider/mapper"
	"github.com/run-x/cloudgrep/pkg/testingutil"
	"github.com/run-x/cloudgrep/pkg/util"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestMapper(t *testing.T) {

	ctx := context.Background()
	logger := zaptest.NewLogger(t)

	provider, err := NewTestProvider(ctx, config.Provider{}, logger)
	assert.NoError(t, err)

	//define some test data
	architecture := "x86_64"
	tr1 := TestResource{
		InstanceId:     "i-123",
		Architecture:   &architecture,
		SomeTags:       []TestTag{{Name: "enabled", Val: "true"}, {Name: "eks:nodegroup", Val: "staging-default"}},
		NeverReturned:  "should not see this",
		unexported:     "not exported",
		SecurityGroups: []string{"sg-1", "sg-2"},
		Limit:          Limit{CPU: MinMax{Min: 1, Max: 3}, Memory: MinMax{Min: 256, Max: 512}},
	}
	tr2 := TestResource{
		InstanceId:   "i-124",
		Architecture: nil,
		SomeTags:     []TestTag{},
	}

	r1 := model.Resource{
		Id: "i-123", Region: "us-east-1", Type: "test.Instance",
		Tags: []model.Tag{
			{Key: "enabled", Value: "true"},
			{Key: "eks:nodegroup", Value: "staging-default"}},
		RawData: datatypes.JSON([]byte(`{"InstanceId":"i-123","Architecture":"x86_64","SomeTags":[{"Name":"enabled","Val":"true"},{"Name":"eks:nodegroup","Val":"staging-default"}],"NeverReturned":"should not see this","SecurityGroups":["sg-1","sg-2"],"Limit":{"CPU":{"Min":1,"Max":3},"Memory":{"Min":256,"Max":512}}}`)),
	}
	r2 := model.Resource{
		Id: "i-124", Region: "us-east-1", Type: "test.Instance",
		Tags:    []model.Tag(nil),
		RawData: datatypes.JSON([]byte(`{"InstanceId":"i-124","Architecture":null,"SomeTags":[],"NeverReturned":"","SecurityGroups":null,"Limit":{"CPU":{"Min":0,"Max":0},"Memory":{"Min":0,"Max":0}}}`)),
	}

	// test conversion
	_r1, err := provider.GetMapper().ToResource(ctx, tr1, "us-east-1")
	assert.NoError(t, err)
	testingutil.AssertEqualsResource(t, r1, _r1)

	// fetch the resources
	resources, err := fetchResources(context.WithValue(ctx, Return("FetchTestResources"), []TestResource{tr1, tr2}), provider)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resources))
	// expect the resources object to be generated - the data is generated by FetchTestResources
	testingutil.AssertEqualsResource(t, r1, *resources[0])
	testingutil.AssertEqualsResource(t, r2, *resources[1])

}

//go:embed mapper_test.yaml
var embedConfig []byte

type TestProvider struct {
	*zap.Logger
	mapper.Mapper
}

type TestResource struct {
	InstanceId     string
	Architecture   *string
	SomeTags       []TestTag
	NeverReturned  string
	unexported     string
	SecurityGroups []string
	Limit          Limit
}
type TestTag struct {
	Name string
	Val  string
}

type Limit struct {
	CPU    MinMax
	Memory MinMax
}

type MinMax struct {
	Min int
	Max int
}

func NewTestProvider(ctx context.Context, cfg config.Provider, logger *zap.Logger) (Provider, error) {
	p := TestProvider{}
	p.Logger = logger
	var err error
	p.Mapper, err = mapper.New(embedConfig, logger, reflect.ValueOf(p))
	if err != nil {
		return nil, err
	}
	return p, nil
}
func (p TestProvider) GetMapper() mapper.Mapper {
	return p.Mapper
}

func (p TestProvider) Region() string {
	return "us-east-1"
}

type Return string

func (TestProvider) FetchTestResources(ctx context.Context, output chan<- TestResource) error {
	return util.SendAllFromSlice(ctx, output, ctx.Value(Return("FetchTestResources")).([]TestResource))
}
