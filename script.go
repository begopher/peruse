// Copyright 2025 Abdulrahman Abdulhamid. All rights reserved.
// Use of this source code is governed by Apache-2.0 
// license that can be found in the LICENSE file.

package peruse

import(
	"unicode"
	"github.com/begopher/peruse/internal/numeric"
	"github.com/begopher/peruse/internal/numeric/ints"
	"github.com/begopher/peruse/internal/numeric/floats"	
)

var integer = numeric.Numbers(ints.Signed(), ints.Unsigned())
var float = numeric.Numbers(floats.Signed(), floats.Unsigned())

func Script(origin, content string) Text {
	return &text{
		origin: origin,
		integer: integer,
		float: float,
		line: 1,
		lineReset: 1,
		column: 1,
		columnReset: 1,
		content: []rune(content),
	}
}

type text struct {
	origin string
	integer numeric.Number
	float numeric.Number
	line int
	lineReset int
	column int
	columnReset int
	content []rune
}

func (t *text) Origin() string {
	return t.origin
}

func (t *text) Column() int {
	return t.column
}

func (t *text) Line() int {
	return t.line
}

func (t *text) Location() Location{
	return NewLocation(t.origin, t.line, t.column)	
}

func (t *text) Length() int {
	return len(t.content)
}

func (t *text) Empty() bool {
	return len(t.content) == 0
}


func (t *text) Remain() string {
	return string(t.content)
}

func (t *text) EatSpaces() {
	for _, r := range t.content {
		if !unicode.IsSpace(r) {
			break
		}
		if '\n' == r {
			t.line++
			t.column = t.columnReset
			t.content = t.content[1:]
			continue
		}
		if ' ' == r || '\t' == r {
			t.column++
			t.content = t.content[1:]		
			continue
		}
		t.content = t.content[1:]
		continue
	}	
}

func (t *text) BeginWith(prefix string) bool {
	if prefix == "" {
		return false
	}
	sub := []rune(prefix)
	if len(sub) > len(t.content) {
		return false
	}
	return string(t.content[:len(sub)]) == prefix
}


func (t *text) Eat(prefix string) bool {
	if !t.BeginWith(prefix) {
		return false
	}
	column, line := t.column, t.line
	for _, r := range prefix {
		if r == '\n' {
			line++
			column = t.columnReset
			continue
		}
		column++
	}
	sub := []rune(prefix)
	t.content = t.content[len(sub):]
	t.column = column
	t.line = line
	return true	
}

func (t *text) EatSymbol(symbol string) bool {
	if t.Eat("("+symbol+" ") {
		return true
	}
	if t.Eat("("+symbol+"\n") {
		return true
	}
	if t.BeginWith("("+symbol+")") && t.Eat("("+symbol) {
		return true
	}
	if len("("+symbol) == t.Length() && t.Eat("("+symbol) {
		return true
	}
	return false
}

func (t *text) EatString() (string, bool) {
	if t.Length() == 0 {
		return "", false
	}
	if t.content[0] != '"' {
		return "", false
	}
	//line := s.lineReset
	line := 0
	col, offset := 2, 2	// for first(") and last(")
	buffer := t.content[1:]
	closed := false
	for i := 0; i<len(buffer); i++ {
		r := buffer[i]
		if r == '\n' {
			offset++
			line++
			col = t.columnReset			
			continue
		}
		if r == '\\' && i < len(buffer) && buffer[i+1] == '"' {
			offset+=2
			col+=2
			//col+=1
			i++			
			continue
		}
		if r == '"' {
			//offset++
			//col++
			closed = true
			break
		}
		offset++
		col++
	}
	if !closed {
		return "", false
	}
	result := t.content[1:offset-1]
	t.content = t.content[offset:]	
	if line == 0 {		
		t.column += col
	} else {
		t.line += line
		t.column = col
	}
	return string(result), true
}

func (t *text) EatWord() string {
	if len(t.content) == 0 {
		return ""
	}
	if !unicode.IsLetter(t.content[0]) {
		return ""
	}
	offset := 0
	for _, r := range t.content {
		if r == ' ' || r == ')' || r =='\n' {
			break
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return ""
		}					
		offset++
	}
	t.column += offset
	result := t.content[:offset]
	t.content = t.content[offset:]
	return string(result)
}

func (t *text) EatKeyword() string {
	if len(t.content) < 2 {
		return ""
	}
	if t.content[0] != ':' {
		return ""
	}
	if !unicode.IsLetter(t.content[1]) {
		return ""
	}
	offset := 2
	content := t.content[offset:]
	for _, r := range content {
		if r == ' ' {
			break
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return ""
		}					
		offset++
	}
	t.column += offset
	result := t.content[:offset]
	t.content = t.content[offset:]
	return string(result)
}

func (t *text) EatInteger() string {
	var digits []rune = t.integer.Scan(t.content)
	if len(digits) == 0 {
		return ""
	}
	t.content = t.content[len(digits):]
	t.column+= len(digits)
	return string(digits)
}

func (t *text) EatFloat() string {
	var digits []rune = t.float.Scan(t.content)
	if len(digits) == 0 {
		return ""
	}
	t.content = t.content[len(digits):]
	t.column+= len(digits)
	return string(digits)
}
