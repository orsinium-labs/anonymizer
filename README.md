# anonymizer

Go package for anonymizing text. It removes all kinds of PII: names, places, phone numbers, etc.

The main design principle is "better safe than sorry": if it's not sure if a word should be anonymized, it gets anonymized. It includes all non-dictionary words and words starting with a capital letter (which aren't at the beginning of a sentence).

## Example

Input:

> Good morning, doctor. My name is Gram. I live in amsterdam, at kerkstraat 42. My social number is 123-456.

Output:

> Good morning, doctor. My name is Xxxx. I live in xxxxxxxxx, at xxxxxxxxxx 00. My social number is 000-000.

## Installation

```bash
go get github.com/orsinium-labs/anonymizer
```

Make sure you have dictionaries installed for the language you're going to anonymize. For example, for American English:

```bash
sudo apt install wamerican
```

To list dictionaries that you already have installed:

```bash
ls /usr/share/dict
```

To list all dictionaries that can be installed:

```bash
sudo apt install aptitude
aptitude search '?provides(wordlist)'
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
