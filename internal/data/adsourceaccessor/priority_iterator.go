// thread-safeless

package adsourceaccessor

import (
	"geniusrabbit.dev/sspserver/internal/adsource"
)

type priorityIterator struct {
	index   int
	request *adsource.BidRequest
	sources []adsource.Source
}

// NewPriorityIterator from request and source
func NewPriorityIterator(request *adsource.BidRequest, sources []adsource.Source) adsource.SourceIterator {
	return &priorityIterator{
		index:   0,
		request: request,
		sources: sources,
	}
}

func (iter *priorityIterator) Next() (src adsource.Source) {
	if iter.index < len(iter.sources) {
		src = iter.sources[iter.index]
		iter.index++
	}
	return src
}

var _ adsource.SourceIterator = &priorityIterator{}
