package anonymizer

import (
	"testing"

	"github.com/matryer/is"
)

func TestAnonymize_Default(t *testing.T) {
	t.Parallel()
	is := is.New(t)
	a := New(nil)
	is.Equal(a.Anonymize("Albert"), "█▄▄▄▄▄")
	is.Equal(a.Anonymize("Albert hello"), "█▄▄▄▄▄ hello")
	is.Equal(a.Anonymize("Hello Albert"), "Hello █▄▄▄▄▄")
	is.Equal(a.Anonymize("Hello, Albert. How are you?"), "Hello, █▄▄▄▄▄. How are you?")
	is.Equal(a.Anonymize("MyLoKo"), "█▄█▄█▄")
	is.Equal(a.Anonymize("12 34"), "00 00")
	is.Equal(a.Anonymize("12-34"), "00-00")
	is.Equal(a.Anonymize("Lelystad"), "█▄▄▄▄▄▄▄")
	is.Equal(a.Anonymize("Flёur"), "█▄▄▄▄")
}

func TestAnonymize_Dutch(t *testing.T) {
	t.Parallel()
	is := is.New(t)
	d, err := LoadDict("nl")
	is.NoErr(err)
	a := New(d)

	is.Equal(a.Anonymize("Albert"), "█▄▄▄▄▄")
	is.Equal(a.Anonymize("hoi Albert!"), "hoi █▄▄▄▄▄!")
	is.Equal(a.Anonymize("Hoi Albert!"), "Hoi █▄▄▄▄▄!")
	is.Equal(a.Anonymize("Hoi Hoi!"), "Hoi █▄▄!")
	is.Equal(a.Anonymize("hoe gaat het, asdasd?"), "hoe gaat het, ▄▄▄▄▄▄?")

	is.Equal(a.Anonymize("MyLoKo"), "█▄█▄█▄")
	is.Equal(a.Anonymize("12 34"), "00 00")
	is.Equal(a.Anonymize("12-34"), "00-00")
	is.Equal(a.Anonymize("Lelystad"), "█▄▄▄▄▄▄▄")
	is.Equal(a.Anonymize("Flёur"), "█▄▄▄▄")
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
