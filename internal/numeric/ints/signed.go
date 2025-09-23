// Copyright 2025 Abdulrahman Abdulhamid. All rights reserved.
// Use of this source code is governed by Apache-2.0 
// license that can be found in the LICENSE file.

package ints

import(
	"unicode"
	"github.com/begopher/peruse/internal/numeric"
)

func Signed() numeric.Number {
	return signed{}
}

type signed struct {}

func (d signed) Scan(content []rune) []rune {
	if len(content) < 2 {
		return []rune{}
	}
	signed := content[0] 
	if  signed != rune('-') && signed != rune('+') {
		return []rune{}
	}
	offset := 1
	for _, d := range content[1:] {
		if d == ' ' || d == ')' || d == '\n' {
			break
		}
		if !unicode.IsDigit(d) {
			return []rune{}
		}
		offset++
	}	
	if offset == 1 {
		return []rune{}
	}
	
	return content[:offset]	
}
