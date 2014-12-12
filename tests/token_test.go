package punkt

import (
  "fmt"
  "testing"
  "github.com/harrisj/punkt"
  . "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type TokenSuite struct{}

var tokenSuite = Suite(&TokenSuite{})

func (s *TokenSuite) TestTokenInit(c *C) {
  token := punkt.MakeToken("Test")

  c.Check(token.IsAbbr(), Equals, false)
}

func (s *TokenSuite) TestFlags(c *C) {
  token := punkt.MakeToken("Test")

  token.SetAbbr(true)
  c.Check(token.IsAbbr(), Equals, true)
  token.SetAbbr(false)
  c.Check(token.IsAbbr(), Equals, false)
}

func (s *TokenSuite) TestTypeAttributes(c *C) {
  token := punkt.MakeToken("Test")
  c.Check(token.Type, Equals, "test")

  token = punkt.MakeToken("Test.")
  c.Check(token.Type, Equals, "test.")

  token = punkt.MakeToken("Índico")
  c.Check(token.Type, Equals, "índico")   // does this need to be plain ASCII?
}

func (s *TokenSuite) TestTypeWithoutPeriod(c *C) {
  token := punkt.MakeToken("Test")
  c.Check(token.TypeWithoutPeriod(), Equals, "test")

  token = punkt.MakeToken("Test.")
  c.Check(token.TypeWithoutPeriod(), Equals, "test")

  token = punkt.MakeToken("123.")
  c.Check(token.TypeWithoutPeriod(), Equals, "##number##")
}

func (s *TokenSuite) TestEndWithPeriod(c *C) {
  token := punkt.MakeToken("Test")
  c.Check(token.EndsWithPeriod(), Equals, false)

  token = punkt.MakeToken("Test.")
  c.Check(token.EndsWithPeriod(), Equals, true) 
}

func (s *TokenSuite) TestTypeWithoutSentencePeriod(c *C) {
  token := punkt.MakeToken("Test")
  c.Check(token.TypeWithoutSentencePeriod(), Equals, "test")

  token = punkt.MakeToken("Test.")
  token.SetSentenceBreak(true)
  c.Check(token.TypeWithoutSentencePeriod(), Equals, "test")
}

func (s *TokenSuite) TestFirstUpper(c *C) {
  token := punkt.MakeToken("Test")
  c.Assert(token.FirstUpper(), Equals, true)

  token = punkt.MakeToken("Índico")
  c.Assert(token.FirstUpper(), Equals, true)

  token = punkt.MakeToken("test.")
  c.Assert(token.FirstUpper(), Equals, false)
}
  
func (s *TokenSuite) TestFirstLower(c *C) {
  token := punkt.MakeToken("Test")
  c.Assert(token.FirstLower(), Equals, false)

  token = punkt.MakeToken("índico")
  c.Assert(token.FirstLower(), Equals, true)

  token = punkt.MakeToken("test.")
  c.Assert(token.FirstLower(), Equals, true) 
}

func (s *TokenSuite) TestIsEllipsis(c *C) {
  token := punkt.MakeToken("...")
  c.Assert(token.MatchEllipsis(), Equals, true)

  token = punkt.MakeToken("..")
  c.Assert(token.MatchEllipsis(), Equals, true)

  token = punkt.MakeToken("..foo")
  c.Assert(token.MatchEllipsis(), Equals, false) 
}

func (s *TokenSuite) TestIsInitial(c *C) {
  token := punkt.MakeToken("C.")
  c.Assert(token.MatchInitial(), Equals, true)

  token = punkt.MakeToken("B.M.")
  c.Assert(token.MatchInitial(), Equals, false)
}

func (s *TokenSuite) TestIsAlpha(c *C) {
  token := punkt.MakeToken("foo")
  c.Assert(token.MatchAlpha(), Equals, true)

  token = punkt.MakeToken("!")
  c.Assert(token.MatchAlpha(), Equals, false)
}

func (s *TokenSuite) TestMatchNonPunctuation(c *C) {
  token := punkt.MakeToken("foo")
  c.Assert(token.MatchNonPunctuation(), Equals, true)

  token = punkt.MakeToken("!")
  c.Assert(token.MatchNonPunctuation(), Equals, false)
}

func (s *TokenSuite) TestString(c *C) {
  token := punkt.MakeToken("foo")
  token.SetAbbr(true)
  token.SetSentenceBreak(true)
  token.SetEllipsis(true)

  tokenStr := fmt.Sprintf("%v", token)
  c.Assert(tokenStr, Equals, "<<foo<A><E><S>>>")
}


