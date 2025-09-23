// Copyright 2025 Abdulrahman Abdulhamid. All rights reserved.
// Use of this source code is governed by Apache-2.0 
// license that can be found in the LICENSE file.

package ints

import(
	"unicode"
	"github.com/begopher/peruse/internal/numeric"
)

func Unsigned() numeric.Number {
	return unsigned{}
}

type unsigned struct {}

func (d unsigned) Scan(content []rune) []rune {
	offset := 0
	for _, d := range content {
		if d == ' ' || d == ')' || d == '\n' {
			break
		}
		if !unicode.IsDigit(d) {
			return []rune{}
		}
		offset++
	}
	if offset == 0 {
		return []rune{}
	}
	return content[:offset]	
}

