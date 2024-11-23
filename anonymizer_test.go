package anonymizer //nolint:testpackage

import (
	"testing"

	"github.com/matryer/is"
)

func TestAnonymize(t *testing.T) {
	t.Parallel()
	is := is.New(t)
	is.Equal(Anonymize("Albert"), "Xxxxxx")
	is.Equal(Anonymize("hoi Albert!"), "hoi Xxxxxx!")
	is.Equal(Anonymize("Hoi Albert!"), "Hoi Xxxxxx!")
	is.Equal(Anonymize("Hoi Hoi!"), "Hoi Xxx!")
	is.Equal(Anonymize("hoe gaat het, asdasd?"), "hoe gaat het, xxxxxx?")
	is.Equal(Anonymize("MyLoKo"), "XxXxXx")
	is.Equal(Anonymize("12 34"), "00 00")
	is.Equal(Anonymize("12-34"), "00-00")
	is.Equal(Anonymize("Lelystad"), "Xxxxxxxx")
	is.Equal(Anonymize("Flёur"), "Xxxxx")
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

func TestSentenceStart(t *testing.T) {
	t.Parallel()
	is := is.New(t)
	is.True(isSentenceStart([]rune("hello"), 0))
	is.True(isSentenceStart([]rune(". Hello"), 2))
	is.True(isSentenceStart([]rune("? Hello"), 2))
	is.True(isSentenceStart([]rune("⁉ Hello"), 2))
	is.True(isSentenceStart([]rune("h.Hello"), 2))
	is.True(isSentenceStart([]rune("..Hello"), 2))
	is.True(isSentenceStart([]rune(".  Hello"), 3))

	is.True(!isSentenceStart([]rune(", Hello"), 2))
	is.True(!isSentenceStart([]rune("' Hello"), 2))
	is.True(!isSentenceStart([]rune(".,Hello"), 2))
}
