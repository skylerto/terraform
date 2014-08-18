package schema

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/terraform"
)

func TestResourceDataGet(t *testing.T) {
	cases := []struct {
		Schema map[string]*Schema
		State  *terraform.ResourceState
		Diff   *terraform.ResourceDiff
		Key    string
		Value  interface{}
	}{
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: &terraform.ResourceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "foo",
						RequiresNew: true,
					},
				},
			},

			Key: "availability_zone",

			Value: "foo",
		},

		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: &terraform.ResourceState{
				Attributes: map[string]string{
					"availability_zone": "bar",
				},
			},

			Diff: nil,

			Key: "availability_zone",

			Value: "bar",
		},

		{
			Schema: map[string]*Schema{
				"port": &Schema{
					Type:     TypeInt,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: &terraform.ResourceState{
				Attributes: map[string]string{
					"port": "80",
				},
			},

			Diff: nil,

			Key: "port",

			Value: 80,
		},

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: &terraform.ResourceState{
				Attributes: map[string]string{
					"ports.#": "3",
					"ports.0": "1",
					"ports.1": "2",
					"ports.2": "5",
				},
			},

			Key: "ports.1",

			Value: 2,
		},

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: &terraform.ResourceState{
				Attributes: map[string]string{
					"ports.#": "3",
					"ports.0": "1",
					"ports.1": "2",
					"ports.2": "5",
				},
			},

			Key: "ports.#",

			Value: 3,
		},

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: nil,

			Key: "ports.#",

			Value: 0,
		},

		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: &terraform.ResourceState{
				Attributes: map[string]string{
					"ports.#": "3",
					"ports.0": "1",
					"ports.1": "2",
					"ports.2": "5",
				},
			},

			Key: "ports",

			Value: []interface{}{1, 2, 5},
		},

		{
			Schema: map[string]*Schema{
				"ingress": &Schema{
					Type:     TypeList,
					Required: true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"from": &Schema{
								Type:     TypeInt,
								Required: true,
							},
						},
					},
				},
			},

			State: nil,

			Diff: &terraform.ResourceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ingress.#": &terraform.ResourceAttrDiff{
						Old: "",
						New: "1",
					},
					"ingress.0.from": &terraform.ResourceAttrDiff{
						Old: "",
						New: "8080",
					},
				},
			},

			Key: "ingress.0",

			Value: map[string]interface{}{
				"from": 8080,
			},
		},

		{
			Schema: map[string]*Schema{
				"ingress": &Schema{
					Type:     TypeList,
					Required: true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"from": &Schema{
								Type:     TypeInt,
								Required: true,
							},
						},
					},
				},
			},

			State: nil,

			Diff: &terraform.ResourceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ingress.#": &terraform.ResourceAttrDiff{
						Old: "",
						New: "1",
					},
					"ingress.0.from": &terraform.ResourceAttrDiff{
						Old: "",
						New: "8080",
					},
				},
			},

			Key: "ingress",

			Value: []interface{}{
				map[string]interface{}{
					"from": 8080,
				},
			},
		},

		// Computed get
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Computed: true,
				},
			},

			State: &terraform.ResourceState{
				Attributes: map[string]string{
					"availability_zone": "foo",
				},
			},

			Key: "availability_zone",

			Value: "foo",
		},

		// Full object
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: &terraform.ResourceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "foo",
						RequiresNew: true,
					},
				},
			},

			Key: "",

			Value: map[string]interface{}{
				"availability_zone": "foo",
			},
		},

		// List of maps
		{
			Schema: map[string]*Schema{
				"config_vars": &Schema{
					Type:     TypeList,
					Optional: true,
					Computed: true,
					Elem: &Schema{
						Type: TypeMap,
					},
				},
			},

			State: nil,

			Diff: &terraform.ResourceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"config_vars.#": &terraform.ResourceAttrDiff{
						Old: "0",
						New: "2",
					},
					"config_vars.0.foo": &terraform.ResourceAttrDiff{
						Old: "",
						New: "bar",
					},
					"config_vars.1.bar": &terraform.ResourceAttrDiff{
						Old: "",
						New: "baz",
					},
				},
			},

			Key: "config_vars",

			Value: []interface{}{
				map[string]interface{}{
					"foo": "bar",
				},
				map[string]interface{}{
					"bar": "baz",
				},
			},
		},

		// List of maps in state
		{
			Schema: map[string]*Schema{
				"config_vars": &Schema{
					Type:     TypeList,
					Optional: true,
					Computed: true,
					Elem: &Schema{
						Type: TypeMap,
					},
				},
			},

			State: &terraform.ResourceState{
				Attributes: map[string]string{
					"config_vars.#":     "2",
					"config_vars.0.foo": "baz",
					"config_vars.1.bar": "bar",
				},
			},

			Diff: nil,

			Key: "config_vars",

			Value: []interface{}{
				map[string]interface{}{
					"foo": "baz",
				},
				map[string]interface{}{
					"bar": "bar",
				},
			},
		},
	}

	for i, tc := range cases {
		d, err := schemaMap(tc.Schema).Data(tc.State, tc.Diff)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		v := d.Get(tc.Key)
		if !reflect.DeepEqual(v, tc.Value) {
			t.Fatalf("Bad: %d\n\n%#v", i, v)
		}
	}
}

