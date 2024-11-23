# anonymizer

Go package for anonymizing text. It removes all kinds of PII: names, places, phone numbers, etc.

The main design principle is "better safe than sorry": if it's not sure if a word should be anonymized, it gets anonymized. It includes all non-dictionary words and words starting with a capital letter (which aren't at the beginning of a sentence).

Supported languages:

* NL (Dutch)

## Installation

```bash
go get github.com/orsinium-labs/anonymizer
```

## Usage

```go
input := "Hi, my name is Gram."
a := anonymizer.New()
a.Language = "en"
output := a.Anonymize(input)
fmt.Println(output)
// Output: Hi, my name is Xxxx.
```
