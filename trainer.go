package punkt

import (
  "regexp"
  "strings"
  "fmt"
  "math"
)

const (
  // cut-off value whether a 'token' is an abbreviation
  ABBREV_CUTOFF = 0.3

  // allows the disabling of the abbreviation penalty heuristic, which
  // exponentially disadvantages words that are found at times without a
  // final period.
  IGNORE_ABBREV_PENALTY = false

  // upper cut-off for Mikheev's(2002) abbreviation detection algorithm
  ABBREV_BACKOFF = 5

  // minimal log-likelihood value that two tokens need to be considered as a colocation
  COLLOCATION_CUTOFF = 7.88

  // minimal log-likelihood value that a token requires to be considered as a FrequentSentence starter
  SENT_STARTER_CUTOFF = 30

  // this includes as potential collocations all word pairs where the first
  // word ends in a period. It may be useful in corpora where there is a lot
  // of variation that makes abbreviations like Mr difficult to identify.
  INCLUDE_ALL_COLLOCS = true

  // this includes as potential collocations all word pairs where the first
  // word is an abbreviation. Such collocations override the orthographic
  // heuristic, but not the sentence starter heuristic. This is overridden by
  // INCLUDE_ALL_COLLOCS, and if both are false, only collocations with initials
  // and ordinals are considered.
  INCLUDE_ABBREV_COLLOCS = false

  // this sets a minimum bound on the number of times a bigram needs to
  // appear before it can be considered a collocation, in addition to log
  // likelihood statistics. This is useful when INCLUDE_ALL_COLLOCS is True.
  MIN_COLLOC_FREQ = 1
)

type Trainer struct {
  TypeFdist            FrequencyDistribution
  CollocationFdist     FrequencyDistribution
  SentenceStarterFdist FrequencyDistribution
  PeriodTokensCount    int
  SentenceBreakCount   int
  Finalized            bool
}

type AbbrevClassification struct {
  Type string
  Score float64
  IsAdd bool
}

func (t *Trainer) TrainWithText(text string) *LanguageParameters {
  tokens := TokenizeTextToWords(text)
  return t.TrainWithTokenizedText(tokens)
}

func (t *Trainer) TrainWithTokenizedText(textTokens []string) *LanguageParameters {
  tokens := make([]Token, len(textTokens))
  for i := range textTokens {
    tokens[i] = *MakeToken(textTokens[i])
  }

  return t.trainFromTokens(tokens)
}

// private methods
func (t *Trainer) trainFromTokens(tokens []Token) *LanguageParameters {
  uniqueTypes := map[string]bool{}
  parameters := new(LanguageParameters)

  for _, tok := range tokens {
    t.TypeFdist.Inc(tok.Type)
    uniqueTypes[tok.Type] = true

    if tok.EndsWithPeriod() {
      t.PeriodTokensCount += 1
    }
  }

  // reclassify abbeviation types
  abbr_types := t.ReclassifyAbbreviationTypes(parameters, uniqueTypes)

  for _, ac := range abbr_types {
    if ac.Score >= ABBREV_CUTOFF {
      if ac.IsAdd {
        parameters.SaveAbbrevType(ac.Type)
      }  
    } else {
      if !(ac.IsAdd) {
        parameters.DeleteAbbrevType(ac.Type)
      }
    }
  }

  tokens = AnnotateFirstPass(parameters, tokens)

  t.BuildOrthographyTables(parameters, tokens)

  for _, tok := range tokens {
    if tok.IsSentenceBreak() {
      t.SentenceBreakCount += 1
    }
  }

  // in pairs
  for i, tok2 := range tokens {
    if i == 0 {
      continue
    }

    tok1 := tokens[i-1]

    if !(tok1.EndsWithPeriod()) {
      continue
    }

    if t.IsRareAbbrevType(parameters, tok1, tok2) {
      parameters.SaveAbbrevType(tok1.TypeWithoutPeriod())
    }

    // if t.PotentialSentenceStarter(tok1, tok2) {
    //   SentenceStarterFdist.Inc(tok2.Type)
    // }

    // if t.PotentialCollocation(tok1, tok2) {
    //   colloc_key := fmt.Sprintf("%v|%v", tok1.TypeWithoutPeriod, tok2.TypeWithoutPeriod)
    //   CollocationFdist.Inc(colloc_key)
    // }
  }

  return parameters
}

