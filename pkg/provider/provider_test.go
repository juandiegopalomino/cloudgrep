package provider

import (
	"context"
	"embed"
	"reflect"
	"testing"

	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/model"
	"github.com/run-x/cloudgrep/pkg/provider/mapper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestMapper(t *testing.T) {

	ctx := context.Background()
	logger := zaptest.NewLogger(t)

	provider, err := NewTestProvider(ctx, config.Provider{}, logger)
	if err != nil {
		t.Error(err)
	}

	//create a mapper
	mapper, err := mapper.New(provider.GetMapperConfig(), *logger, reflect.ValueOf(provider))
	if err != nil {
		t.Error(err)
	}

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
		Properties: []model.Property{
			{Name: "InstanceId", Value: "i-123"},
			{Name: "Architecture", Value: "x86_64"},
			{Name: "SecurityGroups[0]", Value: "sg-1"},
			{Name: "SecurityGroups[1]", Value: "sg-2"},
			{Name: "Limit[CPU][Min]", Value: "1"},
			{Name: "Limit[CPU][Max]", Value: "3"},
			{Name: "Limit[Memory][Min]", Value: "256"},
			{Name: "Limit[Memory][Max]", Value: "512"},
		},
	}
	r2 := model.Resource{
		Id: "i-124", Region: "us-east-1", Type: "test.Instance",
		Tags: []model.Tag(nil),
		Properties: []model.Property{
			{Name: "InstanceId", Value: "i-124"},
			{Name: "Architecture", Value: ""},
			{Name: "SecurityGroups", Value: ""},
			{Name: "Limit[CPU][Min]", Value: "0"},
			{Name: "Limit[CPU][Max]", Value: "0"},
			{Name: "Limit[Memory][Min]", Value: "0"},
			{Name: "Limit[Memory][Max]", Value: "0"},
		},
	}

	// test conversion
	_r1, err := mapper.ToRessource(tr1, "us-east-1")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, r1, _r1)

	// fetch the resources
	resources, err := FetchResources(context.WithValue(ctx, Return("FetchTestResources"), []TestResource{tr1, tr2}), provider, mapper)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 2, len(resources))
	// expect the resources object to be generated - the data is generated by FetchTestResources
	assert.Equal(t, r1, *resources[0])
	assert.Equal(t, r2, *resources[1])

}

//go:embed mapper_test.yaml
var embedConfig embed.FS

type TestProvider struct {
	*zap.Logger
	mapper.Config
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
	data, err := embedConfig.ReadFile("mapper_test.yaml")
	if err != nil {
		return nil, err
	}
	var config mapper.Config
	config, err = mapper.LoadConfig(data)
	if err != nil {
		return nil, err
	}
	p.Config = config
	return p, nil
}
func (p TestProvider) GetMapperConfig() mapper.Config {
	return p.Config
}

func (p TestProvider) Region() string {
	return "us-east-1"
}

type Return string

func (TestProvider) FetchTestResources(ctx context.Context) ([]TestResource, error) {
	return ctx.Value(Return("FetchTestResources")).([]TestResource), nil
}
