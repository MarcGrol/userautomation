package diffing

import (
	"fmt"
	"sort"
	"strings"

	mapset "github.com/deckarep/golang-set"
)

type Difference struct {
	Removed   []interface{}
	Added     []interface{}
	Unchanged []interface{}
}

func DetectDifferences(beforeSlice []interface{}, afterSlice []interface{}) Difference {
	beforeSet := mapset.NewSetFromSlice(beforeSlice)
	afterSet := mapset.NewSetFromSlice(afterSlice)
	return Difference{
		Removed:   sortit(beforeSet.Difference(afterSet).ToSlice()), // in before, not in after
		Added:     sortit(afterSet.Difference(beforeSet).ToSlice()), // in after, not in before
		Unchanged: sortit(afterSet.Intersect(beforeSet).ToSlice()),  // both in before and in after
	}
}

// We sort to make the result has predefined order for testability
func sortit(slice []interface{}) []interface{} {
	sort.Slice(slice, func(i, j int) bool {
		return strings.Compare(fmt.Sprintf("%s", slice[i]), fmt.Sprintf("%s", slice[j])) < 0
	})
	return slice
}
