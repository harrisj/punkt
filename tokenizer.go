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


    
//     def sentences_from_tokens(tokens)
//       tokens = annotate_tokens(tokens.map { |t| @token_class.new(t) })
      
//       sentences = []
//       sentence = []
//       tokens.each do |t|
//         sentence << t.token
//         if t.sentence_break
//           sentences << sentence
//           sentence = [] 
//         end
//       end
//       sentences << sentence unless sentence.empty?
      
//       return sentences
//     end
    
//     class << self
//       def sentences_text(text, sentences_indexes)
//         sentences_indexes.map { |index| text[index[0]..index[1]] }
//       end
      
//       def tokenized_sentences(text, sentences_indexes)
//         tokenizer = Punkt::Base.new()
//         self.sentences_text(text, sentences_indexes).map { |text| tokenizer.tokenize_words(text, :output => :string) }
//       end
//     end
    
//   private
  
//     def train(train_text)
//       @trainer = Punkt::Trainer.new(@language_vars, @token_class) unless @trainer
//       @trainer.train(train_text)
//       @parameters = @trainer.parameters
//     end
  
//     def split_in_sentences(text)
//       result = []
//       last_break = 0
//       current_sentence_start = 0
//       while match = @language_vars.re_period_context.match(text, last_break)
//         context = match[0] + match[:after_tok]
//         if text_contains_sentence_break?(context)
//           result << [current_sentence_start, (match.end(0)-1)]
//           match[:next_tok] ? current_sentence_start = match.begin(:next_tok) : current_sentence_start = match.end(0)          
//         end
//         if match[:next_tok]
//           last_break = match.begin(:next_tok)
//         else
//           last_break = match.end(0)
//         end
//       end
//       result << [current_sentence_start, (text.size-1)]
//     end
    
//     def text_contains_sentence_break?(text)
//       found = false
//       annotate_tokens(tokenize_words(text)).each do |token|
//         return true if found
//         found = true if token.sentence_break
//       end
//       return false
//     end
    
//     def annotate_tokens(tokens)
//       tokens = annotate_first_pass(tokens)
//       tokens = annotate_second_pass(tokens)
//       return tokens
//     end
    
//     def annotate_second_pass(tokens)
//       pair_each(tokens) do |tok1, tok2|
//         next unless tok2
//         next unless tok1.ends_with_period?
        
//         token            = tok1.token
//         type             = tok1.type_without_period
//         next_token       = tok2.token
//         next_type        = tok2.type_without_sentence_period
//         token_is_initial = tok1.is_initial?

//         if @parameters.collocations.include?([type, next_type])
//           tok1.sentence_break = false
//           tok1.abbr           = true
//           next
//         end

//         if (tok1.abbr || tok1.ellipsis) && !token_is_initial
//           is_sentence_starter = orthographic_heuristic(tok2)
//           if is_sentence_starter == true
//             tok1.sentence_break = true
//             next
//           end
          
//           if tok2.first_upper? && @parameters.sentence_starters.include?(next_type)
//             tok1.sentence_break = true
//             next
//           end
//         end
        
//         if token_is_initial || type == "##number##"
//           is_sentence_starter = orthographic_heuristic(tok2)
//           if is_sentence_starter == false
//             tok1.sentence_break = false
//             tok1.abbr           = true
//             next
//           end
          
//           if is_sentence_starter == :unknown && token_is_initial &&
//              tok2.first_upper? && !(@parameters.orthographic_context[next_type] & Punkt::ORTHO_LC != 0)
//              tok1.sentence_break = false
//              tok1.abbr           = true
//           end
//         end
        
//       end
//       return tokens
//     end
    
//     def orthographic_heuristic(aug_token)
//       return false if [';', ',', ':', '.', '!', '?'].include?(aug_token.token)
      
//       orthographic_context = @parameters.orthographic_context[aug_token.type_without_sentence_period]
//       return true if aug_token.first_upper? && (orthographic_context & Punkt::ORTHO_LC != 0) && !(orthographic_context & Punkt::ORTHO_MID_UC != 0)
//       return false if aug_token.first_lower? && ((orthographic_context & Punkt::ORTHO_UC != 0) || !(orthographic_context & Punkt::ORTHO_BEG_LC != 0))
//       return :unknown
//     end
  
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
