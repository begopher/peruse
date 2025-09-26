// Copyright 2025 Abdulrahman Abdulhamid. All rights reserved.
// Use of this source code is governed by Apache-2.0 
// license that can be found in the LICENSE file.

package peruse

type Text interface {
	Origin() string
	Column() int
	Line() int
	Location() Location
	Length() int
	Empty() bool
	Remain() string

	BeginWith(string) bool
	Eat(string) bool
	EatFunctionName(string) bool
	EatSpaces()
	EatString() (string, bool)
	
	EatWord() string
	EatPrefixedWord(prefix string) (word, prefixed_word string)
	// EatSuffixedWord(suffix) (string, string)
	// EatAffixedWord(prefix, suffix) (string, string)
	EatWords() (string, string)
	IsWord(string) bool
	
	EatSymbol() string
	EatPrefixedSymbol(prefix string) (symbol, prefixed_symbol string)
	EatSymbols() (string, string)
	IsSymbol(string) bool

	EatKeyword() string
	//EatKey() string
	//IsKeyword(string) bool
	
	//EatKeysymbol(string) bool
	//IsKeysymbol(string) bool
	EatInteger() string
	EatFloat() string
}