func TestResourceDataGetChange(t *testing.T) {
	cases := []struct {
		Schema   map[string]*Schema
		State    *terraform.ResourceState
		Diff     *terraform.ResourceDiff
		Key      string
		OldValue interface{}
		NewValue interface{}
	}{
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: &terraform.ResourceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "foo",
						RequiresNew: true,
					},
				},
			},

			Key: "availability_zone",

			OldValue: nil,
			NewValue: "foo",
		},

		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: &terraform.ResourceState{
				Attributes: map[string]string{
					"availability_zone": "foo",
				},
			},

			Diff: &terraform.ResourceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "foo",
						RequiresNew: true,
					},
				},
			},

			Key: "availability_zone",

			OldValue: "foo",
			NewValue: "foo",
		},
	}

	for i, tc := range cases {
		d, err := schemaMap(tc.Schema).Data(tc.State, tc.Diff)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		o, n := d.GetChange(tc.Key)
		if !reflect.DeepEqual(o, tc.OldValue) {
			t.Fatalf("Old Bad: %d\n\n%#v", i, o)
		}
		if !reflect.DeepEqual(n, tc.NewValue) {
			t.Fatalf("New Bad: %d\n\n%#v", i, n)
		}
	}
}

func TestResourceDataHasChange(t *testing.T) {
	cases := []struct {
		Schema map[string]*Schema
		State  *terraform.ResourceState
		Diff   *terraform.ResourceDiff
		Key    string
		Change bool
	}{
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: &terraform.ResourceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "foo",
						RequiresNew: true,
					},
				},
			},

			Key: "availability_zone",

			Change: true,
		},

		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: &terraform.ResourceState{
				Attributes: map[string]string{
					"availability_zone": "foo",
				},
			},

			Diff: &terraform.ResourceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "foo",
						RequiresNew: true,
					},
				},
			},

			Key: "availability_zone",

			Change: false,
		},
	}

	for i, tc := range cases {
		d, err := schemaMap(tc.Schema).Data(tc.State, tc.Diff)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := d.HasChange(tc.Key)
		if actual != tc.Change {
			t.Fatalf("Bad: %d %#v", i, actual)
		}
	}
}

