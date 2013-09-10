package decisiontrees

import (
	"math/rand"
	"sort"
)

type Example struct {
	Label         float64
	WeightedLabel float64
	Features      map[int64]float64
}

func (e Example) asBool() bool {
	if e.Label > 0 {
		return true
	}
	return false
}

type Examples []*Example

func (e Examples) subsampleExamples(samplingRate float64) Examples {
	for i := range e {
		j := rand.Intn(i + 1)
		e[i], e[j] = e[j], e[i]
	}

	return e[:int64(float64(len(e))*samplingRate)]
}

func (e Examples) crossValidationSamples(folds int) []Examples {
	crossValidatedSamples := make([]Examples, folds)
	for i, _ := range crossValidatedSamples {
		crossValidatedSamples[i] = make([]*Example, 0, len(e)/folds)
	}

	// Do a Fischer-Yates shuffle of the input array
	for i := range e {
		j := rand.Intn(i + 1)
		e[i], e[j] = e[j], e[i]
	}

	for i, ex := range e {
		fold := i % len(crossValidatedSamples)
		crossValidatedSamples[fold] = append(crossValidatedSamples[fold], ex)
	}
	return crossValidatedSamples
}

func (e Examples) boostrapFeatures(size int) []int64 {
	subsample := make([]int64, size)
	allFeatures := e.getFeatures()
	for i, _ := range subsample {
		subsample[i] = allFeatures[i]
	}

	for i := size + 1; i < len(allFeatures); i++ {
		j := int(rand.Int31n(int32(i)))
		if j < size {
			subsample[j] = allFeatures[i]
		}
	}
	return subsample
}

type By func(e1, e2 *Example) bool

func (by By) Sort(examples Examples) {
	es := &exampleSorter{
		examples: examples,
		by:       by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(es)
}

type exampleSorter struct {
	examples Examples
	by       By
}

func (s *exampleSorter) Len() int {
	return len(s.examples)
}

func (e *exampleSorter) Swap(i int, j int) {
	e.examples[i], e.examples[j] = e.examples[j], e.examples[i]
}

func (e *exampleSorter) Less(i int, j int) bool {
	return e.by(e.examples[i], e.examples[j])
}

func (e Examples) getFeatures() []int64 {
	vals := make(map[int64]bool)
	for _, example := range e {
		for k, _ := range example.Features {
			vals[k] = true
		}
	}
	res := make([]int64, 0, len(vals))
	for k, _ := range vals {
		res = append(res, k)
	}
	return res
}
