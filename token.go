package punkt

import (
  "strings"
  "regexp"
  "unicode"
  "unicode/utf8"
)

type Token struct {
  Value string
  Type string
  Flags map[string]bool
}

func MakeToken(token string, flags ...string) *Token {
  typeRegexp := regexp.MustCompile("^-?[\\.,]?\\d[\\d,\\.-]*\\.?$")
  t_type := typeRegexp.ReplaceAllString(strings.ToLower(token), "##number##")
  t_flags := make(map[string]bool)

  for _, f := range flags {
    t_flags[f] = true
  }
  
  return &Token{Value: token, Type: t_type, Flags: t_flags}
}

func (t Token) TypeWithoutPeriod() string {
  if len(t.Type) > 1 && strings.HasSuffix(t.Type, ".") {
    return strings.TrimRight(t.Type, ".")
  } else {
    return t.Type
  }
}

func (t Token) TypeWithoutSentencePeriod() string {
  v, ok := t.Flags["SentenceBreak"]

  if ok && v {
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

func (t Token) IsEllipsis() bool {
  matched, _ := regexp.MatchString("^\\.\\.+$", t.Value)
  return matched
}

func (t Token) IsNumber() bool {
  return strings.HasPrefix(t.Value, "##number##")
}

func (t Token) IsInitial() bool {
  matched, _ := regexp.MatchString("^[^\\W\\d]\\.$", t.Value)
  return matched
}

func (t Token) IsAlpha() bool {
  matched, _ := regexp.MatchString("^[^\\W\\d]+$", t.Value)
  return matched
}

func (t Token) IsNonPunctuation() bool {
  matched, _ := regexp.MatchString("[^\\W\\d]", t.Value)
  return matched
}
//  
//    def first_case
//      return :lower if first_lower?
//      return :upper if first_upper?
//      return :none
//    end
//  
//    def ends_with_period?
//      @period_final
//    end
//  
//    def is_ellipsis?
//      !(@token =~ /^\.\.+$/).nil?
//    end
//  
//    def is_number?
//      @type.start_with?("##number##")
//    end
//  
//    def is_initial?
//      !(@token =~ /^[^\W\d]\.$/).nil?
//    end
//  
//    def is_alpha?
//      !(@token =~ /^[^\W\d]+$/).nil?
//    end
//  
//    def is_non_punctuation?
//      !(@type =~ /[^\W\d]/).nil?
//    end
//  
//    def to_s
//      result = @token
//      result += '<A>' if @abbr
//      result += '<E>' if @ellipsis
//      result += '<S>' if @sentence_break
//      result
//    end

//    def inspect
//      "<#{to_s}>"
//    end
//  end
//end
