# anonymizer

Go package for anonymizing text. It removes all kinds of PII: names, places, phone numbers, etc.

The main design principle is "better safe than sorry": if it's not sure if a word should be anonymized, it gets anonymized. It includes all non-dictionary words and words starting with a capital letter (which aren't at the beginning of a sentence).

## Example

Input:

> Good morning, doctor. My name is Gram. I live in amsterdam, at kerkstraat 42. My social number is 123-456.

Output:

> Good morning, doctor. My name is █▄▄▄. I live in ▄▄▄▄▄▄▄▄▄, at ▄▄▄▄▄▄▄▄▄▄ 00. My social number is 000-000.

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

If the language is not found or not provided, the default one will be used. Run `sudo select-default-wordlist` to change the system default.

## Usage

```go
input := "Hi, my name is Gram."
dict, err := anonymizer.LoadDict("en")
if err != nil {
    panic(err)
}
a := anonymizer.New(dict)
a.Language = "en"
output := a.Anonymize(input)
fmt.Println(output)
// Output: Hi, my name is Xxxx.
```
