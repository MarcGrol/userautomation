package diffing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectDifferences(t *testing.T) {
	testCases := []struct {
		name      string
		before    []interface{}
		after     []interface{}
		removed   []interface{}
		added     []interface{}
		unchanged []interface{}
	}{
		{
			name:   "all nil",
			before: nil, after: nil,
			removed: empty(), added: empty(), unchanged: empty(),
		},
		{
			name:   "before nil",
			before: nil, after: slice("a", "b"),
			removed: empty(), added: slice("a", "b"), unchanged: empty(),
		},
		{
			name:   "after nil",
			before: slice("a", "b"), after: nil,
			removed: slice("a", "b"), added: empty(), unchanged: empty(),
		},
		{
			name:   "all empty",
			before: empty(), after: empty(),
			removed: empty(), added: empty(), unchanged: empty(),
		},
		{
			name:   "all same",
			before: slice("a", "b"), after: slice("a", "b"),
			removed: empty(), added: empty(), unchanged: slice("a", "b"),
		},
		{
			name:   "one removed",
			before: slice("a", "b"), after: slice("a"),
			removed: slice("b"), added: empty(), unchanged: slice("a"),
		},
		{
			name:   "all removed",
			before: slice("a", "b"), after: empty(),
			removed: slice("a", "b"), added: empty(), unchanged: empty(),
		},
		{
			name:   "one added",
			before: slice("a", "b"), after: slice("a", "b", "c"),
			removed: empty(), added: slice("c"), unchanged: slice("a", "b"),
		},
		{
			name:   "all added",
			before: empty(), after: slice("a", "b"),
			removed: empty(), added: slice("a", "b"), unchanged: empty(),
		},
		{
			name:   "one added, one removed",
			before: slice("a", "b"), after: slice("b", "c"),
			removed: slice("a"), added: slice("c"), unchanged: slice("b"),
		},
		{
			name:   "all added, all removed",
			before: slice("a", "b"), after: slice("c", "d"),
			removed: slice("a", "b"), added: slice("c", "d"), unchanged: empty(),
		},
		{
			name:   "mixed types",
			before: slice("a", "b"), after: []interface{}{"a", 1, 2},
			removed: slice("b"), added: []interface{}{1, 2}, unchanged: slice("a"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			diff := DetectDifferences(tc.before, tc.after)
			assert.Equal(t, tc.removed, diff.Removed)
			assert.Equal(t, tc.added, diff.Added)
			assert.Equal(t, tc.unchanged, diff.Unchanged)
		})
	}
}

func slice(varargs ...string) []interface{} {
	result := []interface{}{}
	for _, arg := range varargs {
		result = append(result, arg)
	}
	return result
}

func empty() []interface{} {
	return []interface{}{}
}