func (t *Trainer) ReclassifyAbbreviationTypes(parameters *LanguageParameters, uniqueTypes map[string]bool) (out []AbbrevClassification) {
  punctRegexp := regexp.MustCompile("[^\\W\\d]")
  isAdd := false

  for key, _ := range uniqueTypes {
    // if there is punctuation or is a number, continue. This will be processed later
    if key == "##number##" || punctRegexp.MatchString(key) {
      continue
    }

    if strings.HasSuffix(key, ".") {
      if parameters.HasAbbrevType(key) { // FIXME
        continue
      }

      key = key[0:len(key)]  // chop
      isAdd = true
    } else {
      if !(parameters.HasAbbrevType(key)) {
        continue
      }

      isAdd = false
    }

    periodsCount := strings.Count(key, ".") + 1
    nonPeriodsCount := len(key) - periodsCount + 1

    withPeriodsCount := t.TypeFdist.Get(fmt.Sprintf("%s.", key))
    withoutPeriodsCount := t.TypeFdist.Get(key)

    ll := DunningLogLikelihood(withPeriodsCount + withoutPeriodsCount, t.PeriodTokensCount, withPeriodsCount, t.TypeFdist.N)

    fLength  := math.Exp(float64(-nonPeriodsCount))
    fPeriods := periodsCount
    fPenalty := float64(0)

    if !(IGNORE_ABBREV_PENALTY) {
      fPenalty = math.Pow(float64(nonPeriodsCount), float64(-withoutPeriodsCount))
    }

    score := ll * fLength * float64(fPeriods) * fPenalty

    out = append(out, AbbrevClassification{key, score, isAdd})
  }

  return out
}

func DunningLogLikelihood(count_a, count_b, count_ab, n int) float64 {
  p1 := float64(count_b) / float64(n)
  p2 := 0.99

  null_hypo := float64(count_ab) * math.Log(p1) + float64(count_a - count_ab) * math.Log(1.0 - p1)
  alt_hypo := float64(count_ab) * math.Log(p2) + float64(count_a - count_ab) * math.Log(1.0 - p2)

  likelihood := null_hypo - alt_hypo
  return (-2.0 * likelihood)
}

func (t *Trainer) IsRareAbbrevType(parameters *LanguageParameters, current Token, next Token) bool {
  if current.IsAbbr() || !(current.IsSentenceBreak()) {
    return false
  }

  tType := current.TypeWithoutSentencePeriod()
  count := t.TypeFdist.Get(tType) + t.TypeFdist.Get(tType[0:len(tType)])

  if parameters.HasAbbrevType(tType) || count >= ABBREV_BACKOFF {
    return false
  }

  if next.Value[0] == ',' || next.Value[0] == ':' || next.Value[0] == ';' {
    return true
  } else if next.FirstLower() {
    tType2 := next.TypeWithoutSentencePeriod()
    tType2Ortho := parameters.GetOrthographicContext(tType2)

    return (tType2Ortho & ORTHO_BEG_UC != 0) && (tType2Ortho & ORTHO_MID_UC != 0)
  } else {
    return false
  }
}

func (t *Trainer) BuildOrthographyTables(parameters *LanguageParameters, tokens []Token) {
  context := "internal"

  for _, tok := range tokens {
    if tok.IsParagraphStart() && context != "unknown" {
      context = "initial"
    }

    if tok.IsLineStart() && context == "internal" {
      context = "unknown"
    }

    tType := tok.TypeWithoutSentencePeriod()
    orthoKey := fmt.Sprintf("%s|%s", tType, tok.FirstCase())
    flag := ORTHO_MAP[orthoKey]

    if flag > 0 {
      parameters.AddOrthographicContext(tType, flag)
    }

    if tok.IsSentenceBreak() {
      if !(tok.MatchNumber() || tok.MatchInitial()) {
        context = "initial"
      } else {
        context = "unknown"
      }
    } else if tok.IsEllipsis() || tok.IsAbbr() {
      context = "unknown"
    } else {
      context = "internal"
    }
  }
}


