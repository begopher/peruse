// Copyright 2025 Abdulrahman Abdulhamid. All rights reserved.
// Use of this source code is governed by Apache-2.0 
// license that can be found in the LICENSE file.

package floats

import(
	"unicode"
	"github.com/begopher/peruse/internal/numeric"
)

func Unsigned() numeric.Number {
	return unsigned{}
}

type unsigned struct {}

func (unsigned) Scan(content []rune) []rune {
	required := 2
	if len(content) < required {
		return []rune{}
	}	
	offset, dots := 0, 0
	for _, d := range content {
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
