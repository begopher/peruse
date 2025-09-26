// Copyright 2025 Abdulrahman Abdulhamid. All rights reserved.
// Use of this source code is governed by Apache-2.0 
// license that can be found in the LICENSE file.

package peruse

import(
	"unicode"
	"strings"
	"github.com/begopher/peruse/internal/numeric"
	"github.com/begopher/peruse/internal/numeric/ints"
	"github.com/begopher/peruse/internal/numeric/floats"	
)

var integer = numeric.Numbers(ints.Signed(), ints.Unsigned())
var float = numeric.Numbers(floats.Signed(), floats.Unsigned())

func Script(origin, content string) Text {
	return &script{
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

type script struct {
	origin string
	integer numeric.Number
	float numeric.Number
	line int
	lineReset int
	column int
	columnReset int
	content []rune
}

func (s *script) Origin() string {
	return s.origin
}

func (s *script) Column() int {
	return s.column
}

func (s *script) Line() int {
	return s.line
}

func (s *script) Location() Location{
	return NewLocation(s.origin, s.line, s.column)	
}

func (s *script) Length() int {
	return len(s.content)
}

func (s *script) Empty() bool {
	return len(s.content) == 0
}


func (s *script) Remain() string {
	return string(s.content)
}

func (s *script) EatSpaces() {
	for _, r := range s.content {
		if !unicode.IsSpace(r) {
			break
		}
		if '\n' == r {
			s.line++
			s.column = s.columnReset
			s.content = s.content[1:]
			continue
		}
		if ' ' == r || '\t' == r {
			s.column++
			s.content = s.content[1:]		
			continue
		}
		s.content = s.content[1:]
		continue
	}	
}

func (s *script) BeginWith(prefix string) bool {
	if prefix == "" {
		return false
	}
	sub := []rune(prefix)
	if len(sub) > len(s.content) {
		return false
	}
	return string(s.content[:len(sub)]) == prefix
}


func (s *script) Eat(prefix string) bool {
	if !s.BeginWith(prefix) {
		return false
	}
	column, line := s.column, s.line
	for _, r := range prefix {
		if r == '\n' {
			line++
			column = s.columnReset
			continue
		}
		column++
	}
	sub := []rune(prefix)
	s.content = s.content[len(sub):]
	s.column = column
	s.line = line
	return true	
}

func (s *script) EatFunctionName(name string) bool {
	if s.Eat("("+name+" ") {
		return true
	}
	if s.Eat("("+name+"\n") {
		return true
	}
	if s.BeginWith("("+name+")") && s.Eat("("+name) {
		return true
	}
	if len("("+name) == s.Length() && s.Eat("("+name) {
		return true
	}
	return false
}

func (s *script) EatString() (string, bool) {
	if s.Length() == 0 {
		return "", false
	}
	if s.content[0] != '"' {
		return "", false
	}
	//line := s.lineReset
	line := 0
	col, offset := 2, 2	// for first(") and last(")
	buffer := s.content[1:]
	closed := false
	for i := 0; i<len(buffer); i++ {
		r := buffer[i]
		if r == '\n' {
			offset++
			line++
			col = s.columnReset			
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
	result := s.content[1:offset-1]
	s.content = s.content[offset:]	
	if line == 0 {		
		s.column += col
	} else {
		s.line += line
		s.column = col
	}
	return string(result), true
}

func (s *script) EatWord() string {
	if len(s.content) == 0 {
		return ""
	}
	if !unicode.IsLetter(s.content[0]) {
		return ""
	}
	offset := 0
	for _, r := range s.content {
		if r == ' ' || r == ')' || r =='\n' {
			break
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return ""
		}					
		offset++
	}
	s.column += offset
	result := s.content[:offset]
	s.content = s.content[offset:]
	return string(result)
}

func (s *script) EatPrefixedWord(prefix string) string {
	if len(prefix) == 0 {
		return s.EatWord()
	}
	if len(s.content) <= len(prefix) { 
		return ""
	}
	if string(s.content[:len(prefix)]) != prefix {
		return ""
	}
	content := s.content[len(prefix):]
	if !unicode.IsLetter(content[0]) {
		return ""
	}
	offset := len(prefix) // +1 maybe ?
	for _, r := range content {
		if r == ' ' || r == ')' || r =='\n' {
			break
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return ""
		}					
		offset++
	}
	result := s.content[:offset]
	s.column += offset	
	s.content = s.content[offset:]
	return string(result)
}

func (s *script) EatWords() (string, string) {
	if len(s.content) == 0 {
		return "", ""
	}
	if !unicode.IsLetter(s.content[0]) {
		return "", ""
	}
	offset := 0
	colons := 0
	for _, r := range s.content {
		if r == ' ' || r == ')' || r =='\n' {
			break
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != ':' {
			return "", ""
		}
		if r == ':' {
			colons++
		}	
		offset++
	}
	result := s.content[:offset]
	words := strings.Split(string(result), ":")
	if len(words) != 2 {
		return "", ""
	}
	first  := words[0]
	second := words[1]
	if !s.IsWord(first) || !s.IsWord(second) {
		//if len(letters) == 0 || unicode.IsDigit(letters[0]) {
		return "", ""
	}
	s.column += offset
	s.content = s.content[offset:]
	return words[0], words[1]
}

func (s *script) IsWord(value string) bool {
	runes := []rune(value)
	if len(runes) == 0 {
		return false
	}
	if !unicode.IsLetter(runes[0]) {
		return false
	}	
	for _,  r := range runes {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}	
	return true	
}

func (s *script) EatSymbol() string {
	if len(s.content) == 0 {
		return ""
	}
	if !unicode.IsLetter(s.content[0]) {
		return ""
	}
	offset := 0
	for _, r := range s.content {
		if r == ' ' || r == ')' || r =='\n' {
			break
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' {
			return ""
		}					
		offset++
	}
	result := s.content[:offset]
	last := result[len(result)-1]
	if  last == '-' {
		return ""
	}
	s.column += offset
	s.content = s.content[offset:]
	return string(result)
}

func (s *script) EatPrefixedSymbol(prefix string) string {
	if len(prefix) == 0 {
		return s.EatSymbol()
	}
	if len(s.content) <= len(prefix) { 
		return ""
	}
	if string(s.content[:len(prefix)]) != prefix {
		return ""
	}
	content := s.content[len(prefix):]	
	if !unicode.IsLetter(content[0]) {
		return ""
	}
	offset := len(prefix) 
	for _, r := range content {
		if r == ' ' || r == ')' || r =='\n' {
			break
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' {
			return ""
		}					
		offset++
	}
	result := s.content[:offset]
	last := result[len(result)-1]
	if  last == '-' {
		return ""
	}
	s.column += offset
	s.content = s.content[offset:]
	return string(result)	
}

func (s *script) EatSymbols() (string, string) {
	if len(s.content) == 0 {
		return "", ""
	}
	if !unicode.IsLetter(s.content[0]) {
		return "", ""
	}
	offset := 0
	colons := 0
	for _, r := range s.content {
		if r == ' ' || r == ')' || r =='\n' {
			break
		}
		if !unicode.IsLetter(r) &&
			!unicode.IsDigit(r) &&
			r != ':' &&
			r != '-' {
			return "", ""
		}
		if r == ':' {
			colons++
		}	
		offset++
	}
	result := s.content[:offset]
	words := strings.Split(string(result), ":")
	if len(words) != 2 {
		return "", ""
	}
	first  := words[0]
	second := words[1]
	if !s.IsSymbol(first) || !s.IsSymbol(second) {
		//if len(letters) == 0 || unicode.IsDigit(letters[0]) {
		return "", ""
	}
	s.column += offset
	s.content = s.content[offset:]
	return words[0], words[1]
}

func (s *script) IsSymbol(value string) bool {
	runes := []rune(value)
	if len(runes) == 0 {
		return false
	}
	if !unicode.IsLetter(runes[0]) {
		return false
	}
	for _,  r := range runes {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' {
			return false
		}
	}	
	last := runes[len(runes)-1]
	if  last == '-' {
		return false
	}
	return true	
}

func (s *script) EatKeyword() string {
	if len(s.content) < 2 {
		return ""
	}
	if s.content[0] != ':' {
		return ""
	}
	if !unicode.IsLetter(s.content[1]) {
		return ""
	}
	offset := 2
	content := s.content[offset:]
	for _, r := range content {
		if r == ' ' {
			break
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return ""
		}					
		offset++
	}
	s.column += offset
	result := s.content[:offset]
	s.content = s.content[offset:]
	return string(result)
}

func (s *script) EatInteger() string {
	var digits []rune = s.integer.Scan(s.content)
	if len(digits) == 0 {
		return ""
	}
	s.content = s.content[len(digits):]
	s.column+= len(digits)
	return string(digits)
}

func (s *script) EatFloat() string {
	var digits []rune = s.float.Scan(s.content)
	if len(digits) == 0 {
		return ""
	}
	s.content = s.content[len(digits):]
	s.column+= len(digits)
	return string(digits)
}
