package diffing

import (
	mapset "github.com/deckarep/golang-set"
)

type Difference struct {
	Removed []interface{}
	Added   []interface{}
	InBoth  []interface{}
}

func DetectDifferences(beforeSlice []interface{}, afterSlice []interface{}) Difference {
	beforeSet := mapset.NewSetFromSlice(beforeSlice)
	afterSet := mapset.NewSetFromSlice(afterSlice)
	return Difference{
		Removed: beforeSet.Difference(afterSet).ToSlice(), // in before, not in after
		Added:   afterSet.Difference(beforeSet).ToSlice(), // in after, not in before
		InBoth:  afterSet.Intersect(beforeSet).ToSlice(),  // both in before and in after
	}
}
