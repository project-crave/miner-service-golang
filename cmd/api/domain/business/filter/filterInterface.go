package business

import (
	craveModel "crave/shared/model"
	"sync"
)

type IBusiness interface {
	GetFilterFlag() craveModel.Filter
	Apply(html *string) craveModel.Filter //return filter flag if it is filtered, return 0 if it is passed
}

type FilterChain struct {
	filters []IBusiness
}

type FilterStrategy struct {
	filterMap map[craveModel.Filter]IBusiness
}

func NewStrategy(filterMap map[craveModel.Filter]IBusiness) *FilterStrategy {
	return &FilterStrategy{filterMap: filterMap}
}

func (strat *FilterStrategy) GetNameFilterBusiness() IBusiness {
	return strat.filterMap[craveModel.NAME]
}
func (strat *FilterStrategy) GetFilterChain(filter craveModel.Filter) *FilterChain {
	var filterChain FilterChain
	for f, business := range strat.filterMap {
		if filter&f != 0 {
			filterChain.filters = append(filterChain.filters, business)
		}
	}
	return &filterChain
}

func (fc *FilterChain) Apply(html *string) craveModel.Filter {
	var filteredBy craveModel.Filter
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, filter := range fc.filters {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := filter.Apply(html)
			if result != 0 {
				mu.Lock()
				filteredBy |= result
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return filteredBy //return 0 when it is pass
}
