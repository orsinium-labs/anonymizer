package anonymizer

import (
	"testing"

	"github.com/matryer/is"
)

func TestAnonymize(t *testing.T) {
	t.Parallel()
	is := is.New(t)
	a := New()
	is.Equal(a.Anonymize("Albert"), "Xxxxxx")
	is.Equal(a.Anonymize("hoi Albert!"), "hoi Xxxxxx!")
	is.Equal(a.Anonymize("Hoi Albert!"), "Hoi Xxxxxx!")
	is.Equal(a.Anonymize("Hoi Hoi!"), "Hoi Xxx!")
	is.Equal(a.Anonymize("hoe gaat het, asdasd?"), "hoe gaat het, xxxxxx?")
	is.Equal(a.Anonymize("MyLoKo"), "XxXxXx")
	is.Equal(a.Anonymize("12 34"), "00 00")
	is.Equal(a.Anonymize("12-34"), "00-00")
	is.Equal(a.Anonymize("Lelystad"), "Xxxxxxxx")
	is.Equal(a.Anonymize("Flёur"), "Xxxxx")
}

func TestIterWords(t *testing.T) {
	t.Parallel()
	is := is.New(t)

	listWords := func(text string) []string {
		words := make([]string, 0)
		runes := []rune(text)
		for span := range iterWords(runes) {
			word := runes[span.start:span.end]
			words = append(words, string(word))
		}
		return words
	}

	is.Equal(listWords("hello world"), []string{"hello", "world"})
	is.Equal(listWords("HeLlO wOrLd"), []string{"HeLlO", "wOrLd"})
	is.Equal(listWords("привет, мир!"), []string{"привет", "мир"})
	is.Equal(listWords("Привет, мир!"), []string{"Привет", "мир"})
}

// func TestSentenceStart(t *testing.T) {
// 	t.Parallel()
// 	is := is.New(t)
// 	is.True(isSentenceStart([]rune("hello"), 0))
// 	is.True(isSentenceStart([]rune(". Hello"), 2))
// 	is.True(isSentenceStart([]rune("? Hello"), 2))
// 	is.True(isSentenceStart([]rune("⁉ Hello"), 2))
// 	is.True(isSentenceStart([]rune("h.Hello"), 2))
// 	is.True(isSentenceStart([]rune("..Hello"), 2))
// 	is.True(isSentenceStart([]rune(".  Hello"), 3))

// 	is.True(!isSentenceStart([]rune(", Hello"), 2))
// 	is.True(!isSentenceStart([]rune("' Hello"), 2))
// 	is.True(!isSentenceStart([]rune(".,Hello"), 2))
// }
