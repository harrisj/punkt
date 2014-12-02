package punkt

import (
	//"fmt"
	"regexp"
	"strings"
)

// for debugging reasons why it exits
const (
	REASON_DEFAULT_DECISION                            = "default decision"
	REASON_KNOWN_COLLOCATION                           = "known collocation (both words)"
	REASON_ABBR_WITH_ORTHOGRAPHIC_HEURISTIC            = "abbreviation + orthographic heuristic"
	REASON_ABBR_WITH_SENTENCE_STARTER                  = "abbreviation + frequent sentence starter"
	REASON_INITIAL_WITH_ORTHOGRAPHIC_HEURISTIC         = "initial + orthographic heuristic"
	REASON_NUMBER_WITH_ORTHOGRAPHIC_HEURISTIC          = "number + orthographic heuristic"
	REASON_INITIAL_WITH_SPECIAL_ORTHOGRAPHIC_HEURISTIC = "initial + special orthographic heuristic"
)

// WordTokenization code
type WordTokenizeState uint8

const (
	NONE WordTokenizeState = iota
	WORD
	PUNCT
)

// The original Python and Ruby versions use PCRE regexps with lookaheads. This functionality is not supported by Go's native Regexp
// class and I don't want to compile in a PCRE library, so I wrote this as a crude state machine. This code could no doubt be optimized
// but at least it has tests unlike the original Regexp in word_tokenize_test.go
func SplitTextIntoWords(input string) (tokens []string) {
	RE_WORD_START := regexp.MustCompile("^[^\\(\"\\`{\\s\\[:;&\\#\\*@\\)}\\]\\-,]")
	RE_NON_WORD := regexp.MustCompile("^[\\?!\\)\";}\\]\\*:@'\\({\\[,]")
	RE_WHITESPACE := regexp.MustCompile("^\\s+")
	RE_NONWHITESPACE := regexp.MustCompile("^\\S")
	RE_SPECIAL_PUNCT := regexp.MustCompile("^[\\.\\-]")

	curState := NONE
	stateIndex := 0
	v := ""
	last_v := ""
	//oldState := NONE
	
	for i, r := range input {
		v = string(r)

		//oldState = curState

		switch curState {
		case NONE:
			if RE_WORD_START.MatchString(v) || RE_SPECIAL_PUNCT.MatchString(v) {
				curState = WORD
				stateIndex = i
			} else if RE_NONWHITESPACE.MatchString(v) || RE_SPECIAL_PUNCT.MatchString(v) {
				tokens = append(tokens, input[i:i+1])
			}

		case WORD:
			if RE_SPECIAL_PUNCT.MatchString(v) {
				// if only a single - or ., treat as part of word. Otherwise, dump any word and switch to punct state
				if v == last_v {
					if stateIndex < i - 1 {
						tokens = append(tokens, input[stateIndex:i-1])
					}

					curState = PUNCT
					stateIndex = i-1
				}
			} else if RE_NON_WORD.MatchString(v) || RE_WHITESPACE.MatchString(v) {
				tokens = append(tokens, input[stateIndex:i])
				curState = NONE

				// append punct
				if RE_NON_WORD.MatchString(v) {
					tokens = append(tokens, input[i:i+1])
				}
			}

		case PUNCT:
			if !(RE_SPECIAL_PUNCT.MatchString(v)) {
				tokens = append(tokens, input[stateIndex:i])

				if RE_WORD_START.MatchString(v) {
					curState = WORD
					stateIndex = i
				} else {
					if !(RE_WHITESPACE.MatchString(v)) {
						tokens = append(tokens, input[i:i+1])
					}

					curState = NONE
				}
			}
		}


		//fmt.Printf("%v %s %s->%s %v\n", i, v, oldState, curState, stateIndex)
		last_v = v
	}

	if curState == WORD || curState == PUNCT {
		tokens = append(tokens, input[stateIndex:])
	}

	return
}

func TokenizeText(plainText string) []*Token {
	paragraphStart := false

	lines := strings.Split(plainText, "\n")
	out := make([]*Token, 0)

	for _, line := range lines {
		if len(line) == 0 {
			paragraphStart = true
		} else {
			words := SplitTextIntoWords(line)

			for j, v := range words {
				t := MakeToken(v)

				if j == 0 {
					t.SetParagraphStart(paragraphStart)
					t.SetLineStart(true)
				}

				paragraphStart = false
				out = append(out, t)
			}
		}
	}

	return out 
}
