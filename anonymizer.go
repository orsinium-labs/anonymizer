package anonymizer

import (
	_ "embed"
	"fmt"
	"iter"
	"unicode"

	"github.com/derekparker/trie"
)

type Anonymizer struct {
	Dict      *Dict
	Uppercase rune
	Lowercase rune
	Digit     rune
}

func New(dict *Dict) Anonymizer {
	if dict == nil {
		var err error
		dict, err = LoadDict("")
		if err != nil {
			panic(fmt.Errorf("failed to load default dictionary: %v", err))
		}
	}
	return Anonymizer{
		Dict:      dict,
		Uppercase: 'X',
		Lowercase: 'x',
		Digit:     '0',
	}
}

// Replace with a placeholder all non-dictionary words in the text.
func (a Anonymizer) Anonymize(text string) string {
	runes := []rune(text)
	for i, r := range runes {
		if unicode.IsDigit(r) {
			runes[i] = a.Digit
		}
	}
	for span := range iterWords(runes) {
		if shouldAnonymize(a.Dict, span) {
			a.mask(runes, span)
		}
	}
	return string(runes)
}

// Mask the given span in the slice of runes.
func (a Anonymizer) mask(runes []rune, span span) {
	for i := span.start; i < span.end; i++ {
		r := runes[i]
		isUpper := unicode.IsUpper(r)
		if isUpper {
			runes[i] = a.Uppercase
		} else {
			runes[i] = a.Lowercase
		}
	}
}

// Check if the word in the given span should be anonymized.
func shouldAnonymize(dict *trie.Trie, span span) bool {
	word := span.word
	if unicode.IsUpper(word[0]) {
		if span.initial {
			word = toLower(word)
		} else {
			return true
		}
	}
	_, knownWord := dict.Find(string(word))
	return !knownWord
}

func toLower(word []rune) []rune {
	return append([]rune{unicode.ToLower(word[0])}, word[1:]...)
}

type span struct {
	// The index of the first rune of the word.
	start int
	// The index of the first rune after the word.
	end int
	// The star of the show, the complete word in its original case.
	word []rune
	// True if it is the first word of a sentence.
	initial bool
}

func iterWords(runes []rune) iter.Seq[span] {
	return func(yield func(span) bool) {
		start := 0
		end := 0
		terminal := -2
		for i, r := range runes {
			if unicode.IsLetter(r) {
				end = i + 1
				continue
			}
			if start < end {
				keepGoing := yield(span{
					start:   start,
					end:     end,
					word:    runes[start:end],
					initial: terminal != -1,
				})
				if !keepGoing {
					break
				}
				terminal = -1
			}
			if unicode.In(r, unicode.Sentence_Terminal) {
				terminal = i
			} else if terminal != -1 && !unicode.IsSpace(r) {
				terminal = -1
			}
			start = i + 1
		}

		if start < end {
			yield(span{
				start:   start,
				end:     end,
				word:    runes[start:end],
				initial: terminal != -1,
			})
		}
	}
}
