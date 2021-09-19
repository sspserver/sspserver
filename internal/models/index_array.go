package models

import "github.com/geniusrabbit/gosql"

func indexSliceCmp(a1, a2 gosql.NullableOrderedUintArray) int {
	l1, l2 := a1.Len(), a2.Len()
	if l1 == 0 || l2 == 0 {
		switch {
		case l1 == l2:
			return 0
		case l1 == 0:
			return 1
		}
		return -1
	}
	a, b := a1[0], a2[0]
	switch {
	case a < b:
		return -1
	case b < a:
		return 1
	}
	return 0
}

///////////////////////////////////////////////////////////////////////////////
// Methods for campare object with filter
///////////////////////////////////////////////////////////////////////////////

// If slice less then val -1 ; 0 ; 1
func indexSliceOneCmp(s gosql.NullableOrderedUintArray, v uint) int {
	if s.Len() > 0 && s[s.Len()-1] < v {
		return -1
	}
	return 0
}

// If slice less then val -1 ; 0 ; 1
func indexStringSliceOneCmp(s gosql.NullableStringArray, v string) int {
	if s.Len() > 0 && s[s.Len()-1] < v {
		return -1
	}
	return 0
}
