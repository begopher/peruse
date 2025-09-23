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
	EatSymbol(string) bool 
	EatSpaces()
	EatString() (string, bool)
	EatWord() string
	EatKeyword() string
	EatInteger() string
	EatFloat() string
}
