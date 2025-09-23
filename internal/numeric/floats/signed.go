// Copyright 2025 Abdulrahman Abdulhamid. All rights reserved.
// Use of this source code is governed by Apache-2.0 
// license that can be found in the LICENSE file.

package floats

import(
	"unicode"
	"github.com/begopher/peruse/internal/numeric"
)

func Signed() numeric.Number {
	return signed{}
}

type signed struct {}

func (signed) Scan(content []rune) []rune {
	required := 3
	if len(content) < required {
		return []rune{}
	}	
	signed := content[0] 
	if  signed != rune('-') && signed != rune('+') {
		return []rune{}
	}
	offset := 1
	dots := 0
	for _, d := range content[1:] {
		if d == '.' {
			dots++
			if dots > 1 { return []rune{} }
			offset++
			continue
		}
		if d == ' ' || d == ')' || d == '\n' {
			break
		}		
		if !unicode.IsDigit(d) {
			return []rune{}
		}		
		offset++
	}	
	if offset < required {
		return []rune{}
	}	
	return content[:offset]	
}
