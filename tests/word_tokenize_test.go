package tests

import (
  "github.com/harrisj/punkt"
  . "gopkg.in/check.v1"
)

type WordTokenizeSuite struct{}

var wordTokenizeSuite = Suite(&WordTokenizeSuite{})

func (s *WordTokenizeSuite) TestWordTokenize(c *C) {
  c.Check(punkt.SplitTextIntoWords("apple pears"), DeepEquals, []string{"apple", "pears"})
  c.Check(punkt.SplitTextIntoWords("    apple pears    "), DeepEquals, []string{"apple", "pears"})
  c.Check(punkt.SplitTextIntoWords("[@apple]"), DeepEquals, []string{"[", "@", "apple", "]"})
  c.Check(punkt.SplitTextIntoWords("apple,pears"), DeepEquals, []string{"apple", ",", "pears"})
  c.Check(punkt.SplitTextIntoWords("apple! pears"), DeepEquals, []string{"apple", "!", "pears"})
  c.Check(punkt.SplitTextIntoWords("pears..."), DeepEquals, []string{"pears", "..."})
  c.Check(punkt.SplitTextIntoWords("self-conscious"), DeepEquals, []string{"self-conscious"})
  c.Check(punkt.SplitTextIntoWords("außer über"), DeepEquals, []string{"außer", "über"})
  c.Check(punkt.SplitTextIntoWords("apple (pears)"), DeepEquals, []string{"apple", "(", "pears", ")"})
  c.Check(punkt.SplitTextIntoWords("apple. pears."), DeepEquals, []string{"apple.", "pears."})
  c.Check(punkt.SplitTextIntoWords("apple.pears."), DeepEquals, []string{"apple.pears."})
  c.Check(punkt.SplitTextIntoWords("apple... pears"), DeepEquals, []string{"apple", "...", "pears"})
  c.Check(punkt.SplitTextIntoWords("apple -- pears"), DeepEquals, []string{"apple", "--", "pears"})

  sentence := "For example, the word \"abbreviation\" can itself be represented by the abbreviation abbr., abbrv. or abbrev."
  tokens := punkt.SplitTextIntoWords(sentence)
  c.Check(tokens, DeepEquals, []string{"For", "example", ",", "the", "word", "\"", "abbreviation", "\"", "can", "itself", "be", "represented", "by", "the", "abbreviation", "abbr.", ",", "abbrv.", "or", "abbrev."})
}
