package diffing

import (
	mapset "github.com/deckarep/golang-set"
)

type Difference struct {
	Removed  []interface{}
	Added    []interface{}
	Remained []interface{}
}

func DetectDifferences(beforeSlice []interface{}, afterSlice []interface{}) Difference {
	beforeSet := mapset.NewSetFromSlice(beforeSlice)
	afterSet := mapset.NewSetFromSlice(afterSlice)
	return Difference{
		Removed:  beforeSet.Difference(afterSet).ToSlice(), // in before, not in after
		Added:    afterSet.Difference(beforeSet).ToSlice(), // in after, not in before
		Remained: beforeSet.Intersect(afterSet).ToSlice(),  // both in before and in after
	}
}
