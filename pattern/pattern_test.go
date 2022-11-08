package pattern_test

import (
	"encoding/json"
	"testing"

	"github.com/xenomote/json_matcher/pattern"
)

func TestMatches(t *testing.T) {
	tests := []struct {
		a interface{}
		b interface{}
	}{
		{json.RawMessage(`{"x": 1, "y": 2}`), json.RawMessage(`{"y": 2, "x": 1}`)},
		{map[string]interface{}{"x": json.RawMessage("1"), "y": json.RawMessage("2")}, json.RawMessage(`{"y": 2, "x": 1}`)},
	}

	for _, test := range tests {
		if !pattern.Matches(test.a, test.b) {
			t.Errorf(`%s != %s`, test.a, test.b)
		}
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldParse bool
	}{
		{"empty object", `{}`, true},
		{"empty array", `[]`, true},

		{"object with one field", `{"a": 123}`, true},
		{"object with one binding", `{"a": <=x>}`, true},
		{"object with bound field", `{"a": <=x> 123}`, true},

		{"object with optional field", `{"a"?: 123}`, true},
		{"object with optional binding", `{"a"?: <=x>}`, true},
		{"object with optional bound field", `{"a"?: <=x> 123}`, true},

		{"object with reference", `{"a": <=x>, "b": <x>}`, true},
		{"object with self reference", `{"a": <=x> <x>}`, false},
		{"object with nested reference", `{"a": {"b": <=x>}, "c": <x>}`, true},
		{"object with nested self reference", `{"a": <=x> {"b": <x>}}`, false},

		{"object with nested object", `{"a": {}}`, true},
		{"object with nested array", `{"a": []}`, true},

		{"object with duplicate key", `{"a": 123, "a": 456}`, false},
		{"object with numeric key", `{1: 123}`, false},

		{"array with one field", `[0: 1]`, true},
		{"array with one binding", `[0: <=x>]`, true},
		{"array with bound field", `[0: <=x> 1]`, true},

		{"array with reference", `[0: <=x>, 1: <x>]`, true},
		{"array with self reference", `[0: <=x> <x>]`, false},
		{"array with nested reference", `[0: [0: <=x>], 1: <x>]`, true},
		{"array with nested self reference", `[0: <=x> [0: <x>]]`, false},

		{"array with optional field", `[0?: 1]`, true},
		{"array with optional binding", `[0?: <=x>]`, true},
		{"array with optional bound field", `[0?: <=x> 1]`, true},

		{"array with nested object", `[0: {}]`, true},
		{"array with nested array", `[0: []]`, true},

		{"array with duplicate index", `[0: 1, 0: 2]`, false},
		{"array with string index", `["a": 123]`, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := pattern.Parse(test.input)
			parsed := (err == nil)
			if parsed != test.shouldParse {
				if test.shouldParse {
					t.Fatalf("%s failed to parse: %s", test.name, err)
				} else {
					t.Fatalf("%s should not have parsed", test.name)
				}

			}
		})
	}
}

type bindings = map[string]interface{}

