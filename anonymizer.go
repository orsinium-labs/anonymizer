package anonymizer

import (
	_ "embed"
	"iter"
	"strings"
	"sync"
	"unicode"

	"github.com/derekparker/trie"
)

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
func Anonymize(text string) string {
	// We work with binary instead of strings to reduce memory allocations.
	runes := []rune(text)
	for i, r := range runes {
		if unicode.IsDigit(r) {
			runes[i] = '0'
		}
	}
	dict := getDict()
	for span := range iterWords(runes) {
		if !shouldAnonymize(dict, runes, span) {
			continue
		}
		mask(runes, span)
	}
	return string(runes)
}

// Mask the given span in the slice of runes.
func mask(runes []rune, span span) {
	for i := span.start; i < span.end; i++ {
		c := runes[i]
		isUpper := unicode.IsUpper(c)
		if isUpper {
			runes[i] = 'X'
		} else {
			runes[i] = 'x'
		}
	}
}

// Check if the word in the given span should be anonymized.
func shouldAnonymize(dict *trie.Trie, runes []rune, span span) bool {
	word := runes[span.start:span.end]
	// If the first letter is uppercase and it's the beginning of a sentence,
	// we want to make it lowercase
	if unicode.IsUpper(runes[span.start]) {
		if isSentenceStart(runes, span.start) {
			word = append([]rune{unicode.ToLower(word[0])}, word[1:]...)
		}
	}
	_, knownWord := dict.Find(string(word))
	return !knownWord
}

// Detect if the character at the given position is the first in a sentence.
func isSentenceStart(runes []rune, i int) bool {
	if i < 2 {
		return true
	}
	for j := i - 1; j >= 0; j-- {
		r := runes[j]
		if !unicode.IsSpace(r) {
			return unicode.In(r, unicode.Sentence_Terminal)
		}
	}
	return true
}

type span struct {
	start int
	end   int
}

func iterWords(runes []rune) iter.Seq[span] {
	return func(yield func(span) bool) {
		start := 0
		end := 0
		for i, r := range runes {
			if unicode.IsLetter(r) {
				end = i + 1
			} else {
				if start < end {
					keepGoing := yield(span{start, end})
					if !keepGoing {
						break
					}
				}
				start = i + 1
			}
		}

		if start < end {
			yield(span{start, end})
		}
	}
}
