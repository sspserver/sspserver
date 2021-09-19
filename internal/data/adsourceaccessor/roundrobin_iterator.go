// thread-safeless

package adsourceaccessor

import (
	"math/rand"

	"geniusrabbit.dev/sspserver/internal/adsource"
)

type roundrobinIterator struct {
	index    int
	endIndex int
	request  *adsource.BidRequest
	sources  []adsource.Source
}

// NewRoundrobinIterator from request and source
func NewRoundrobinIterator(request *adsource.BidRequest, sources []adsource.Source) adsource.SourceIterator {
	startIndex := rand.Int() % len(sources)
	return &roundrobinIterator{
		index:    startIndex,
		endIndex: startIndex,
		request:  request,
		sources:  sources,
	}
}

func (iter *roundrobinIterator) Next() adsource.Source {
	if iter.index >= len(iter.sources) {
		return nil
	}
	src := iter.sources[iter.index]
	if iter.index++; iter.index > len(iter.sources) {
		iter.index = 0
	}
	if iter.index == iter.endIndex {
		return nil
	}
	return src
}

var _ adsource.SourceIterator = &roundrobinIterator{}
