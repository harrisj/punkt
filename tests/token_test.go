package tests

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
  opts := punkt.TokenOptions {
    "ParagraphStart": true, 
    "LineStart": true, 
    "SentenceBreak": true, 
    "Abbr": true,
  }

  token := punkt.MakeToken("Test", opts)
  c.Check(token.Flags["ParagraphStart"], Equals, true)
  c.Check(token.Flags["LineStart"], Equals, true)
  c.Check(token.Flags["SentenceBreak"], Equals, true)
  c.Check(token.Flags["Abbr"], Equals, true)

  v := token.Flags["Ellipsis"]
  c.Check(v, Equals, false)
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

func (s *TokenSuite) TestTypeWithoutSentencePeriod(c *C) {
  token := punkt.MakeToken("Test")
  c.Check(token.TypeWithoutSentencePeriod(), Equals, "test")

  opts := make(punkt.TokenOptions)
  opts["SentenceBreak"] = true
  token = punkt.MakeToken("Test.", opts)
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
  c.Assert(token.IsEllipsis(), Equals, true)

  token = punkt.MakeToken("..")
  c.Assert(token.IsEllipsis(), Equals, true)

  token = punkt.MakeToken("..foo")
  c.Assert(token.IsEllipsis(), Equals, false) 
}

func (s *TokenSuite) TestIsInitial(c *C) {
  token := punkt.MakeToken("C.")
  c.Assert(token.IsInitial(), Equals, true)

  token = punkt.MakeToken("B.M.")
  c.Assert(token.IsInitial(), Equals, false)
}

func (s *TokenSuite) TestIsAlpha(c *C) {
  token := punkt.MakeToken("foo")
  c.Assert(token.IsAlpha(), Equals, true)

  token = punkt.MakeToken("!")
  c.Assert(token.IsAlpha(), Equals, false)
}

func (s *TokenSuite) TestIsNonPunctuation(c *C) {
  token := punkt.MakeToken("foo")
  c.Assert(token.IsNonPunctuation(), Equals, true)

  token = punkt.MakeToken("!")
  c.Assert(token.IsNonPunctuation(), Equals, false)
}

func (s *TokenSuite) TestString(c *C) {
  opts := punkt.TokenOptions {
    "Abbr": true,
    "SentenceBreak": true,
    "Ellipsis": true,
  }

  token := punkt.MakeToken("foo", opts)

  tokenStr := fmt.Sprintf("%v", token)
  c.Assert(tokenStr, Equals, "foo<A><E><S>")
}

//   def test_first_case    
//     token = Punkt::Token.new("Test")
//     assert_equal :upper, token.first_case

//     token = Punkt::Token.new("Índico")
//     assert_equal :upper, token.first_case

//     token = Punkt::Token.new("test.")
//     assert_equal :lower, token.first_case
    
//     token = Punkt::Token.new("@")
//     assert_equal :none, token.first_case
//   end
  
  
//   def test_to_s_and_inspect
//     token = Punkt::Token.new("foo", :abbr => true, :sentence_break => true, :ellipsis => true)
    
//     assert_equal "<foo<A><E><S>>", token.inspect
//   end
  
// end