func TestInterpret(t *testing.T) {
	tests := []struct {
		pattern     string
		input       string
		shouldMatch bool
		output      string
	}{
		{`{}`, `[]`, false, ``},
		{`[]`, `{}`, false, ``},

		{`{}`, `{}`, true, `{}`},
		{`{}`, `{"a": 1, "b": 2, "c": 3}`, true, `{}`},

		{`{"a": 1}`, `{}`, false, ``},
		{`{"a": 1}`, `{"b": 1}`, false, ``},
		{`{"a": 1}`, `{"a": 2}`, false, ``},
		{`{"a": 1}`, `{"a": 1}`, true, `{}`},

		{`{"a"?: 1}`, `{}`, true, `{}`},
		{`{"a"?: 1}`, `{"a": 1}`, true, `{}`},
		{`{"a"?: 1}`, `{"a": 2}`, false, ``},

		{`{"a": <=x>}`, `{"a": 1}`, true, `{"x": 1}`},
		{`{"a": <=x>, "b": <=y>}`, `{"a": 1, "b": 2}`, true, `{"x": 1, "y": 2}`},
		{`{"a": <=x>, "b": <=y>}`, `{"b": 2, "a": 1}`, true, `{"x": 1, "y": 2}`},
		{`{"b": <=y>, "a": <=x>}`, `{"a": 1, "b": 2}`, true, `{"x": 1, "y": 2}`},
		{`{"b": <=y>, "a": <=x>}`, `{"b": 2, "a": 1}`, true, `{"x": 1, "y": 2}`},

		{`{"a"?: <=x>}`, `{"a": 1}`, true, `{"x": 1}`},
		{`{"a"?: <=x>}`, `{}`, true, `{}`},

		{`[]`, `[]`, true, `{}`},
		{`[]`, `[1, 2, 3]`, true, `{}`},

		{`[1: 1]`, `[]`, false, ``},
		{`[1: 1]`, `[1]`, false, ``},
		{`[1: 1]`, `[1, 2]`, false, ``},
		{`[1: 1]`, `[2, 1]`, true, `{}`},

		{`[0?: 1]`, `[]`, true, `{}`},
		{`[0?: 1]`, `[1]`, true, `{}`},
		{`[0?: 1]`, `[2]`, false, ``},

		{`[1: <=x>]`, `[2, 1]`, true, `{"x": 1}`},
		{`[0: <=x>, 1: <=y>]`, `[1, 2]`, true, `{"x": 1, "y": 2}`},
		{`[1: <=y>, 0: <=x>]`, `[1, 2]`, true, `{"x": 1, "y": 2}`},

		{`[0?: <=x>]`, `[1]`, true, `{"x": 1}`},
		{`[0?: <=x>]`, `[]`, true, `{}`},

		{`{"a": <=x>, "b": <x>}`, `{"a": 1, "b": 1}`, true, `{"x": 1}`},
		{`{"a": <=x>, "b": <x>}`, `{"a": {"x": 1, "y": 2}, "b": {"y": 2, "x": 1}}`, true, `{"x": {"x":1, "y":2}}`},
		{`{"a": <=x>, "b": <x>}`, `{"a": {"y": 2, "x": 1}, "b": {"x": 1, "y": 2}}`, true, `{"x": {"x":1, "y":2}}`},
		{`{"a": <=x>, "b": <x>}`, `{"a": 1, "b": 2}`, false, ``},

		{`{"a": {"b": <=x>}, "c": <x>}`, `{"a": {"b": 1}, "c": 1}`, true, `{"x": 1}`},
		{`{"a": {"b": <=x>}, "c": <x>}`, `{"a": {"b": 1}, "c": 2}`, false, ``},

		{`[0: <=x>, 1: <x>]`, `[1, 1]`, true, `{"x": 1}`},
		{`[0: <=x>, 1: <x>]`, `[1, 2]`, false, ``},

		{`[0: [0: <=x>], 1: <x>]`, `[[1], 1]`, true, `{"x": 1}`},
		{`[0: [0: <=x>], 1: <x>]`, `[[1], 2]`, false, ``},
	}

	for _, test := range tests {
		name := test.pattern + " -> " + test.input
		t.Run(name, func(t *testing.T) {
			p, err := pattern.Parse(test.pattern)
			if err != nil {
				t.Fatal(err)
			}

			b, err := p.Interpret(test.input)

			if test.shouldMatch && err != nil {
				t.Fatalf(`'%s' failed to match: %s`, name, err)
			}

			if !test.shouldMatch && err == nil {
				t.Fatalf(`'%s' should not have matched`, name)
			}

			if !test.shouldMatch {
				return
			}

			var a interface{}
			err = json.Unmarshal([]byte(test.output), &a)
			if err != nil {
				t.Fatalf(`bad test pattern '%s': %s`, test.output, err)
			}

			if !pattern.Matches(a, b) {
				ab, _ := json.Marshal(a)
				bb, _ := json.Marshal(b)

				t.Fatalf("'%s' did not match: \n'%s' != \n'%s'", name, string(ab), string(bb))
			}
		})
	}
}
