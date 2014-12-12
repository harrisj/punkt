package punkt

import (
  "math"
  "github.com/harrisj/punkt"
  . "gopkg.in/check.v1"
)

type FrequencyDistributionSuite struct{
  Words []string
}

var frequencyDistributionSuite = Suite(&FrequencyDistributionSuite{})

func (s *FrequencyDistributionSuite) SetUpSuite(c *C) {
  s.Words = []string{"two", "one", "three", "one", "one", "three", "two", "one", "two"}
}

func (s *FrequencyDistributionSuite) TestIncrementCountOnGivenSample(c *C) {
  Fd := new(punkt.FrequencyDistribution)

  for _, word := range s.Words {
    Fd.Inc(word)
  }

  c.Check(Fd.Get("one"), Equals, 4)
  c.Check(Fd.Get("two"), Equals, 3)
  c.Check(Fd.Get("three"), Equals, 2)
  c.Check(Fd.N, Equals, 9)
}

func (s *FrequencyDistributionSuite) TestIncrementByCountOnGivenSample(c *C) {
  Fd := new(punkt.FrequencyDistribution)

  for _, word := range s.Words {
    Fd.IncBy(word, 2)
  }

  c.Check(Fd.Get("one"), Equals, 8)
  c.Check(Fd.Get("two"), Equals, 6)
  c.Check(Fd.Get("three"), Equals, 4)
  c.Check(Fd.N, Equals, 18)
}

func (s *FrequencyDistributionSuite) TestDirectCountAttribution(c *C) {
  Fd := new(punkt.FrequencyDistribution)

  Fd.Set("one", 10)
  Fd.Set("two", 20)
  Fd.Set("three", 30)

  c.Check(Fd.Get("one"), Equals, 10)
  c.Check(Fd.Get("two"), Equals, 20)
  c.Check(Fd.Get("three"), Equals, 30)
  c.Check(Fd.N, Equals, 60)
}

func (s *FrequencyDistributionSuite) TestGetSampleFrequencies(c *C) {
  Fd := new(punkt.FrequencyDistribution)

  for _, word := range s.Words {
    Fd.Inc(word)
  }

  total := Fd.FrequencyOf("one") + Fd.FrequencyOf("two") + Fd.FrequencyOf("three")
  c.Check(math.Ceil(total), Equals, 1.0)
}

func (s *FrequencyDistributionSuite) TestSampleWithMaxOccurrences(c *C) {
  Fd := new(punkt.FrequencyDistribution)

  for _, word := range s.Words {
    Fd.Inc(word)
  }

  c.Check(Fd.Max(), DeepEquals, punkt.SampleCount{"one",4})
}

func (s *FrequencyDistributionSuite) TestOrderedKeyRetrieval(c *C) {
  Fd := new(punkt.FrequencyDistribution)

  for _, word := range s.Words {
    Fd.Inc(word)
  }

  c.Check(Fd.OrderedSamples(), DeepEquals, []punkt.SampleCount{{"one",4},{"two",3},{"three",2}})
}

// FIXME: Do we need delete?

func (s *FrequencyDistributionSuite) TestEmptyDistribution(c *C) {
  Fd := new(punkt.FrequencyDistribution)

  c.Check(Fd.Get("a sample"), Equals, 0)
  c.Check(Fd.N, Equals, 0)
  c.Check(Fd.FrequencyOf("a sample"), Equals, float64(0))
}
