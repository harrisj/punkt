package punkt

import (
  "strings"
  "regexp"
)

// A map from context position and first-letter case to the
// appropriate orthographic context flag.
var ORTHO_MAP = map[string]OrthoContext {
    "initial|upper": ORTHO_BEG_UC,
    "internal|upper": ORTHO_MID_UC,
    "unknown|upper": ORTHO_UNK_UC,
    "initial|lower": ORTHO_BEG_LC,
    "internal|lower": ORTHO_MID_LC,
    "unknown|lower": ORTHO_UNK_LC,
}

type OrthoHeuristicResult byte

const (
  ORTHO_BOUND_TRUE OrthoHeuristicResult = iota
  ORTHO_BOUND_FALSE
  ORTHO_BOUND_UNK
)

// the orthographic heuristic, which decides for a token following an abbreviation or an ellipsis on the basis 
// of the orthographic statistics gathered for all word types whether it represents good evidence for a preceding
// sentence boundary or not.
func GuessOrthographicBoundary(parameters *LanguageParameters, token Token) OrthoHeuristicResult {
  punctRegexp := regexp.MustCompile("[;,:.!?]")

  if punctRegexp.MatchString(token.Value) {
    return ORTHO_BOUND_FALSE
  }

  ortho_context := parameters.GetOrthographicContext(token.TypeWithoutSentencePeriod())

  if token.FirstUpper() && (ortho_context & ORTHO_LC != 0) && !(ortho_context & ORTHO_MID_UC != 0) {
    return ORTHO_BOUND_TRUE
  } else if token.FirstLower() && ((ortho_context & ORTHO_UC != 0) || !(ortho_context & ORTHO_BEG_LC != 0)) {
    return ORTHO_BOUND_FALSE
  } else {
    return ORTHO_BOUND_UNK
  }
}

func AnnotateFirstPass(parameters *LanguageParameters, tokens []Token) []Token {
  for i := range tokens {
    str := tokens[i].Value

    switch {
    case str == "." || str == "?" || str == "!":
      tokens[i].SetSentenceBreak(true)
    case tokens[i].MatchEllipsis():
      tokens[i].SetEllipsis(true)
    case tokens[i].EndsWithPeriod():
      tokLow := strings.ToLower(str[0:len(str)-1])

      tokSplit := strings.Split(tokLow, "-")

      if parameters.HasAbbrevType(tokLow) || (len(tokSplit) > 1 && parameters.HasAbbrevType(tokSplit[len(tokSplit)-1])) {
        tokens[i].SetAbbr(true)
      } else {
        tokens[i].SetSentenceBreak(true)
      }
    }
  }

  //fmt.Println("First Pass: ", tokens)
  return tokens
}

func AnnotateSecondPass(parameters *LanguageParameters, tokens []Token) []Token {
  for i := range tokens {
    if i == 0 {
      continue
    }

    // make sure to use pointers
    tok1 := &tokens[i-1]
    tok2 := &tokens[i]

    if tok1.EndsWithPeriod() {
      continue
    }

    // FIXME: Check these are really not used
    //t1Value := tok1.Value
    t1Type := tok1.Type
    //t2Value := tok2.Value
    t2Type := tok2.TypeWithoutSentencePeriod()
    t1Initial := tok1.MatchInitial()

    if parameters.HasCollocation(t1Type, t2Type) {
      tok1.SetSentenceBreak(false)
      tok1.SetAbbr(true)
      continue
    }

    if (tok1.IsAbbr() || tok1.IsEllipsis()) && !(t1Initial) {
      if GuessOrthographicBoundary(parameters, *tok2) == ORTHO_BOUND_TRUE {
        tok1.SetSentenceBreak(true)
        continue
      }

      if tok1.FirstUpper() && parameters.HasSentenceStarter(t2Type) {
        tok1.SetSentenceBreak(true)
        continue
      }
    }

    if t1Initial || t1Type == "##number##" {
      ot := GuessOrthographicBoundary(parameters, *tok2)
      if ot == ORTHO_BOUND_FALSE {
        tok1.SetSentenceBreak(false)
        tok1.SetAbbr(true)
        continue
      }

      if ot == ORTHO_BOUND_UNK && t1Initial && tok2.FirstUpper() {
        ot2 := parameters.GetOrthographicContext(t2Type)
        if !(ot2 & ORTHO_LC != 0) {
          tok1.SetSentenceBreak(false)
          tok1.SetAbbr(true)
        }
      }
    }
  }

  return tokens
}