func (t *Trainer) FinalizeTraining(parameters *LanguageParameters) {
  // clear sentence starters

  // clear collocations

  t.Finalized = true
}
        

    
//     def train(text_or_tokens)
//       if text_or_tokens.kind_of?(String)
//         tokens = tokenize_words(text_or_tokens) 
//       elsif text_or_tokens.kind_of?(Array)
//         tokens = text_or_tokens.map { |t| @token_class.new(t) }
//       end
//       train_tokens(tokens)
//     end
    
//     def parameters
//       finalize_training unless @finalized
//       return @parameters
//     end
    
//     def finalize_training
//       @parameters.clear_sentence_starters 
//       find_sentence_starters do |type, ll|
//         @parameters.sentence_starters << type
//       end
      
//       @parameters.clear_collocations
//       find_collocations do |types, ll|
//         @parameters.collocations << [types[0], types[1]]
//       end

//       @finalized = true
//     end
    
//   private 
  

    
//     def reclassify_abbreviation_types(types, &block)
//       types.each do |type|
//         # if there is punctuation or is a number, continue. This will be processed later
//         next if (type =~ /[^\W\d]/).nil? || type == "##number##" 
          
//         if type.end_with?(".")
//           next if @parameters.abbreviation_types.include?(type)
//           type = type.chop
//           is_add = true
//         else
//           next unless @parameters.abbreviation_types.include?(type)
//           is_add = false
//         end
        
//         periods_count = type.count(".") + 1
//         non_periods_count = type.size - periods_count + 1
        
//         with_periods_count     = @type_fdist[type + "."]
//         without_periods_count  = @type_fdist[type]
        
//         ll = dunning_log_likelihood(with_periods_count + without_periods_count,
//                                     @period_tokens_count, 
//                                     with_periods_count,
//                                     @type_fdist.N)
        
//         f_length  = Math.exp(-non_periods_count)
//         f_periods = periods_count
//         f_penalty = IGNORE_ABBREV_PENALTY ? 0 : non_periods_count**(-without_periods_count).to_f
        
//         score = ll * f_length * f_periods * f_penalty
        
//         yield(type, score, is_add)
//       end
//     end
    
//     def dunning_log_likelihood(count_a, count_b, count_ab, n)
//       p1 = count_b.to_f / n
//       p2 = 0.99

//       null_hypo = (count_ab.to_f * Math.log(p1) +
//                    (count_a - count_ab) * Math.log(1.0 - p1))
//       alt_hypo  = (count_ab.to_f * Math.log(p2) +
//                    (count_a - count_ab) * Math.log(1.0 - p2))

//       likelihood = null_hypo - alt_hypo

//       return (-2.0 * likelihood)
//     end
    
//     def get_orthography_data(tokens)
//       context = :internal
      
//       tokens.each do |aug_token|
//         context = :initial if aug_token.paragraph_start && context != :unknown
//         context = :unknown if aug_token.line_start && context == :internal
        
//         type = aug_token.type_without_sentence_period

//         flag = Punkt::ORTHO_MAP[[context, aug_token.first_case]] || 0
//         @parameters.add_orthographic_context(type, flag) if flag
        
//         if aug_token.sentence_break
//           context = !(aug_token.is_number? || aug_token.is_initial?) ? :initial : :unknown
//         elsif aug_token.ellipsis || aug_token.abbr
//           context = :unknown
//         else
//           context = :internal
//         end
//       end
//     end
    
//     def is_rare_abbreviation_type?(current_token, next_token)
//       return false if current_token.abbr || !current_token.sentence_break
      
