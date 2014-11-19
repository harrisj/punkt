package punkt

import (
	//"fmt"
	"regexp"
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

type OrthoPair struct {
	First  string
	Second string
}

// A map from context position and first-letter case to the
// appropriate orthographic context flag.
// const(
//   ORTHO_MAP = map[OrthoPair]OrthoContext {
//     OrthoPair{"initial","upper"}: ORTHO_BEG_UC,
//     OrthoPair{"internal","upper"}: ORTHO_MID_UC,
//     OrthoPair{"unknown","upper"}: ORTHO_UNK_UC,
//     OrthoPair{"initial","lower"}: ORTHO_BEG_LC,
//     OrthoPair{"internal","lower"}: ORTHO_MID_LC,
//     OrthoPair{"unknown","lower"}: ORTHO_UNK_LC
//   }
// )

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
func WordTokenize(input string) (tokens []string) {
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
