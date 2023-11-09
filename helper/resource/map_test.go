package resource

import (
	"reflect"
	"testing"

	"github.com/outscale-mgo/terraform-plugin-sdk/internal/helper/config"
	"github.com/outscale-mgo/terraform-plugin-sdk/terraform"
)

func TestMapResources(t *testing.T) {
	m := &Map{
		Mapping: map[string]Resource{
			"aws_elb":      {},
			"aws_instance": {},
		},
	}

	rts := m.Resources()

	expected := []terraform.ResourceType{
		{
			Name: "aws_elb",
		},
		{
			Name: "aws_instance",
		},
	}

	if !reflect.DeepEqual(rts, expected) {
		t.Fatalf("bad: %#v", rts)
	}
}

func TestMapValidate(t *testing.T) {
	m := &Map{
		Mapping: map[string]Resource{
			"aws_elb": {
				ConfigValidator: &config.Validator{
					Required: []string{"foo"},
				},
			},
		},
	}

	var c *terraform.ResourceConfig
	var ws []string
	var es []error

	// Valid
	c = testConfigForMap(t, map[string]interface{}{"foo": "bar"})
	ws, es = m.Validate("aws_elb", c)
	if len(ws) > 0 {
		t.Fatalf("bad: %#v", ws)
	}
	if len(es) > 0 {
		t.Fatalf("bad: %#v", es)
	}

	// Invalid
	c = testConfigForMap(t, map[string]interface{}{})
	ws, es = m.Validate("aws_elb", c)
	if len(ws) > 0 {
		t.Fatalf("bad: %#v", ws)
	}
	if len(es) == 0 {
		t.Fatalf("bad: %#v", es)
	}
}

func testConfigForMap(t *testing.T, c map[string]interface{}) *terraform.ResourceConfig {
	return terraform.NewResourceConfigRaw(c)
}
