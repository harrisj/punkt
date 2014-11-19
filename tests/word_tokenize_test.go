package tests

import (
  "github.com/harrisj/punkt"
  . "gopkg.in/check.v1"
)

type WordTokenizeSuite struct{}

var wordTokenizeSuite = Suite(&WordTokenizeSuite{})

func (s *WordTokenizeSuite) TestWordTokenize(c *C) {
  c.Check(punkt.WordTokenize("apple pears"), DeepEquals, []string{"apple", "pears"})
  c.Check(punkt.WordTokenize("    apple pears    "), DeepEquals, []string{"apple", "pears"})
  c.Check(punkt.WordTokenize("[@apple]"), DeepEquals, []string{"[", "@", "apple", "]"})
  c.Check(punkt.WordTokenize("apple,pears"), DeepEquals, []string{"apple", ",", "pears"})
  c.Check(punkt.WordTokenize("apple! pears"), DeepEquals, []string{"apple", "!", "pears"})
  c.Check(punkt.WordTokenize("pears..."), DeepEquals, []string{"pears", "..."})
  c.Check(punkt.WordTokenize("self-conscious"), DeepEquals, []string{"self-conscious"})
  c.Check(punkt.WordTokenize("außer über"), DeepEquals, []string{"außer", "über"})
  c.Check(punkt.WordTokenize("apple (pears)"), DeepEquals, []string{"apple", "(", "pears", ")"})
  c.Check(punkt.WordTokenize("apple. pears."), DeepEquals, []string{"apple.", "pears."})
  c.Check(punkt.WordTokenize("apple.pears."), DeepEquals, []string{"apple.pears."})
  c.Check(punkt.WordTokenize("apple... pears"), DeepEquals, []string{"apple", "...", "pears"})
  c.Check(punkt.WordTokenize("apple -- pears"), DeepEquals, []string{"apple", "--", "pears"})

  sentence := "For example, the word \"abbreviation\" can itself be represented by the abbreviation abbr., abbrv. or abbrev."
  tokens := punkt.WordTokenize(sentence)
  c.Check(tokens, DeepEquals, []string{"For", "example", ",", "the", "word", "\"", "abbreviation", "\"", "can", "itself", "be", "represented", "by", "the", "abbreviation", "abbr.", ",", "abbrv.", "or", "abbrev."})
}
