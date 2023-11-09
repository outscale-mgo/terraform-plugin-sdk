package terraform

import (
	"strings"
	"testing"

	"github.com/outscale-mgo/terraform-plugin-sdk/internal/addrs"
)

func TestRefreshGraphBuilder_configOrphans(t *testing.T) {

	m := testModule(t, "refresh-config-orphan")

	state := MustShimLegacyState(&State{
		Modules: []*ModuleState{
			{
				Path: rootModulePath,
				Resources: map[string]*ResourceState{
					"test_object.foo.0": {
						Type: "test_object",
						Deposed: []*InstanceState{
							{
								ID: "foo",
							},
						},
					},
					"test_object.foo.1": {
						Type: "test_object",
						Deposed: []*InstanceState{
							{
								ID: "bar",
							},
						},
					},
					"test_object.foo.2": {
						Type: "test_object",
						Deposed: []*InstanceState{
							{
								ID: "baz",
							},
						},
					},
					"data.test_object.foo.0": {
						Type: "test_object",
						Deposed: []*InstanceState{ // NOTE: Real-world data resources don't get deposed
							{
								ID: "foo",
							},
						},
					},
					"data.test_object.foo.1": {
						Type: "test_object",
						Deposed: []*InstanceState{ // NOTE: Real-world data resources don't get deposed
							{
								ID: "bar",
							},
						},
					},
					"data.test_object.foo.2": {
						Type: "test_object",
						Deposed: []*InstanceState{ // NOTE: Real-world data resources don't get deposed
							{
								ID: "baz",
							},
						},
					},
				},
			},
		},
	})

	b := &RefreshGraphBuilder{
		Config:     m,
		State:      state,
		Components: simpleMockComponentFactory(),
		Schemas:    simpleTestSchemas(),
	}
	g, err := b.Build(addrs.RootModuleInstance)
	if err != nil {
		t.Fatalf("Error building graph: %s", err)
	}

	actual := strings.TrimSpace(g.StringWithNodeTypes())
	expected := strings.TrimSpace(`
data.test_object.foo[0] - *terraform.NodeRefreshableManagedResourceInstance
  provider.test - *terraform.NodeApplyableProvider
data.test_object.foo[0] (deposed 00000001) - *terraform.NodePlanDeposedResourceInstanceObject
  provider.test - *terraform.NodeApplyableProvider
data.test_object.foo[1] - *terraform.NodeRefreshableManagedResourceInstance
  provider.test - *terraform.NodeApplyableProvider
data.test_object.foo[1] (deposed 00000001) - *terraform.NodePlanDeposedResourceInstanceObject
  provider.test - *terraform.NodeApplyableProvider
data.test_object.foo[2] - *terraform.NodeRefreshableManagedResourceInstance
  provider.test - *terraform.NodeApplyableProvider
data.test_object.foo[2] (deposed 00000001) - *terraform.NodePlanDeposedResourceInstanceObject
  provider.test - *terraform.NodeApplyableProvider
provider.test - *terraform.NodeApplyableProvider
provider.test (close) - *terraform.graphNodeCloseProvider
  data.test_object.foo[0] - *terraform.NodeRefreshableManagedResourceInstance
  data.test_object.foo[0] (deposed 00000001) - *terraform.NodePlanDeposedResourceInstanceObject
  data.test_object.foo[1] - *terraform.NodeRefreshableManagedResourceInstance
  data.test_object.foo[1] (deposed 00000001) - *terraform.NodePlanDeposedResourceInstanceObject
  data.test_object.foo[2] - *terraform.NodeRefreshableManagedResourceInstance
  data.test_object.foo[2] (deposed 00000001) - *terraform.NodePlanDeposedResourceInstanceObject
  test_object.foo - *terraform.NodeRefreshableManagedResource
  test_object.foo[0] (deposed 00000001) - *terraform.NodePlanDeposedResourceInstanceObject
  test_object.foo[1] (deposed 00000001) - *terraform.NodePlanDeposedResourceInstanceObject
  test_object.foo[2] (deposed 00000001) - *terraform.NodePlanDeposedResourceInstanceObject
test_object.foo - *terraform.NodeRefreshableManagedResource
  provider.test - *terraform.NodeApplyableProvider
test_object.foo[0] (deposed 00000001) - *terraform.NodePlanDeposedResourceInstanceObject
  provider.test - *terraform.NodeApplyableProvider
test_object.foo[1] (deposed 00000001) - *terraform.NodePlanDeposedResourceInstanceObject
  provider.test - *terraform.NodeApplyableProvider
test_object.foo[2] (deposed 00000001) - *terraform.NodePlanDeposedResourceInstanceObject
  provider.test - *terraform.NodeApplyableProvider
`)
	if expected != actual {
		t.Fatalf("wrong result\n\ngot:\n%s\n\nwant:\n%s", actual, expected)
	}
}