//       type = current_token.type_without_sentence_period
      
//       count = @type_fdist[type] + @type_fdist[type.chop]
//       return false if (@parameters.abbreviation_types.include?(type) || count >= ABBREV_BACKOFF)

//       if @language_vars.internal_punctuation.include?(next_token.token[0])
//         return true
//       elsif next_token.first_lower?
//         type2 = next_token.type_without_sentence_period
//         type2_orthographic_context = @parameters.orthographic_context[type2]
//         return true if (type2_orthographic_context & Punkt::ORTHO_BEG_UC != 0) && (type2_orthographic_context & Punkt::ORTHO_MID_UC != 0)
//       end      
//     end
    
//     def is_potential_sentence_starter?(current_token, previous_token)
//       return (previous_token.sentence_break && 
//               !(previous_token.is_number? || previous_token.is_initial?) && 
//               current_token.is_alpha?)
//     end
    
//     def is_potential_collocation?(tok1, tok2)
//       return ((INCLUDE_ALL_COLLOCS || 
//                   (INCLUDE_ABBREV_COLLOCS && tok1.abbr) || 
//                   (tok1.sentence_break && 
//                     (tok1.is_number? || tok2.is_initial?))) &&
//                 tok1.is_non_punctuation? &&
//                 tok2.is_non_punctuation?)
//     end
    
//     def find_sentence_starters(&block)
//       @sentence_starter_fdist.each do |type, type_at_break_count|
//         next if !type
        
//         type_count = @type_fdist[type] + @type_fdist[type + "."]
        
//         next if type_count < type_at_break_count
        
//         ll = col_log_likelihood(@sentence_break_count, 
//                                 type_count, 
//                                 type_at_break_count, 
//                                 @type_fdist.N)
              
//         if (ll >= SENT_STARTER && 
//            @type_fdist.N.to_f/@sentence_break_count > type_count.to_f/type_at_break_count)
//           yield(type, ll)
//         end
//       end
//     end
    
//     def col_log_likelihood(count_a, count_b, count_ab, n)
//       p = 1.0 * count_b / n
//       p1 = 1.0 * count_ab / count_a
//       p2 = 1.0 * (count_b - count_ab) / (n - count_a)

//       summand1 = (count_ab * Math.log(p) +
//                   (count_a - count_ab) * Math.log(1.0 - p))

//       summand2 = ((count_b - count_ab) * Math.log(p) +
//                   (n - count_a - count_b + count_ab) * Math.log(1.0 - p))

//       if count_a == count_ab
//           summand3 = 0
//       else
//           summand3 = (count_ab * Math.log(p1) +
//                       (count_a - count_ab) * Math.log(1.0 - p1))
//       end

//       if count_b == count_ab
//           summand4 = 0
//       else
//           summand4 = ((count_b - count_ab) * Math.log(p2) +
//                       (n - count_a - count_b + count_ab) * Math.log(1.0 - p2))
//       end

//       likelihood = summand1 + summand2 - summand3 - summand4

//       return (-2.0 * likelihood)
//     end
    
//     def find_collocations(&block)
//       @collocation_fdist.each do |types, col_count|
//         type1, type2 = types
        
//         next if type1.nil? || type2.nil?
//         next if @parameters.sentence_starters.include?(type2)
        
//         type1_count = @type_fdist[type1] + @type_fdist[type1 + "."]
//         type2_count = @type_fdist[type2] + @type_fdist[type2 + "."]
        
//         if (type1_count > 1 && type2_count > 1 &&
//             MIN_COLLOC_FREQ < col_count &&
//             col_count <= [type1_count, type2_count].min)
          
//           ll = col_log_likelihood(type1_count, type2_count,
//                                   col_count, @type_fdist.N)
                                  
//           if (ll >= COLLOCATION &&
//               @type_fdist.N.to_f/type1_count > type2_count.to_f/col_count)
//             yield([type1, type2], ll)
//           end
//         end
//       end
//     end
    
//   end
// end
