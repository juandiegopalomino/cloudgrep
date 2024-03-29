package generator

import (
	"testing"

	"github.com/juandiegopalomino/cloudgrep/hack/awsgen/config"
	"github.com/juandiegopalomino/cloudgrep/hack/awsgen/writer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator(t *testing.T) {
	cfg := config.Config{
		Services: []config.Service{
			{
				Name:           "foo",
				ServicePackage: "foo2",
				Types: []config.Type{
					{
						Name: "Bar",
						ListAPI: config.ListAPI{
							Call: "ListFoo",
							InputOverrides: config.InputOverrides{
								FieldFuncs: map[string]string{
									"Foo": "fieldInputFoo",
								},
								FullFuncs: []string{
									"fullInputFoo",
								},
							},
							Pagination: true,
							OutputKey:  config.NestedField{config.Field{Name: "Spam"}, config.Field{Name: "Ham"}},
							SDKType:    "Foo",
							IDField: config.Field{
								Name:    "ID",
								Pointer: true,
							},
						},
						GetTagsAPI: config.GetTagsAPI{
							Call: "GetBarTags",
							InputIDField: config.Field{
								Name:      "BarID",
								SliceType: "types.BarIDType",
							},
							Tags: &config.TagField{
								Field: config.NestedField{
									config.Field{Name: "Tags"},
								},
								Key:     "Key",
								Value:   "Value",
								Style:   "struct",
								Pointer: true,
							},
							AllowedAPIErrorCodes: []string{
								"SpamError",
							},
						},
						Transformers: []config.Transformer{
							{
								Name: "foo",
								Expr: "bar",
							},
							{
								Name: "spam",
								Expr: "ham[%type]",
							},
							{
								Expr: "a",
							},
							{
								Expr: "b[%type]",
							},
						},
					},
				},
			},
		},
	}
	err := config.AggregateValidationErrors(cfg.Validate())
	require.NoError(t, err)

	w := writer.NewFakeWriter()
	g := Generator{Format: true}
	err = g.Generate(w, cfg)
	assert.NoError(t, err)
	assert.Len(t, w.Files, 2)
}
