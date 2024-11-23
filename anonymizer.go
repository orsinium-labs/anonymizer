package anonymizer

import (
	_ "embed"
	"iter"
	"strings"
	"sync"
	"unicode"

	"github.com/derekparker/trie"
)

type Anonymizer struct {
	Dictionary *trie.Trie
	Langugage  string
	Uppercase  rune
	Lowercase  rune
	Digit      rune
}

func New() Anonymizer {
	return Anonymizer{
		Uppercase: 'X',
		Lowercase: 'x',
		Digit:     '0',
	}
}

//go:embed words/nl.txt
var dutch string

var getDict = sync.OnceValue(func() *trie.Trie {
	dict := trie.New()
	for _, word := range strings.Split(dutch, "\n") {
		dict.Add(word, nil)
	}
	dutch = ""
	return dict
})

func init() {
	go getDict()
}

// Replace with a placeholder all non-dictionary words in the text.
func (a Anonymizer) Anonymize(text string) string {
	runes := []rune(text)
	for i, r := range runes {
		if unicode.IsDigit(r) {
			runes[i] = a.Digit
		}
	}
	dict := a.Dictionary
	if dict == nil {
		dict = getDict()
	}
	for span := range iterWords(runes) {
		if shouldAnonymize(dict, runes, span) {
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
func shouldAnonymize(dict *trie.Trie, runes []rune, span span) bool {
	// If the first letter is uppercase and it's the beginning of a sentence,
	// we want to make it lowercase.
	word := span.word
	if unicode.IsUpper(word[0]) {
		if span.initial {
			word = append([]rune{unicode.ToLower(word[0])}, word[1:]...)
		}
	}
	_, knownWord := dict.Find(string(word))
	return !knownWord
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
