package tests

import (
  . "github.com/harrisj/punkt"
  . "gopkg.in/check.v1"
)

type AnnotateSuite struct{
  parameters *LanguageParameters
}

var annotateSuite = Suite(&AnnotateSuite{})

func (s *AnnotateSuite) SetUpTest(c *C) {
  s.parameters = new(LanguageParameters)
}

func (s *AnnotateSuite) TestGuessOrthographicBoundary(c *C) {
  token := MakeToken("A;B")
  c.Check(GuessOrthographicBoundary(s.parameters, *token), Equals, ORTHO_BOUND_FALSE)

  token = MakeToken("This")
  c.Check(GuessOrthographicBoundary(s.parameters, *token), Equals, ORTHO_BOUND_UNK)

  // add to parameters
  s.parameters.AddOrthographicContext(token.TypeWithoutSentencePeriod(), ORTHO_LC)
  c.Check(GuessOrthographicBoundary(s.parameters, *token), Equals, ORTHO_BOUND_TRUE)

  token = MakeToken("this")
  c.Check(GuessOrthographicBoundary(s.parameters, *token), Equals, ORTHO_BOUND_UNK)
  s.parameters.AddOrthographicContext(token.TypeWithoutSentencePeriod(), ORTHO_UC)
  c.Check(GuessOrthographicBoundary(s.parameters, *token), Equals, ORTHO_BOUND_FALSE)
}

func (s *AnnotateSuite) TestAnnotateFirstPass(c *C) {
  tokens := make([]Token, 4)
  tokens[0] = *MakeToken("e.g.")
  tokens[1] = *MakeToken(",")
  tokens[2] = *MakeToken("Apple")
  tokens[3] = *MakeToken("Computer.")

  s.parameters.SaveAbbrevType("e.g")

  c.Assert(s.parameters.HasAbbrevType("e.g"), Equals, true)

  tokens = AnnotateFirstPass(s.parameters, tokens)

  c.Check(tokens[0].IsAbbr(), Equals, true)
  c.Check(tokens[0].IsSentenceBreak(), Equals, false)
  c.Check(tokens[1].IsAbbr(), Equals, false)
  c.Check(tokens[3].IsSentenceBreak(), Equals, true)
}

// func GuessOrthographicBoundary(parameters *LanguageParameters, token Token) OrthoHeuristicResult {
//   punctRegexp := regexp.MustCompile("[;,:.!?]")

//   if punctRegexp.MatchString(token.Value) {
//     return ORTHO_BOUND_FALSE
//   }

//   ortho_context := parameters.GetOrthographicContext(token.TypeWithoutSentencePeriod())

//   if token.FirstUpper() && (ortho_context & ORTHO_LC != 0) && !(ortho_context & ORTHO_MID_UC != 0) {
//     return ORTHO_BOUND_TRUE
//   } else if token.FirstLower() && ((ortho_context & ORTHO_UC != 0) || !(ortho_context & ORTHO_BEG_LC != 0)) {
//     return ORTHO_BOUND_FALSE
//   } else {
//     return ORTHO_BOUND_UNK
//   }
// }