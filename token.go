package punkt

import (
  "strings"
  "regexp"
  "unicode"
  "unicode/utf8"
)

type TokenFlags byte

const (
  TOK_PARAGRAPH_START TokenFlags = 1<<(1*iota)
  TOK_LINE_START
  TOK_SENTENCE_BREAK
  TOK_ABBR
  TOK_ELLIPSIS
)

type Token struct {
  Value string
  Type string
  Flags TokenFlags
}

// optional argument for flags
func MakeToken(token string) *Token {
  typeRegexp := regexp.MustCompile("^-?[\\.,]?\\d[\\d,\\.-]*\\.?$")
  t_type := typeRegexp.ReplaceAllString(strings.ToLower(token), "##number##")
 
  return &Token{Value: token, Type: t_type}
}

func (t Token) IsAbbr() bool {
  return t.Flags & TOK_ABBR != 0
}

func (t *Token) SetAbbr(b bool) {
  if b {
    t.Flags |= TOK_ABBR
  } else {
    t.Flags ^= TOK_ABBR
  }
}

func (t Token) IsSentenceBreak() bool {
  return t.Flags & TOK_SENTENCE_BREAK != 0
}

func (t *Token) SetSentenceBreak(b bool) {
  if b {
    t.Flags |= TOK_SENTENCE_BREAK
  } else {
    t.Flags ^= TOK_SENTENCE_BREAK
  }
}

func (t Token) IsEllipsis() bool {
  return t.Flags & TOK_ELLIPSIS != 0
}

func (t *Token) SetEllipsis(b bool) {
  if b {
    t.Flags |= TOK_ELLIPSIS
  } else {
    t.Flags ^= TOK_ELLIPSIS
  }
}

func (t Token) IsParagraphStart() bool {
  return t.Flags & TOK_PARAGRAPH_START != 0
}

func (t *Token) SetParagraphStart(b bool) {
  if b {
    t.Flags |= TOK_PARAGRAPH_START
  } else {
    t.Flags ^= TOK_PARAGRAPH_START
  }
}

func (t Token) IsLineStart() bool {
  return t.Flags & TOK_LINE_START != 0
}

func (t *Token) SetLineStart(b bool) {
  if b {
    t.Flags |= TOK_LINE_START
  } else {
    t.Flags ^= TOK_LINE_START
  }
}

func (t Token) TypeWithoutPeriod() string {
  if len(t.Type) > 1 && strings.HasSuffix(t.Type, ".") {
    return strings.TrimRight(t.Type, ".")
  } else {
    return t.Type
  }
}

func (t Token) TypeWithoutSentencePeriod() string {
  if t.IsSentenceBreak() {
    return t.TypeWithoutPeriod()
  } else {
    return t.Type
  }
}

func (t Token) FirstUpper() bool {
  if len(t.Value) == 0 {
    return false
  } else {
    // approach suggested by 
    runeValue, _ := utf8.DecodeRuneInString(t.Value)
    return unicode.IsUpper(runeValue)
  }
}

func (t Token) FirstLower() bool {
  if len(t.Value) == 0 {
    return false
  } else {
    runeValue, _ := utf8.DecodeRuneInString(t.Value)
    return unicode.IsLower(runeValue)
  }
}

func (t Token) FirstCase() string {
  if t.FirstUpper() {
    return "upper"
  } else if t.FirstLower() {
    return "lower"
  } else {
    return "none"
  }
}

func (t Token) EndsWithPeriod() bool {
  return strings.HasSuffix(t.Value, ".")
}

func (t Token) MatchEllipsis() bool {
  matched, _ := regexp.MatchString("^\\.\\.+$", t.Value)
  return matched
}

func (t Token) MatchNumber() bool {
  return strings.HasPrefix(t.Value, "##number##")
}

func (t Token) MatchInitial() bool {
  matched, _ := regexp.MatchString("^[^\\W\\d]\\.$", t.Value)
  return matched
}

func (t Token) MatchAlpha() bool {
  matched, _ := regexp.MatchString("^[^\\W\\d]+$", t.Value)
  return matched
}

func (t Token) MatchNonPunctuation() bool {
  matched, _ := regexp.MatchString("[^\\W\\d]", t.Value)
  return matched
}

func (t Token) String() (out string) {
  out = "<<"
  out += t.Value

  if t.IsAbbr() {
    out += "<A>"
  }

  if t.IsEllipsis() {
    out += "<E>"
  }

  if t.IsSentenceBreak() {
    out += "<S>"
  }

  out += ">>"
  
  return
}

// Is anybody using this one
//    def first_case
//      return :lower if first_lower?
//      return :upper if first_upper?
//      return :none
//    end
//  
//  
//    def to_s
//      result = @token
//      result += '<A>' if @abbr
//      result += '<E>' if @ellipsis
//      result += '<S>' if @sentence_break
//      result
//    end
//end