func TestResourceDataSet(t *testing.T) {
	cases := []struct {
		Schema   map[string]*Schema
		State    *terraform.ResourceState
		Diff     *terraform.ResourceDiff
		Key      string
		Value    interface{}
		Err      bool
		GetKey   string
		GetValue interface{}
	}{
		// Basic good
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: nil,

			Key:   "availability_zone",
			Value: "foo",

			GetKey:   "availability_zone",
			GetValue: "foo",
		},

		// Basic int
		{
			Schema: map[string]*Schema{
				"port": &Schema{
					Type:     TypeInt,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: nil,

			Key:   "port",
			Value: 80,

			GetKey:   "port",
			GetValue: 80,
		},

		// Invalid type
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: nil,

			Key:   "availability_zone",
			Value: 80,
			Err:   true,

			GetKey:   "availability_zone",
			GetValue: nil,
		},

		// List of primitives, set element
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Computed: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: &terraform.ResourceState{
				Attributes: map[string]string{
					"ports.#": "3",
					"ports.0": "1",
					"ports.1": "2",
					"ports.2": "5",
				},
			},

			Diff: nil,

			Key:   "ports.1",
			Value: 3,

			GetKey:   "ports",
			GetValue: []interface{}{1, 3, 5},
		},

		// List of primitives, set list
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Computed: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: nil,

			Diff: nil,

			Key:   "ports",
			Value: []int{1, 2, 5},

			GetKey:   "ports",
			GetValue: []interface{}{1, 2, 5},
		},

		// List of primitives, set list with error
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Computed: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: nil,

			Diff: nil,

			Key:   "ports",
			Value: []interface{}{1, "NOPE", 5},
			Err:   true,

			GetKey:   "ports",
			GetValue: []interface{}{},
		},

		// List of resource, set element
		{
			Schema: map[string]*Schema{
				"ingress": &Schema{
					Type:     TypeList,
					Computed: true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"from": &Schema{
								Type: TypeInt,
							},
						},
					},
				},
			},

			State: &terraform.ResourceState{
				Attributes: map[string]string{
					"ingress.#":      "2",
					"ingress.0.from": "80",
					"ingress.1.from": "8080",
				},
			},

			Diff: nil,

			Key:   "ingress.1.from",
			Value: 9000,

			GetKey: "ingress",
			GetValue: []interface{}{
				map[string]interface{}{
					"from": 80,
				},
				map[string]interface{}{
					"from": 9000,
				},
			},
		},

		// List of resource, set full resource element
		{
			Schema: map[string]*Schema{
				"ingress": &Schema{
					Type:     TypeList,
					Computed: true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"from": &Schema{
								Type: TypeInt,
							},
						},
					},
				},
			},

			State: &terraform.ResourceState{
				Attributes: map[string]string{
					"ingress.#":      "2",
					"ingress.0.from": "80",
					"ingress.1.from": "8080",
				},
			},

			Diff: nil,

			Key: "ingress.1",
			Value: map[string]interface{}{
				"from": 9000,
			},

			GetKey: "ingress",
			GetValue: []interface{}{
				map[string]interface{}{
					"from": 80,
				},
				map[string]interface{}{
					"from": 9000,
				},
			},
		},

		// List of resource, set full resource element, with error
		{
			Schema: map[string]*Schema{
				"ingress": &Schema{
					Type:     TypeList,
					Computed: true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"from": &Schema{
								Type: TypeInt,
							},
							"to": &Schema{
								Type: TypeInt,
							},
						},
					},
				},
			},

			State: &terraform.ResourceState{
				Attributes: map[string]string{
					"ingress.#":      "2",
					"ingress.0.from": "80",
					"ingress.0.to":   "10",
					"ingress.1.from": "8080",
					"ingress.1.to":   "8080",
				},
			},

			Diff: nil,

			Key: "ingress.1",
			Value: map[string]interface{}{
				"from": 9000,
				"to":   "bar",
			},
			Err: true,

			GetKey: "ingress",
			GetValue: []interface{}{
				map[string]interface{}{
					"from": 80,
					"to":   10,
				},
				map[string]interface{}{
					"from": 8080,
					"to":   8080,
				},
			},
		},

		// Set a list of maps
		{
			Schema: map[string]*Schema{
				"config_vars": &Schema{
					Type:     TypeList,
					Optional: true,
					Computed: true,
					Elem: &Schema{
						Type: TypeMap,
					},
				},
			},

			State: nil,

			Diff: nil,

			Key: "config_vars",
			Value: []interface{}{
				map[string]interface{}{
					"foo": "bar",
				},
				map[string]interface{}{
					"bar": "baz",
				},
			},
			Err: false,

			GetKey: "config_vars",
			GetValue: []interface{}{
				map[string]interface{}{
					"foo": "bar",
				},
				map[string]interface{}{
					"bar": "baz",
				},
			},
		},
	}

	for i, tc := range cases {
		d, err := schemaMap(tc.Schema).Data(tc.State, tc.Diff)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		err = d.Set(tc.Key, tc.Value)
		if (err != nil) != tc.Err {
			t.Fatalf("%d err: %s", i, err)
		}

		v := d.Get(tc.GetKey)
		if !reflect.DeepEqual(v, tc.GetValue) {
			t.Fatalf("Get Bad: %d\n\n%#v", i, v)
		}
	}
}

