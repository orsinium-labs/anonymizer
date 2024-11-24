package anonymizer

import (
	"bufio"
	"fmt"
	"os"
	"unicode"

	"github.com/derekparker/trie/v3"
)

// Dictionary of words.
type Dict = trie.Trie[struct{}]

const (
	dictDir     = "/usr/share/dict/"
	defaultPath = dictDir + "words"
)

var installed map[string]struct{}

var languages = map[string]string{
	"en": "american-english",
	"no": "bokmaal",
	"bg": "bulgarian",
	"ca": "catalan",
	"da": "danish",
	"nl": "dutch",
	"fo": "faroese",
	"fr": "french",
	"gl": "galician",
	"de": "ngerman",
	"it": "italian",
	"pl": "polish",
	"pt": "portuguese",
	"es": "spanish",
	"sv": "swedish",
	"uk": "ukrainian",
}

func init() {
	err := initInstalled()
	if err != nil {
		panic(err)
	}
}

func initInstalled() error {
	installed = make(map[string]struct{})
	files, err := os.ReadDir(dictDir)
	if err != nil {
		return fmt.Errorf("open %s dir: %v", dictDir, err)
	}
	for _, file := range files {
		installed[file.Name()] = struct{}{}
	}
	return nil
}

// A wrapper around [LoadDict] that panics on error.
func MustLoadDict(lang string) *Dict {
	dict, err := LoadDict(lang)
	if err != nil {
		panic(fmt.Errorf("load dictionary: %v", err))
	}
	return dict
}

// Load dictionary for the given language.
//
// If the language is not found or not provided,
// the default one will be used. Run `sudo select-default-wordlist`
// to change the system default.
func LoadDict(lang string) (*Dict, error) {
	path := findDict(lang)
	return loadDict(path)
}

func findDict(lang string) string {
	if lang == "" {
		return defaultPath
	}
	_, knownFile := installed[lang]
	if knownFile {
		return dictDir + lang
	}
	fileName := languages[lang]
	if fileName != "" {
		return dictDir + fileName
	}
	return defaultPath
}

func loadDict(path string) (*Dict, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %v", path, err)
	}
	defer file.Close()
	dict := trie.New[struct{}]()
	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		word := scanner.Text()
		if isLower(word) {
			dict.Add(word, struct{}{})
		}
	}
	return dict, nil
}

// Check if the given word has only lowercase letters.
//
// No uppercase, no symbols, no digits.
func isLower(word string) bool {
	for _, r := range word {
		if !unicode.IsLower(r) {
			return false
		}
	}
	return true
}
