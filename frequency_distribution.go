package punkt

import (
  "sort"
)

type SampleCount struct {
  Sample string
  Count int
}

// ByCount implements sort.Interface for []SampleCount based on
// the Count field.
type ByCount []SampleCount

func (a ByCount) Len() int           { return len(a) }
func (a ByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCount) Less(i, j int) bool { return a[i].Count > a[j].Count } // descending order

type FrequencyDistribution struct {
  N int
  Counts map[string]int
  MaxSample SampleCount
  Sorted []SampleCount
}

func (f *FrequencyDistribution) Clear() {
  f.N = 0
  // FIXME: Is this cool?
  f.Counts = make(map[string]int)
}

func (f *FrequencyDistribution) ClearCaches() {
  f.MaxSample = SampleCount{"", 0}
  f.Sorted = []SampleCount{}
}

func (f *FrequencyDistribution) Get(sample string) int {
  return f.Counts[sample]
}

func (f *FrequencyDistribution) Set(sample string, value int) {
  if len(f.Counts) == 0 {
    f.Counts = make(map[string]int)
  }

  f.N += (value - f.Get(sample))
  f.Counts[sample] = value
  f.ClearCaches()
}

func (f *FrequencyDistribution) Inc(sample string) {
  f.Set(sample, f.Get(sample)+1)
}

func (f *FrequencyDistribution) IncBy(sample string, n int) {
  f.Set(sample, f.Get(sample)+n)
}

func (f *FrequencyDistribution) FrequencyOf(sample string) float64 {
  if f.N == 0 {
    return 0
  }

  return float64(f.Get(sample))/float64(f.N)
}

func (f *FrequencyDistribution) Max() SampleCount {
  if f.MaxSample.Sample == "" {
    maxSample := ""
    maxCount := -1

    for k, v := range f.Counts {
      if v > maxCount {
        maxSample = k
        maxCount = v
      }
    }

    f.MaxSample = SampleCount{maxSample,maxCount}
  }

  return f.MaxSample
}

func (f *FrequencyDistribution) OrderedSamples() []SampleCount {
  if len(f.Sorted) == 0 {
    f.Sorted = make([]SampleCount, 0, len(f.Counts))

    for k, v := range f.Counts {
      f.Sorted = append(f.Sorted, SampleCount{k,v})
    }

    sort.Sort(ByCount(f.Sorted))
  }

  return f.Sorted
}
