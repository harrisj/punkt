package punkt

import (
  "regexp"
)

type Tokenizer struct {
  parameters *LanguageParameters
}

func (t *Tokenizer) SetParameters(l *LanguageParameters) {
  t.parameters = l
}

// A shortcut to set the parameters for a specific language
func (t *Tokenizer) SetLanguage(lang string) {
  t.SetParameters(LoadLanguage(lang))
}

func (t Tokenizer) SentencesFromText(text string) (sentences []string) {
  sentences = t.splitIntoSentences(text)
  //sentences = t.realignBoundaries(text, sentences) //if options[:realign_boundaries]
  return
}

func (t Tokenizer) splitIntoSentences(input string) []string {
  out := make([]string, 0)
  currentSentenceStart := 0

  scanRegexp := regexp.MustCompile("\\S*[.?!]([\\?!\\)\";}\\]\\*:@'\\({\\[,]|\\s+\\S+)")
  periodRegexp := regexp.MustCompile("[.?!]")
  nextTokenRegexp := regexp.MustCompile("^\\s+(\\S+)")

  matches := scanRegexp.FindAllStringIndex(input, -1)

  if matches == nil {
    return out
  }

  for _, a := range matches {
    //fmt.Println("Scan:", input[a[0]:a[1]])

    context := input[a[0]:a[1]]
    if t.textContainsSentenceBreak(context) {
      // now I need to find where the period is in my big match
      pm := periodRegexp.FindStringIndex(context)

      // actual string end is a[0] + pm[1]
      sentence := input[currentSentenceStart:(a[0]+pm[1])]
      //fmt.Println("APPEND", sentence)
      out = append(out, sentence)

      // ugly, but let's look for next token
      nm := nextTokenRegexp.FindStringSubmatchIndex(input[(a[0]+pm[1]):len(input)])
      if nm == nil {
        currentSentenceStart = a[1]
      } else {
        currentSentenceStart = a[0]+pm[1]+nm[2]  // nm[2] = beginning of submatch, rel to a[0]+pm[1]
      }
    }
  }

  out = append(out, input[currentSentenceStart:len(input)])
  return out
}

func (t Tokenizer) textContainsSentenceBreak(text string) bool {
  found := false

  tokens := AnnotateTokens(t.parameters, TokenizeText(text))

  // don't return true if last token is a sentence break
  for _, tok := range tokens {
    if found {
      return true
    }

    if tok.IsSentenceBreak() {
      found = true
    }
  }

  return false
}

  
//     def realign_boundaries(text, sentences)
//       result = []
//       realign = 0
//       pair_each(sentences) do |i1, i2|
//         s1 = text[i1[0]..i1[1]]
//         s2 = i2 ? text[i2[0]..i2[1]] : nil
//         #s1 = s1[realign..(s1.size-1)]
//         unless s2
//           result << [i1[0]+realign, i1[1]] if s1
//           next
//         end
//         if match = @language_vars.re_boundary_realignment.match(s2)
//           result << [i1[0]+realign, i1[1]+match[0].strip.size] #s1 + match[0].strip()
//           realign = match.end(0)
//         else
//           result << [i1[0]+realign, i1[1]] if s1
//           realign = 0
//         end
//       end
//       return result
//     end
    
//   end
// end