func TestResourceDataState(t *testing.T) {
	cases := []struct {
		Schema map[string]*Schema
		State  *terraform.ResourceState
		Diff   *terraform.ResourceDiff
		Set    map[string]interface{}
		Result *terraform.ResourceState
	}{
		// Basic primitive in diff
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: &terraform.ResourceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "foo",
						RequiresNew: true,
					},
				},
			},

			Result: &terraform.ResourceState{
				Attributes: map[string]string{
					"availability_zone": "foo",
				},
			},
		},

		// Basic primitive set override
		{
			Schema: map[string]*Schema{
				"availability_zone": &Schema{
					Type:     TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},
			},

			State: nil,

			Diff: &terraform.ResourceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"availability_zone": &terraform.ResourceAttrDiff{
						Old:         "",
						New:         "foo",
						RequiresNew: true,
					},
				},
			},

			Set: map[string]interface{}{
				"availability_zone": "bar",
			},

			Result: &terraform.ResourceState{
				Attributes: map[string]string{
					"availability_zone": "bar",
				},
			},
		},

		// List
		{
			Schema: map[string]*Schema{
				"ports": &Schema{
					Type:     TypeList,
					Required: true,
					Elem:     &Schema{Type: TypeInt},
				},
			},

			State: &terraform.ResourceState{
				Attributes: map[string]string{
					"ports.#": "1",
					"ports.0": "80",
				},
			},

			Diff: &terraform.ResourceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ports.#": &terraform.ResourceAttrDiff{
						Old: "1",
						New: "2",
					},
					"ports.1": &terraform.ResourceAttrDiff{
						Old: "",
						New: "100",
					},
				},
			},

			Result: &terraform.ResourceState{
				Attributes: map[string]string{
					"ports.#": "2",
					"ports.0": "80",
					"ports.1": "100",
				},
			},
		},

		// List of resources
		{
			Schema: map[string]*Schema{
				"ingress": &Schema{
					Type:     TypeList,
					Required: true,
					Elem: &Resource{
						Schema: map[string]*Schema{
							"from": &Schema{
								Type:     TypeInt,
								Required: true,
							},
						},
					},
				},
			},

			State: &terraform.ResourceState{
				Attributes: map[string]string{
					"ingress.#":      "1",
					"ingress.0.from": "80",
				},
			},

			Diff: &terraform.ResourceDiff{
				Attributes: map[string]*terraform.ResourceAttrDiff{
					"ingress.#": &terraform.ResourceAttrDiff{
						Old: "1",
						New: "2",
					},
					"ingress.0.from": &terraform.ResourceAttrDiff{
						Old: "80",
						New: "150",
					},
					"ingress.1.from": &terraform.ResourceAttrDiff{
						Old: "",
						New: "100",
					},
				},
			},

			Result: &terraform.ResourceState{
				Attributes: map[string]string{
					"ingress.#":      "2",
					"ingress.0.from": "150",
					"ingress.1.from": "100",
				},
			},
		},

		// List of maps
		{
			Schema: map[string]*Schema{
				"config_vars": &Schema{
					Type:     TypeList,
					Optional: true,
					Computed: true,
					Elem: &Schema{
						Type: TypeMap,
					},
				},
			},

			State: &terraform.ResourceState{
				Attributes: map[string]string{
					"config_vars.#":     "2",
					"config_vars.0.foo": "bar",
					"config_vars.1.bar": "baz",
				},
			},

			Set: map[string]interface{}{
				"config_vars.1": map[string]interface{}{
					"baz": "bang",
				},
			},

			Result: &terraform.ResourceState{
				Attributes: map[string]string{
					"config_vars.#":     "2",
					"config_vars.0.foo": "bar",
					"config_vars.1.baz": "bang",
				},
			},
		},
	}

	for i, tc := range cases {
		d, err := schemaMap(tc.Schema).Data(tc.State, tc.Diff)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		for k, v := range tc.Set {
			if err := d.Set(k, v); err != nil {
				t.Fatalf("%d err: %s", i, err)
			}
		}

		actual := d.State()
		if !reflect.DeepEqual(actual, tc.Result) {
			t.Fatalf("Bad: %d\n\n%#v", i, actual)
		}
	}
}

func TestResourceDataSetId(t *testing.T) {
	d := &ResourceData{}
	d.SetId("foo")

	actual := d.State()
	if actual.ID != "foo" {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestResourceDataSetId_clear(t *testing.T) {
	d := &ResourceData{
		state: &terraform.ResourceState{ID: "bar"},
	}
	d.SetId("")

	actual := d.State()
	if actual.ID != "" {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestResourceDataSetId_override(t *testing.T) {
	d := &ResourceData{
		state: &terraform.ResourceState{ID: "bar"},
	}
	d.SetId("foo")

	actual := d.State()
	if actual.ID != "foo" {
		t.Fatalf("bad: %#v", actual)
	}
}
