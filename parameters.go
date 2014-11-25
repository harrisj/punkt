package punkt

import (
  "fmt"
)

type OrthoContext uint32

//#####################################################################
//{ Orthographic Context Constants
//#####################################################################
// The following constants are used to describe the orthographic
// contexts in which a word can occur.  BEG=beginning, MID=middle,
// UNK=unknown, UC=uppercase, LC=lowercase, NC=no case.
const (
  _                         = iota
  ORTHO_BEG_UC OrthoContext = 1 << (1 * iota) // beginning of a sentence with upper case.
  ORTHO_MID_UC                                // middle of a sentence with upper case.
  ORTHO_UNK_UC                                // unknown position in a sentence with upper case.
  ORTHO_BEG_LC                                // beginning of a sentence with lower case
  ORTHO_MID_LC                                // middle of a sentence with lower case
  ORTHO_UNK_LC                                // unknown position in a sentence with lower case
)

const (
  ORTHO_UC = ORTHO_BEG_UC + ORTHO_MID_UC + ORTHO_UNK_UC
  ORTHO_LC = ORTHO_BEG_LC + ORTHO_MID_LC + ORTHO_UNK_LC
)

type LanguageParameters struct {
  AbbrevTypes map[string]bool
  Collocations map[string]bool
  SentenceStarters map[string]bool
  OrthographicContext map[string]OrthoContext
}

func (p LanguageParameters) HasAbbrevType(s string) bool {
  return p.AbbrevTypes[s]
}

func (p *LanguageParameters) SaveAbbrevType(s string) {
  if len(p.AbbrevTypes) == 0 {
    p.ClearAbbrevTypes()
  }

  p.AbbrevTypes[s] = true
}

func (p *LanguageParameters) DeleteAbbrevType(s string) {
  if len(p.AbbrevTypes) == 0 {
    return
  }

  delete(p.AbbrevTypes, s)
}

func (p *LanguageParameters) ClearAbbrevTypes() {
  p.AbbrevTypes = make(map[string]bool)
}

func (p LanguageParameters) HasSentenceStarter(s string) bool {
  return p.SentenceStarters[s]
}

func (p *LanguageParameters) SaveSentenceStarter(s string) {
  if len(p.SentenceStarters) == 0 {
    p.ClearSentenceStarters()
  }

  p.SentenceStarters[s] = true
}

func (p *LanguageParameters) ClearSentenceStarters() {
  p.SentenceStarters = make(map[string]bool)
}

func collocationMapKey(s1, s2 string) (key string) {
  key = fmt.Sprintf("%v|%v", s1, s2)
  return
}

func (p LanguageParameters) HasCollocation(s1, s2 string) bool {
  return p.Collocations[collocationMapKey(s1,s2)]
}

func (p *LanguageParameters) SaveCollocation(s1, s2 string) {
  if len(p.Collocations) == 0 {
    p.ClearCollocations()
  }

  p.Collocations[collocationMapKey(s1,s2)] = true
}

func (p *LanguageParameters) ClearCollocations() {
  p.Collocations = make(map[string]bool)
}

func (p *LanguageParameters) ClearOrthographicContext() {
  p.OrthographicContext = make(map[string]OrthoContext)
}

func (p *LanguageParameters) GetOrthographicContext(s string) OrthoContext {
  return p.OrthographicContext[s]
}

func (p *LanguageParameters) AddOrthographicContext(s string, flag OrthoContext) {
  if len(p.OrthographicContext) == 0 {
    p.ClearOrthographicContext()
  }

  p.OrthographicContext[s] |= flag
}

func (p *LanguageParameters) DeleteOrthographicContext(s string, flag OrthoContext) {
  if len(p.OrthographicContext) == 0 {
    p.ClearOrthographicContext()
  }

  p.OrthographicContext[s] ^= flag
}

