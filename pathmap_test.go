package pathmap_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/raphaelreyna/pathmap"
)

func TestGetAndSet(t *testing.T) {
	type setArgs struct {
		path  string
		value any
	}

	tests := []struct {
		name         string
		setArgsSlice []setArgs
		expected     pathmap.Map
	}{
		{
			name: "simple set and get",
			setArgsSlice: []setArgs{
				{"a", 1},
				{"b", 2},
			},
			expected: pathmap.Map{
				"a": 1,
				"b": 2,
			},
		},
		{
			name: "nested set and get",
			setArgsSlice: []setArgs{
				{"a.b.c", 3},
				{"a.b.d", 4},
				{"a.e", 5},
			},
			expected: pathmap.Map{
				"a": pathmap.Map{
					"b": pathmap.Map{
						"c": 3,
						"d": 4,
					},
					"e": 5,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := pathmap.Map{}
			for _, args := range test.setArgsSlice {
				m.Set(args.path, args.value)
			}
			for k, v := range test.expected {
				if got, ok := m[k]; !ok || !reflect.DeepEqual(got, v) {
					t.Errorf("m[%q] = %v, want %v", k, got, v)
				}
			}
		})
	}
}

func TestFlattenedKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    pathmap.Map
		expected []string
	}{
		{
			name: "simple keys",
			input: pathmap.Map{
				"a": 1,
				"b": 2,
			},
			expected: []string{"a", "b"},
		},
		{
			name: "nested keys",
			input: pathmap.Map{
				"a": pathmap.Map{
					"b": pathmap.Map{
						"c": 3,
						"d": 4,
					},
					"e": 5,
				},
			},
			expected: []string{"a.b.c", "a.b.d", "a.e"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.FlattenedKeys()
			sort.Strings(got)
			sort.Strings(test.expected)
			if !reflect.DeepEqual(got, test.expected) {
				t.Errorf("FlattenedKeys() = %v, want %v", got, test.expected)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name     string
		input    pathmap.Map
		path     string
		expected pathmap.Map
	}{
		{
			name: "delete simple key",
			input: pathmap.Map{
				"a": 1,
				"b": 2,
			},
			path: "a",
			expected: pathmap.Map{
				"b": 2,
			},
		},
		{
			name: "delete nested key",
			input: pathmap.Map{
				"a": pathmap.Map{
					"b": pathmap.Map{
						"c": 3,
						"d": 4,
					},
					"e": 5,
				},
			},
			path: "a.b.c",
			expected: pathmap.Map{
				"a": pathmap.Map{
					"b": pathmap.Map{
						"d": 4,
					},
					"e": 5,
				},
			},
		},
		{
			name: "delete non-existent key",
			input: pathmap.Map{
				"a": pathmap.Map{
					"b": pathmap.Map{
						"c": 3,
						"d": 4,
					},
					"e": 5,
				},
			},
			path: "a.b.x",
			expected: pathmap.Map{
				"a": pathmap.Map{
					"b": pathmap.Map{
						"c": 3,
						"d": 4,
					},
					"e": 5,
				},
			},
		},
		{
			name: "delete root key",
			input: pathmap.Map{
				"a": pathmap.Map{
					"b": pathmap.Map{
						"c": 3,
						"d": 4,
					},
					"e": 5,
				},
				"f": 6,
			},
			path: "a",
			expected: pathmap.Map{
				"f": 6,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m := test.input
			m.Delete(test.path)
			if !reflect.DeepEqual(m, test.expected) {
				t.Errorf("Delete(%q) = %v, want %v", test.path, m, test.expected)
			}
		})
	}
}
