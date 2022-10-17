package pattern_test

import (
	"testing"

	"github.com/xenomote/object_language/pattern"
)

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

func TestInterpret(t *testing.T) {
	tests := []struct {
		pattern     string
		input       string
		shouldMatch bool
		output      map[string]string
	}{
		{`{}`, `[]`, false, nil},
		{`[]`, `{}`, false, nil},

		{`{}`, `{}`, true, nil},
		{`{}`, `{"a": 1, "b": 2, "c": 3}`, true, nil},

		{`{"a": 1}`, `{}`, false, nil},
		{`{"a": 1}`, `{"b": 1}`, false, nil},
		{`{"a": 1}`, `{"a": 2}`, false, nil},
		{`{"a": 1}`, `{"a": 1}`, true, nil},

		{`{"a"?: 1}`, `{}`, true, nil},
		{`{"a"?: 1}`, `{"a": 1}`, true, nil},
		{`{"a"?: 1}`, `{"a": 2}`, false, nil},

		{`{"a": <=x>}`, `{"a": 1}`, true, map[string]string{"x": "1"}},
		{`{"a": <=x>, "b": <=y>}`, `{"a": 1, "b": 2}`, true, map[string]string{"x": "1", "y": "2"}},
		{`{"a": <=x>, "b": <=y>}`, `{"b": 2, "a": 1}`, true, map[string]string{"x": "1", "y": "2"}},
		{`{"b": <=y>, "a": <=x>}`, `{"a": 1, "b": 2}`, true, map[string]string{"x": "1", "y": "2"}},
		{`{"b": <=y>, "a": <=x>}`, `{"b": 2, "a": 1}`, true, map[string]string{"x": "1", "y": "2"}},

		{`{"a"?: <=x>}`, `{"a": 1}`, true, map[string]string{"x": "1"}},
		{`{"a"?: <=x>}`, `{}`, true, nil},

		{`[]`, `[]`, true, nil},
		{`[]`, `[1, 2, 3]`, true, nil},

		{`[1: 1]`, `[]`, false, nil},
		{`[1: 1]`, `[1]`, false, nil},
		{`[1: 1]`, `[1, 2]`, false, nil},
		{`[1: 1]`, `[2, 1]`, true, nil},

		{`[0?: 1]`, `[]`, true, nil},
		{`[0?: 1]`, `[1]`, true, nil},
		{`[0?: 1]`, `[2]`, false, nil},

		{`[1: <=x>]`, `[2, 1]`, true, map[string]string{"x": "1"}},
		{`[0: <=x>, 1: <=y>]`, `[1, 2]`, true, map[string]string{"x": "1", "y": "2"}},
		{`[1: <=y>, 0: <=x>]`, `[1, 2]`, true, map[string]string{"x": "1", "y": "2"}},

		{`[0?: <=x>]`, `[1]`, true, map[string]string{"x": "1"}},
		{`[0?: <=x>]`, `[]`, true, nil},

		{`{"a": <=x>, "b": <x>}`, `{"a": 1, "b": 1}`, true, map[string]string{"x": "1"}},
		{`{"a": <=x>, "b": <x>}`, `{"a": 1, "b": 2}`, false, nil},

		{`{"a": {"b": <=x>}, "c": <x>}`, `{"a": {"b": 1}, "c": 1}`, true, map[string]string{"x": "1"}},
		{`{"a": {"b": <=x>}, "c": <x>}`, `{"a": {"b": 1}, "c": 2}`, false, nil},

		{`[0: <=x>, 1: <x>]`, `[1, 1]`, true, map[string]string{"x": "1"}},
		{`[0: <=x>, 1: <x>]`, `[1, 2]`, false, nil},

		{`[0: [0: <=x>], 1: <x>]`, `[[1], 1]`, true, map[string]string{"x": "1"}},
		{`[0: [0: <=x>], 1: <x>]`, `[[1], 2]`, false, nil},
	}

	for _, test := range tests {
		name := test.pattern + " -> " + test.input
		t.Run(name, func(t *testing.T) {
			p, err := pattern.Parse(test.pattern)
			if err != nil {
				t.Fatal(err)
			}

			b, err := p.Interpret(test.input)
			if test.shouldMatch != (err == nil) {
				if test.shouldMatch {
					t.Fatalf("%s failed to match: %s", name, err)
				} else {
					t.Fatalf("%s should not have matched", name)
				}

			}

			for k := range test.output {
				if _, exists := b[k]; !exists {
					t.Errorf("output did not contain expected binding %s", k)
				}
			}

			for k := range b {
				if _, exists := test.output[k]; !exists {
					t.Errorf("output contained unexpected binding %s", k)
				}
			}

			for k, v1 := range test.output {
				if v2, exists := b[k]; exists && v1 != v2 {
					t.Errorf("output for %s did not match expected value: %s != %s", k, v2, v1)
				}
			}
		})
	}
}
