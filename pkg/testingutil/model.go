package testingutil

import (
	"fmt"
	"strings"
	"testing"

	"github.com/juandiegopalomino/cloudgrep/pkg/model"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

const TestTag = "test"

// AssertResourceCount asserts that there is a specific number of given resources with the "test" tag.
// If tagValue is not an empty string, it also filters on resources that have the "test" tag with that value.
func AssertResourceCount(t TestingTB, resources []model.Resource, tagValue string, count int) {
	t.Helper()
	if tagValue == "" {
		resources = ResourceFilterTagKey(resources, TestTag)
	} else {
		resources = ResourceFilterTagKeyValue(resources, TestTag, tagValue)
	}

	assert.Lenf(t, resources, count, "expected %d resource(s) with tag %s=%s", count, TestTag, tagValue)
}

// ResourceFilterTagKey filters a slice of model.Resources based on a given tag key being present on that resource.
func ResourceFilterTagKey(in []model.Resource, key string) []model.Resource {
	return FilterFunc(in, func(r model.Resource) bool {
		for _, tag := range r.Tags {
			if tag.Key == key {
				return true
			}
		}

		return false
	})
}

// ResourceFilterTagKey filters a slice of model.Resources based on a given tag key/value pair being present on that resource.
func ResourceFilterTagKeyValue(in []model.Resource, key, value string) []model.Resource {
	return FilterFunc(in, func(r model.Resource) bool {
		for _, tag := range r.Tags {
			if tag.Key == key && tag.Value == value {
				return true
			}
		}

		return false
	})
}

func AssertEqualsResources(t *testing.T, a, b model.Resources) {
	assert.Equal(t, len(a), len(b))
	for _, resourceA := range a {
		resourceB := b.FindById(resourceA.Id)
		if resourceB == nil {
			t.Errorf("can't find a resource with id %v", resourceA.Id)
			return
		}
		AssertEqualsResource(t, *resourceA, *resourceB)
	}
}

func AssertEqualsResourcePter(t *testing.T, a, b *model.Resource) {
	AssertEqualsResource(t, *a, *b)
}

func AssertEqualsResource(t *testing.T, a, b model.Resource) {
	assert.Equal(t, a.Id, b.Id)
	assert.Equal(t, a.Region, b.Region)
	assert.Equal(t, a.Type, b.Type)
	jsonsEqual, err := JSONBytesEqual(a.RawData, b.RawData)
	assert.NoError(t, err)
	assert.True(t, jsonsEqual)
	assert.ElementsMatch(t, a.Tags.Clean(), b.Tags.Clean())
}

func AssertEqualsField(t *testing.T, a, b model.Field) {
	assert.Equal(t, a.Name, b.Name)
	assert.Equal(t, a.Count, b.Count)
	assert.ElementsMatch(t, a.Values, b.Values)
}

func AssertEqualsTag(t *testing.T, a, b *model.Tag) {
	if a == nil {
		assert.Nil(t, b)
		return
	}
	assert.Equal(t, a.Key, b.Key)
	assert.Equal(t, a.Value, b.Value)
}

func AssertEqualsTags(t *testing.T, a, b model.Tags) {
	assert.Equal(t, len(a), len(b))
	for _, tagA := range a {
		tagB := b.Find(tagA.Key)
		if tagB == nil {
			t.Errorf("can't find a tag with key %v", tagA.Key)
			return
		}
		AssertEqualsTag(t, &tagA, tagB)
	}
}

func AssertResourceFilteredCount(t TestingTB, resources []model.Resource, count int, filter ResourceFilter) []model.Resource {
	t.Helper()

	filtered := filter.Filter(resources)

	success := assert.Lenf(t, filtered, count, "expected %d resource(s) with filter %s", count, filter)
	if !success {
		partialFiltered := filter.PartialFilter(resources)

		names := maps.Keys(partialFiltered)
		slices.Sort(names)

		var matches []string
		for _, name := range names {
			resources := partialFiltered[name]
			matches = append(matches,
				fmt.Sprintf("%s=%d", name, len(resources)),
			)
		}

		t.Errorf("filter %s partial matches: %s", filter, strings.Join(matches, ", "))
	}
	return filtered
}

func AssertEvent(t *testing.T, ee, ae model.Event) {
	assert.Equal(t, ee.Type, ae.Type)
	assert.Equal(t, ee.Status, ae.Status)
	assert.Equal(t, ee.ProviderName, ae.ProviderName)
	assert.Equal(t, ee.ResourceType, ae.ResourceType)
	assert.Equal(t, ee.Error, ae.Error)
	if ee.ChildEvents != nil {
		assert.Equal(t, len(ee.ChildEvents), len(ae.ChildEvents))
		for _, e := range ee.ChildEvents {
			for _, a := range ae.ChildEvents {
				switch e.Type {
				case model.EventTypeEngine:
					if e.Type == a.Type {
						AssertEvent(t, e, a)
					}
				case model.EventTypeProvider:
					if e.Type == a.Type && e.ProviderName == a.ProviderName {
						AssertEvent(t, e, a)
					}
				case model.EventTypeResource:
					if e.Type == a.Type && e.ProviderName == a.ProviderName && e.ResourceType == a.ResourceType {
						AssertEvent(t, e, a)
					}
				}
			}
		}
	}
}
