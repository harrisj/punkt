package tests

import (
  "github.com/harrisj/punkt"
  . "gopkg.in/check.v1"
)

type WordTokenizeSuite struct{}

var wordTokenizeSuite = Suite(&WordTokenizeSuite{})

func (s *WordTokenizeSuite) TestWordTokenize(c *C) {
  c.Check(punkt.TokenizeTextToWords("apple pears"), DeepEquals, []string{"apple", "pears"})
  c.Check(punkt.TokenizeTextToWords("    apple pears    "), DeepEquals, []string{"apple", "pears"})
  c.Check(punkt.TokenizeTextToWords("[@apple]"), DeepEquals, []string{"[", "@", "apple", "]"})
  c.Check(punkt.TokenizeTextToWords("apple,pears"), DeepEquals, []string{"apple", ",", "pears"})
  c.Check(punkt.TokenizeTextToWords("apple! pears"), DeepEquals, []string{"apple", "!", "pears"})
  c.Check(punkt.TokenizeTextToWords("pears..."), DeepEquals, []string{"pears", "..."})
  c.Check(punkt.TokenizeTextToWords("self-conscious"), DeepEquals, []string{"self-conscious"})
  c.Check(punkt.TokenizeTextToWords("außer über"), DeepEquals, []string{"außer", "über"})
  c.Check(punkt.TokenizeTextToWords("apple (pears)"), DeepEquals, []string{"apple", "(", "pears", ")"})
  c.Check(punkt.TokenizeTextToWords("apple. pears."), DeepEquals, []string{"apple.", "pears."})
  c.Check(punkt.TokenizeTextToWords("apple.pears."), DeepEquals, []string{"apple.pears."})
  c.Check(punkt.TokenizeTextToWords("apple... pears"), DeepEquals, []string{"apple", "...", "pears"})
  c.Check(punkt.TokenizeTextToWords("apple -- pears"), DeepEquals, []string{"apple", "--", "pears"})

  sentence := "For example, the word \"abbreviation\" can itself be represented by the abbreviation abbr., abbrv. or abbrev."
  tokens := punkt.TokenizeTextToWords(sentence)
  c.Check(tokens, DeepEquals, []string{"For", "example", ",", "the", "word", "\"", "abbreviation", "\"", "can", "itself", "be", "represented", "by", "the", "abbreviation", "abbr.", ",", "abbrv.", "or", "abbrev."})
}
