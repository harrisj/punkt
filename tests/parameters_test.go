package tests

import (
  . "github.com/harrisj/punkt"
  . "gopkg.in/check.v1"
)

type LanguageParametersSuite struct{
}

var languageParametersSuite = Suite(&LanguageParametersSuite{})

func (s *LanguageParametersSuite) TestOrthographicContext(c *C) {
  p := new(LanguageParameters)

  context := p.GetOrthographicContext("Dog")
  c.Check(context, Equals, OrthoContext(0))
  p.AddOrthographicContext("Dog", ORTHO_LC)
  c.Check(p.GetOrthographicContext("Dog"), Equals, ORTHO_LC)
  p.DeleteOrthographicContext("Dog", ORTHO_LC)
  c.Check(p.GetOrthographicContext("Dog"), Equals, OrthoContext(0))
}