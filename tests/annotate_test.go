package tests

import (
  . "github.com/harrisj/punkt"
  . "gopkg.in/check.v1"
  "fmt"
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
  c.Check(GuessOrthographicBoundary(s.parameters, token), Equals, ORTHO_BOUND_FALSE)

  token = MakeToken("This")
  c.Check(GuessOrthographicBoundary(s.parameters, token), Equals, ORTHO_BOUND_UNK)

  // add to parameters
  s.parameters.AddOrthographicContext(token.TypeWithoutSentencePeriod(), ORTHO_LC)
  c.Check(GuessOrthographicBoundary(s.parameters, token), Equals, ORTHO_BOUND_TRUE)

  token = MakeToken("this")
  c.Check(GuessOrthographicBoundary(s.parameters, token), Equals, ORTHO_BOUND_UNK)
  s.parameters.AddOrthographicContext(token.TypeWithoutSentencePeriod(), ORTHO_UC)
  c.Check(GuessOrthographicBoundary(s.parameters, token), Equals, ORTHO_BOUND_FALSE)
}

func (s *AnnotateSuite) TestAnnotateFirstPass(c *C) {
  tokens := make([]*Token, 4)
  tokens[0] = MakeToken("e.g.")
  tokens[1] = MakeToken(",")
  tokens[2] = MakeToken("Apple")
  tokens[3] = MakeToken("Computer.")

  s.parameters.SaveAbbrevType("e.g")

  c.Assert(s.parameters.HasAbbrevType("e.g"), Equals, true)

  tokens = AnnotateFirstPass(s.parameters, tokens)

  c.Check(tokens[0].IsAbbr(), Equals, true)
  c.Check(tokens[0].IsSentenceBreak(), Equals, false)
  c.Check(tokens[1].IsAbbr(), Equals, false)
  c.Check(tokens[3].IsSentenceBreak(), Equals, true)
}

func (s *AnnotateSuite) TestAnnotateSecondPass(c *C) {
  // words := TokenizeTextToWords("At 9 P.M., I went to bed.")

  // tokens := make([]Token, len(words))
  // for i, w := range words {
  //   tokens[i] = *MakeToken(w)
  // }
  // 
  
  // parameters := t.TrainWithText(str)
}

func (s *AnnotateSuite) TestTokenizer(c *C) {
  str := "When Mr. Gregor Samsa woke up one morning from unsettling dreams, he found himself changed in his bed into a monstrous vermin. He was lying on his back as hard as armor plate, and when he lifted his head a little, he saw his vaulted brown belly, sectioned by arch-shaped ribs, to whose dome the cover, about to slide off completely, could barely cling. His many legs, pitifully thin compared with the size of the rest of him, were waving helplessly before his eyes."

  t := new(Tokenizer)
  t.SetLanguage("english")

  fmt.Printf("%#v\n", t.SentencesFromText(str))
}
