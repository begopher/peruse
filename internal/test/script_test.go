package test

import(
	"testing"
	"github.com/begopher/peruse"
)

func TestRemain(t *testing.T) {
	content := `buffer content`
	buffer := peruse.Script("", content)
	if buffer.Remain() != content {
		t.Errorf("Buffer does not have the right content")
	}

}

func TestLength(t *testing.T) {
	table := []struct {
		content string
		length int
	}{
		{"", 0},
		{"a", 1},
		{"ab", 2},
		{"any", 3},
	}

	for _, data := range table {
		buffer := peruse.Script("", data.content)
		expected := data.length
		got := buffer.Length()
		if got != expected {
			t.Errorf("Buffer does not return correct length, got(%d) expected (%d)", got, expected)
		}
	}
}

func TestPath(t *testing.T) {
	table := []string {
		"main.twq",
			"student.twq",
			"teacher.twq",
			"any.twq",
		}

	for _, path := range table {
		text := peruse.Script(path, "any")
		expected := path
		got := text.Origin()
		if got != expected {
			t.Errorf("Script does not return correct origin, got(%s) expected (%s)", got, expected)
		}
	}
}

func TestEat(t *testing.T) {
	table := []struct {
		content string
		expected bool
		eat string
		remain string
		column int
		line int
	}{
		{
			content: "universe\n programming language",
			expected: true,
			eat: "universe",
			remain: "\n programming language",			
			column: 9,
			line: 1,
		},
		{
			content: "universe\n programming language",
			expected: true,
			eat: "universe\n p",
			remain: "rogramming language",
			column: 3,
			line: 2,
		},
		{
			content: "universe\n programming language",
			expected: false,
			eat: "uuuniverse\n p",
			remain: "universe\n programming language",
			column: 1,
			line: 1,
		},
	}
	for _, data := range table {
		buffer := peruse.Script("", data.content)
		

		if got, expected := buffer.Eat(data.eat), data.expected; got != expected {
			t.Errorf("Buffer does not eat prefix correctly, expected (%v) got (%v) ",expected, got)
		}

		if got, expected := buffer.Remain(), data.remain; got != expected {
			t.Errorf("Buffer does not eat prefix correctly, remain: got(%s) expected (%s)", got, expected)
		}

		if got, expected := buffer.Column(), data.column; got != expected {
			t.Errorf("Buffer.Eat(prefix) dose not count column currectly, got(%d) expected (%d)", got, expected)
		}

		if got, expected := buffer.Line(), data.line; got != expected {
			t.Errorf("Buffer.Eat(prefix) dose not count line currectly, got(%d) expected (%d)", got, expected)
		}
	}
}

func TestString(t *testing.T) {
	table := []struct {
		content string
		expected string
		remain string
		column int
		line int 
	}{
		{content: ``, expected: "", remain: "", column: 1 , line: 1},
		{content: `""any2`, expected: "", remain: "any2", column: 3, line: 1},
		{
			content: `"any\"value3"any3`,
			expected: `any\"value3`,
			remain: "any3",
			column: 14,
			line: 1,
		},
		{
			content: `"any
 Value4\""any4`,
			expected: `any
 Value4\"`,
			remain: "any4",
			column: 10,
			line: 2,
		},
		{
			content: "\"any \n Value5\"any5",
			expected: "any \n Value5",
			remain: "any5",
			column: 8,
			line: 2,
		},
		
	}

	for _, data := range table {
		buffer := peruse.Script("", data.content)
		expected := data.expected
		got, _ := buffer.EatString()
		if got != expected {
			t.Errorf("Buffer does not eat string correctly, got(%s) expected (%s)", got, expected)
		}

		remain := buffer.Remain()
		if remain != data.remain {
			t.Errorf("Buffer does not eat string correctly, remain: got(%s) expected (%s)", remain, data.remain)
		}

		if got, expected := buffer.Column(), data.column; got != expected {
			t.Errorf("Buffer.String dose not count column currectly, got(%d) expected (%d)", got, expected)
		}

		if got, expected := buffer.Line(), data.line; got != expected {
			t.Errorf("Buffer.String dose not count line currectly, got(%d) expected (%d)", got, expected)
		}
	}
}

func TestEatSpaces(t *testing.T) {
	table := []struct {
		content string
		remain string
		column int
		line int
	}{
		{content: " any1", remain: "any1", column: 2, line: 1},
		{content: "    any2", remain: "any2", column: 5, line: 1},
		{content: "\tany3", remain: "any3", column: 2, line: 1},
		{content: " \t any4", remain: "any4", column: 4, line: 1},
		{content: "\nany5", remain: "any5", column: 1, line: 2},
		{content: "\n any6", remain: "any6", column: 2, line: 2},
		{content: "\n\n\n\t\t\n    any5", remain: "any5", column: 5, line: 5},
 	}

	for _, data := range table {
		buffer := peruse.Script("any", data.content)
		buffer.EatSpaces()
		expected := data.remain
		got := buffer.Remain()
		if got != expected {
			t.Errorf("Buffer does not eat spaces correctly, got(%s) expected (%s)", got, expected)
		}

		if got, expected := buffer.Column(), data.column; got != expected {
			t.Errorf("Buffer does not eat spaces correctly, column: got(%d) expected (%d)", got, expected)
		}

		if got, expected := buffer.Line(), data.line; got != expected {
			t.Errorf("Buffer does not eat spaces correctly, line: got(%d) expected (%d)", got, expected)
		}
	}
}

func TestEatWord(t *testing.T) {
	table := []struct {
		content string
		expected string
		remain string
		col int
		line int 
	}{
		{"A", "A", "", 2, 1},
		{"Any more1", "Any", " more1", 4, 1},
		{"AnyAny more1", "AnyAny", " more1", 7, 1},
		{"A-ny more1", "", "A-ny more1", 1, 1},
		{"A_ny more1", "", "A_ny more1", 1, 1},
		{"any more2", "any", " more2", 4, 1},
		{"any1 more3", "any1", " more3", 5, 1},
		{"1any", "", "1any", 1, 1},
		{" Any more5", "", " Any more5", 1, 1},
		{":any", "", ":any", 1, 1},
		{"_any", "", "_any", 1, 1},
		{"111234", "", "111234", 1 ,1},
	}
	for _, data := range table {
		buffer := peruse.Script("", data.content)
		got := buffer.EatWord()
		expected := data.expected		
		if got != expected {
			t.Errorf("Buffer does not eat word correctly, got(%s) expected (%s)", got, expected)
		}
		expected = data.remain
		got = buffer.Remain()
		if got != expected {
			t.Errorf("Buffer does not eat word correctly, remain: got(%s) expected (%s)", got, expected)
		}

		if expected, got := data.col, buffer.Column(); expected != got {
			t.Errorf("Buffer does not eat word correctly, column: got(%d) expected (%d)", got, expected)
		}
		if expected, got := data.line, buffer.Line(); expected != got {
			t.Errorf("Buffer does not eat word correctly, line: got(%d) expected (%d)", got, expected)
		}
	}
}

func TestIsWord(t *testing.T) {
	table := []struct {
		content string
		answer bool
	}{
		{content: "a", answer: true},
		{content: "A", answer: true},
		{content: "abc", answer: true},
		{content: "ABC", answer: true},
		{content: "a1", answer: true},
		{content: "B1", answer: true},
		{content: "A1b2C3d4", answer: true},
		{content: "0", answer: false},
		{content: "1", answer: false},
		{content: "3456789", answer: false},
		{content: "4A", answer: false},
		{content: "A-", answer: false},
		{content: "A-A", answer: false},
		{content: "A-3", answer: false},
		
	}
	for _, data := range table {
		buffer := peruse.Script("", data.content)
		if got, expected := buffer.IsWord(data.content), data.answer; got != expected {
			t.Errorf(`IsWord("%v") returns (%v) expected (%v)`, data.content, got, expected)
		}		
	}
}

func TestEatPrefixedWord(t *testing.T) {
	table := []struct {
		content string
		prefix  string
		word string
		prefixedWord string
		remain string
		col int
		line int 
	}{
		{content: "a", prefix: "", word: "a", prefixedWord: "a", remain:"", col: 2, line: 1},
		{content: "A", prefix: "", word: "A", prefixedWord: "A", remain:"", col: 2, line: 1},
		{content: "1234", prefix: ":", word: "", prefixedWord: "", remain:"1234", col: 1, line: 1},
		{content: ":1234", prefix: "@", word: "", prefixedWord: "", remain:":1234", col: 1, line: 1},
		{content: ":A)",  prefix: ":",  word: "A", prefixedWord: ":A", remain:")", col: 3, line: 1},
		{content: `:A
`,  prefix: ":",  word: "A", prefixedWord: ":A", remain:"\n", col: 3, line: 1},
		{content: "#A",  prefix: "#",  word: "A", prefixedWord: "#A", remain:"", col: 3, line: 1},
		{content: "@name",  prefix: "@",  word: "name", prefixedWord: "@name", remain:"", col: 6, line: 1},
		{content: "$any word", prefix: "$", word: "any", prefixedWord: "$any", remain:" word", col: 5, line: 1},
		{content: "__why lisp", prefix: "__", word: "why", prefixedWord: "__why", remain:" lisp", col: 6, line: 1},
		{content: "__why-lisp", prefix: "__", word: "", prefixedWord: "", remain:"__why-lisp", col: 1, line: 1},
		{content: "#:why lisp", prefix: "#:", word: "why", prefixedWord: "#:why", remain:" lisp", col: 6, line: 1},
		
	}
	for _, data := range table {
		buffer := peruse.Script("any", data.content)
		word, prefixedWord := buffer.EatPrefixedWord(data.prefix)
		expected := data.word
		if got, expected := word, data.word;  got != expected {
			t.Errorf("EatPrefixedWord(%v), expected word is (%s) got (%s)", data.content, expected, got)
		}
		if got, expected := prefixedWord, data.prefixedWord;  got != expected {
			t.Errorf("EatPrefixedWord(%v), expected prefixedWord is (%s) got (%s)", data.content, expected, got)
		}
		expected = data.remain
		got := buffer.Remain()
		if got != expected {
			t.Errorf("Buffer does not eat prefixed word correctly, remain: got(%s) expected (%s)", got, expected)
		}

		if expected, got := data.col, buffer.Column(); expected != got {
			t.Errorf("Buffer does not eat prefixed word correctly, column: got(%d) expected (%d)", got, expected)
		}
		if expected, got := data.line, buffer.Line(); expected != got {
			t.Errorf("Buffer does not eat prefixed word correctly, line: got(%d) expected (%d)", got, expected)
		}
	}

}

func TestEatWords(t *testing.T) {
	table := []struct {
		content string
		first   string
		second  string
		remain   string
		col int
		line int 
	}{
		{content: "", first: "", second: "", remain: "", col:1, line:1},
		{content: "A", first: "", second: "", remain: "A", col:1, line:1},		
		{content: "A-1:B2", first: "", second: "", remain: "A-1:B2", col:1, line:1},
		{content: "A1:B-2", first: "", second: "", remain: "A1:B-2", col:1, line:1},
		{content: "A1:B2 ", first: "A1", second: "B2", remain: " ", col:6, line:1},
		{content: "A1:B2)", first: "A1", second: "B2", remain: ")", col:6, line:1},
		{content: `A1:B2
`, first: "A1", second: "B2", remain: "\n", col:6, line:1},
		
		{content: "lispFirst:lispLast", first: "lispFirst", second: "lispLast", remain: "", col:19, line:1},		
	}
	
	for _, data := range table {
		buffer := peruse.Script("", data.content)
		first, second := buffer.EatWords()
		if first != data.first {
			t.Errorf("Buffer does not eat words correctly, expected first is (%s) got (%s)", data.first, first)
		}
		if second != data.second {
			t.Errorf("Buffer does not eat words correctly, expected second is (%s) got (%s)", data.second, second)
		}
		expected := data.remain
		got := buffer.Remain()
		if got != expected {
			t.Errorf("Buffer does not eat words correctly, remain: got(%s) expected (%s)", got, expected)
		}

		if expected, got := data.col, buffer.Column(); expected != got {
			t.Errorf("Buffer does not eat word correctly, column: got(%d) expected (%d)", got, expected)
		}
		if expected, got := data.line, buffer.Line(); expected != got {
			t.Errorf("Buffer does not eat word correctly, line: got(%d) expected (%d)", got, expected)
		}
	}
}


func TestEatPrefixedSymbol(t *testing.T) {
	table := []struct {
		content string
		prefix  string
		symbol string
		prefixedSymbol string
		remain string
		col int
		line int 
	}{

		{content: "", prefix: "", symbol: "", prefixedSymbol: "", remain:"", col: 1, line: 1},
		{content: "a", prefix: "", symbol: "a", prefixedSymbol: "a", remain:"", col: 2, line: 1},
		{content: "A", prefix: "", symbol: "A", prefixedSymbol: "A", remain:"", col: 2, line: 1},
		{content: "a-a", prefix: "", symbol: "a-a", prefixedSymbol: "a-a", remain:"", col: 4, line: 1},
		{content: "1a-a", prefix: "", symbol: "", prefixedSymbol: "", remain:"1a-a", col: 1, line: 1},
		{content: "a0-1a2", prefix: "", symbol: "a0-1a2", prefixedSymbol: "a0-1a2", remain:"", col: 7, line: 1},

		{content: ":A",  prefix: ":",  symbol: "A", prefixedSymbol: ":A", remain:"", col: 3, line: 1},
		{content: "#A",  prefix: "#",  symbol: "A", prefixedSymbol: "#A", remain:"", col: 3, line: 1},
		{content: "@A",  prefix: "@",  symbol: "A", prefixedSymbol: "@A", remain:"", col: 3, line: 1},
		
		{content: ":1a", prefix: ":", symbol: "", prefixedSymbol: "", remain:":1a", col: 1, line: 1},
		{content: ":a1", prefix: ":", symbol: "a1", prefixedSymbol: ":a1", remain:"", col: 4, line: 1},
		{content: ":a1 any", prefix: ":", symbol: "a1", prefixedSymbol: ":a1", remain:" any", col: 4, line: 1},
		{content: ":a1- any", prefix: ":", symbol: "",  prefixedSymbol: "", remain:":a1- any", col: 1, line: 1},
		{content: ":a1-any", prefix: ":", symbol: "a1-any",  prefixedSymbol: ":a1-any", remain:"", col: 8, line: 1},
		{content: ":a1-any lisp", prefix: ":", symbol: "a1-any", prefixedSymbol: ":a1-any", remain:" lisp", col: 8, line: 1},
		{content: ":a1-2any3 lisp", prefix: ":", symbol: "a1-2any3", prefixedSymbol: ":a1-2any3", remain:" lisp", col: 10, line: 1},
		
		{content: "__why-lisp", prefix: "__", symbol: "why-lisp", prefixedSymbol: "__why-lisp", remain:"", col: 11, line: 1},
		{content: "__why-lisp25", prefix: "__", symbol: "why-lisp25", prefixedSymbol: "__why-lisp25", remain:"", col: 13, line: 1},
		{content: "#:why-lisp )", prefix: "#:", symbol: "why-lisp", prefixedSymbol: "#:why-lisp", remain:" )", col: 11, line: 1},

		

		{content: "1234", prefix: ":", symbol: "", prefixedSymbol: "", remain:"1234", col: 1, line: 1},		
		{content: ":1234", prefix: "@", symbol: "", prefixedSymbol: "", remain:":1234", col: 1, line: 1},
		{content: ":A)",  prefix: ":",  symbol: "A", prefixedSymbol: ":A", remain:")", col: 3, line: 1},
		{content: ":A-B)",  prefix: ":",  symbol: "A-B", prefixedSymbol: ":A-B", remain:")", col: 5, line: 1},
		{content: `:A
`,  prefix: ":",  symbol: "A", prefixedSymbol: ":A", remain:"\n", col: 3, line: 1},
		{content: `:A-B
`,  prefix: ":", symbol: "A-B",  prefixedSymbol: ":A-B", remain:"\n", col: 5, line: 1},
		
	}
	for _, data := range table {
		buffer := peruse.Script("any", data.content)
		symbol, prefixedSymbol := buffer.EatPrefixedSymbol(data.prefix)
		if got, expected := symbol, data.symbol; got != expected {
			t.Errorf("EatPrefixedSymbol(%v) returned symbol is (%s) expected (%s)", data.content, got, expected)
		}
		if got, expected := prefixedSymbol, data.prefixedSymbol; got != expected {
			t.Errorf("EatPrefixedSymbol(%v) returned prefixedSymbol is (%s) expected (%s)", data.content, got, expected)
		}
		expected := data.remain
		got := buffer.Remain()
		if got != expected {
			t.Errorf("Buffer does not eat prefixed symbol correctly, remain: got(%s) expected (%s)", got, expected)
		}

		if expected, got := data.col, buffer.Column(); expected != got {
			t.Errorf("Buffer does not eat prefixed symbol correctly, column: got(%d) expected (%d)", got, expected)
		}
		if expected, got := data.line, buffer.Line(); expected != got {
			t.Errorf("Buffer does not eat prefixed symbol correctly, line: got(%d) expected (%d)", got, expected)
		}
	}

}

func TestEatSymbol(t *testing.T) {
	table := []struct {
		content string
		expected string
		remain string
		col int
		line int 
	}{
		{"A", "A", "", 2, 1},
		{"Any more1", "Any", " more1", 4, 1},
		{"AnyAny more1", "AnyAny", " more1", 7, 1},
		{"A-ny more1", "A-ny", " more1", 5, 1},
		{"A-ny- more1", "", "A-ny- more1", 1, 1},
		{"A_ny more1", "", "A_ny more1", 1, 1},
		{"any1 more3", "any1", " more3", 5, 1},
		{"any-1 more3", "any-1", " more3", 6, 1},
		{"1any", "", "1any", 1, 1},
		{" Any more5", "", " Any more5", 1, 1},
		{":any", "", ":any", 1, 1},
		{"_any", "", "_any", 1, 1},
		{"111234", "", "111234", 1 ,1},
	}
	for _, data := range table {
		buffer := peruse.Script("", data.content)
		got := buffer.EatSymbol()
		expected := data.expected		
		if got != expected {
			t.Errorf("Buffer does not eat symbol correctly, got(%s) expected (%s)", got, expected)
		}
		expected = data.remain
		got = buffer.Remain()
		if got != expected {
			t.Errorf("Buffer does not eat symbol correctly, remain: got(%s) expected (%s)", got, expected)
		}

		if expected, got := data.col, buffer.Column(); expected != got {
			t.Errorf("Buffer does not eat symbol correctly, column: got(%d) expected (%d)", got, expected)
		}
		if expected, got := data.line, buffer.Line(); expected != got {
			t.Errorf("Buffer does not eat symbol correctly, line: got(%d) expected (%d)", got, expected)
		}
	}
}

func TestIsSymbol(t *testing.T) {
	table := []struct {
		content string
		answer bool
	}{
		{content: "a", answer: true},
		{content: "A", answer: true},
		{content: "abc", answer: true},
		{content: "ABC", answer: true},
		{content: "a1", answer: true},
		{content: "B1", answer: true},
		{content: "A1b2C3d4", answer: true},
		{content: "0", answer: false},
		{content: "1", answer: false},
		{content: "3456789", answer: false},
		{content: "4A", answer: false},
		{content: "A-", answer: false},
		//
		{content: "A-A", answer: true},
		{content: "A-A-", answer: false},
		{content: "A-A-1234", answer: true},
		{content: "A-4", answer: true},
		
	}
	for _, data := range table {
		buffer := peruse.Script("", data.content)
		if got, expected := buffer.IsSymbol(data.content), data.answer; got != expected {
			t.Errorf(`IsSymbol("%v") returns (%v) expected (%v)`, data.content, got, expected)
		}		
	}
}

func TestEatSymbols(t *testing.T) {
	table := []struct {
		content string
		first   string
		second  string
		remain   string
		col int
		line int 
	}{
		{content: "", first: "", second: "", remain: "", col:1, line:1},
		{content: "A", first: "", second: "", remain: "A", col:1, line:1},		
		{content: "A1:B-2", first: "A1", second: "B-2", remain: "", col:7, line:1},
		{content: "A1:B-2 ", first: "A1", second: "B-2", remain: " ", col:7, line:1},
		{content: "A1:B2)", first: "A1", second: "B2", remain: ")", col:6, line:1},
		{content: `A1:B2
`, first: "A1", second: "B2", remain: "\n", col:6, line:1},
		
		{content: "lisp-first:lisp-last", first: "lisp-first", second: "lisp-last", remain: "", col:21, line:1},
		{
			content: "lisp0-1first3:lisp4-5last6",
			first:   "lisp0-1first3",
			second:  "lisp4-5last6",
			remain:  "",
			col:     27,
			line:    1,
		},
	}
	
	for _, data := range table {
		buffer := peruse.Script("", data.content)
		first, second := buffer.EatSymbols()
		if first != data.first {
			t.Errorf("Buffer does not eat symbols correctly, expected first is (%s) got (%s)", data.first, first)
		}
		if second != data.second {
			t.Errorf("Buffer does not eat symbols correctly, expected second is (%s) got (%s)", data.second, second)
		}
		expected := data.remain
		got := buffer.Remain()
		if got != expected {
			t.Errorf("Buffer does not eat symbols correctly, remain: got(%s) expected (%s)", got, expected)
		}

		if expected, got := data.col, buffer.Column(); expected != got {
			t.Errorf("Buffer does not eat symbol correctly, column: got(%d) expected (%d)", got, expected)
		}
		if expected, got := data.line, buffer.Line(); expected != got {
			t.Errorf("Buffer does not eat symbol correctly, line: got(%d) expected (%d)", got, expected)
		}
	}
}

func TestKeyword(t *testing.T) {
	table := []struct {
		content string
		expected string
		remain string
		col int
		line int 
	}{
		{content: "", expected: "", remain:"", col: 1, line: 1},
		{content: "A", expected: "", remain:"A", col: 1, line: 1},
		{content: "1234", expected: "", remain:"1234", col: 1, line: 1},
		{content: ":1234", expected: "", remain:":1234", col: 1, line: 1},
		{content: ":A", expected: ":A", remain:"", col: 3, line: 1},
		{content: ":any n", expected: ":any", remain:" n", col: 5, line: 1},
		{content: ":a-ny", expected: "", remain:":a-ny", col: 1, line: 1},
		{content: ":a_ny", expected: "", remain:":a_ny", col: 1, line: 1},
		{content: ":-any", expected: "", remain:":-any", col: 1, line: 1},
		{content: ":_any", expected: "", remain:":_any", col: 1, line: 1},
	}
	for _, data := range table {
		buffer := peruse.Script("any", data.content)
		got := buffer.EatKeyword()
		expected := data.expected		
		if got != expected {
			t.Errorf("Buffer does not eat keyword correctly, got(%s) expected (%s)", got, expected)
		}
		expected = data.remain
		got = buffer.Remain()
		if got != expected {
			t.Errorf("Buffer does not eat keyword correctly, remain: got(%s) expected (%s)", got, expected)
		}

		if expected, got := data.col, buffer.Column(); expected != got {
			t.Errorf("Buffer does not eat word correctly, column: got(%d) expected (%d)", got, expected)
		}
		if expected, got := data.line, buffer.Line(); expected != got {
			t.Errorf("Buffer does not eat word correctly, line: got(%d) expected (%d)", got, expected)
		}
	}

}


func TestInteger(t *testing.T) {
	table := []struct {
		content string
		expected string
		remain string
		col int
		line int 
	}{
		{content: "", expected: "", remain:"", col: 1, line: 1},
		{content: "A", expected: "", remain:"A", col: 1, line: 1},
		// signed
		{content: "+0", expected: "+0", remain:"", col: 3, line: 1},
		{content: "-0", expected: "-0", remain:"", col: 3, line: 1},
		{content: "+1", expected: "+1", remain:"", col: 3, line: 1},
		{content: "-1", expected: "-1", remain:"", col: 3, line: 1},
		{content: "+01234", expected: "+01234", remain:"", col: 7, line: 1},
		{content: "-12345", expected: "-12345", remain:"", col: 7, line: 1},
		{content: "+01234 any", expected: "+01234", remain:" any", col: 7, line: 1},
		{content: "-12345 any", expected: "-12345", remain:" any", col: 7, line: 1},
		{content: "+01234-any", expected: "", remain:"+01234-any", col: 1, line: 1},
		{content: "-12345-123", expected: "", remain:"-12345-123", col: 1, line: 1},
		{content: "+01234-123", expected: "", remain:"+01234-123", col: 1, line: 1},
		{content: "+01234_123", expected: "", remain:"+01234_123", col: 1, line: 1},
		{content: "-12345_123", expected: "", remain:"-12345_123", col: 1, line: 1},
		// unsigned
		{content: "0", expected: "0", remain:"", col: 2, line: 1},
		{content: "1", expected: "1", remain:"", col: 2, line: 1},
		{content: "01234", expected: "01234", remain:"", col: 6, line: 1},
		{content: "01234any", expected: "", remain:"01234any", col: 1, line: 1},
		{content: "01234 any", expected: "01234", remain:" any", col: 6, line: 1},
		{content: "12345-123", expected: "", remain:"12345-123", col: 1, line: 1},
		{content: "12345_123", expected: "", remain:"12345_123", col: 1, line: 1},

	}
	for _, data := range table {
		buffer := peruse.Script("any", data.content)
		got := buffer.EatInteger()
		expected := data.expected		
		if got != expected {
			t.Errorf("Buffer does not eat integer correctly, got(%s) expected (%s)", got, expected)
		}
		expected = data.remain
		got = buffer.Remain()
		if got != expected {
			t.Errorf("Buffer does not eat integer correctly, remain: got(%s) expected (%s)", got, expected)
		}
		if expected, got := data.col, buffer.Column(); expected != got {
			t.Errorf("Buffer does not eat integer correctly, column: got(%d) expected (%d)", got, expected)
		}
		if expected, got := data.line, buffer.Line(); expected != got {
			t.Errorf("Buffer does not eat integer correctly, line: got(%d) expected (%d)", got, expected)
		}
	}

}


func TestFloat(t *testing.T) {
	table := []struct {
		content string
		expected string
		remain string
		col int
		line int 
	}{
		// invalid floats
		{content: "", expected: "", remain:"", col: 1, line: 1},
		{content: "A", expected: "", remain:"A", col: 1, line: 1},
		{content: ".", expected: "", remain:".", col: 1, line: 1},
		{content: "-.", expected: "", remain:"-.", col: 1, line: 1},
		{content: "+.", expected: "", remain:"+.", col: 1, line: 1},
		{content: "+1..1", expected: "", remain:"+1..1", col: 1, line: 1},
		{content: "-1.1.", expected: "", remain:"-1.1.", col: 1, line: 1},
		// underscore is not allowed 
		{content: "-1.1000_000", expected: "", remain:"-1.1000_000", col: 1, line: 1},
		{content: "+2.1000_000", expected: "", remain:"+2.1000_000", col: 1, line: 1},
		{content: "-3000_000.00", expected: "", remain:"-3000_000.00", col: 1, line: 1},
		{content: "+3000_000.00", expected: "", remain:"+3000_000.00", col: 1, line: 1},
		{content: "1.1000_000", expected: "", remain:"1.1000_000", col: 1, line: 1},
		{content: "2.1000_000", expected: "", remain:"2.1000_000", col: 1, line: 1},
		{content: "3000_000.00", expected: "", remain:"3000_000.00", col: 1, line: 1},
		{content: "3000_000.00", expected: "", remain:"3000_000.00", col: 1, line: 1},
		

		// integer
		{content: "+0", expected: "", remain:"+0", col: 1, line: 1},
		{content: "-0", expected: "", remain:"-0", col: 1, line: 1},
		{content: "+1", expected: "", remain:"+1", col: 1, line: 1},
		{content: "-1", expected: "", remain:"-1", col: 1, line: 1},
		
		// signed with a digit
		{content: "+.0", expected: "+.0", remain:"", col: 4, line: 1},
		{content: "-.0", expected: "-.0", remain:"", col: 4, line: 1},
		{content: "+.5", expected: "+.5", remain:"", col: 4, line: 1},
		{content: "-.9", expected: "-.9", remain:"", col: 4, line: 1},
		{content: "+0.", expected: "+0.", remain:"", col: 4, line: 1},
		{content: "-0.", expected: "-0.", remain:"", col: 4, line: 1},
		{content: "+5.", expected: "+5.", remain:"", col: 4, line: 1},
		{content: "-9.", expected: "-9.", remain:"", col: 4, line: 1},				
		// signed with a digits
		{content: "+0.1234", expected: "+0.1234", remain:"", col: 8, line: 1},
		{content: "+1.1234 any0", expected: "+1.1234", remain:" any0", col: 8, line: 1},
		{content: "-1234.5", expected: "-1234.5", remain:"", col: 8, line: 1},
		{content: "-1234.5 any5", expected: "-1234.5", remain:" any5", col: 8, line: 1},
		{content: "+88.555 any6", expected: "+88.555", remain:" any6", col: 8, line: 1},
		{content: "-77.45 any7", expected: "-77.45", remain:" any7", col: 7, line: 1},		
		// unsigned with a digit
		{content: ".0", expected: ".0", remain:"", col: 3, line: 1},
		{content: "0.", expected: "0.", remain:"", col: 3, line: 1},
		{content: ".5", expected: ".5", remain:"", col: 3, line: 1},
		{content: ".9", expected: ".9", remain:"", col: 3, line: 1},
		{content: "7. any", expected: "7.", remain:" any", col: 3, line: 1},
		// unsigned with digits
		{content: "2.4321", expected: "2.4321", remain:"", col: 7, line: 1},
		{content: "3.4321 any1", expected: "3.4321", remain:" any1", col: 7, line: 1},
		{content: "41.909 any2", expected: "41.909", remain:" any2", col: 7, line: 1},
		{content: "12.21 any3", expected: "12.21", remain:" any3", col: 6, line: 1},		

	}
	for _, data := range table {
		path := "float.twq"
		script := peruse.Script(path, data.content)
		got := script.EatFloat()
		expected := data.expected		
		if got != expected {
			t.Errorf("Buffer does not eat float correctly, got(%s) expected (%s)", got, expected)
		}
		expected = data.remain
		got = script.Remain()
		if got != expected {
			t.Errorf("Buffer does not eat float correctly, remain: got(%s) expected (%s)", got, expected)
		}
		if expected, got := data.col, script.Column(); expected != got {
			t.Errorf("Buffer does not eat float correctly, column: got(%d) expected (%d)", got, expected)
		}
		if expected, got := data.line, script.Line(); expected != got {
			t.Errorf("Buffer does not eat float correctly, line: got(%d) expected (%d)", got, expected)
		}
		// addition
		location := script.Location()
		if expected, got := data.col, location.Column(); expected != got {
			t.Errorf("Script does not eat float correctly, column: got(%d) expected (%d)", got, expected)
		}
		if expected, got := data.line, location.Line(); expected != got {
			t.Errorf("Script does not eat float correctly, line: got(%d) expected (%d)", got, expected)
		}		

		if expected, got := path, location.Origin(); got != expected {
			t.Errorf("Script return invalid origin, got(%s) expected (%s)", got, expected)
		}
		
	}

}


func TestBeginWith(t *testing.T) {
	table := []struct {
		content string
		prefix string
		expected bool
	}{
		{content: "", prefix: "", expected: false},
		{content: "any text", prefix: "a", expected: true},
		{content: "any text", prefix: "an", expected: true},
		{content: "any text", prefix: "any", expected: true},
		{content: "any text", prefix: "any ", expected: true},
		{content: "any text", prefix: "any t", expected: true},
		{content: "any text", prefix: "any te", expected: true},
		{content: "any text", prefix: "any tex", expected: true},
		{content: "any text", prefix: "any text", expected: true},
		{content: "any text", prefix: "any text.", expected: false},
		{content: "any text", prefix: "text any", expected: false},
		{content: "any text", prefix: "ext any", expected: false},
		{content: "any text", prefix: "xt any", expected: false},
		{content: "any text", prefix: "t any", expected: false},
		{content: "any text", prefix: " any", expected: false},

	}
	for _, data := range table {
		buffer := peruse.Script("any", data.content)		
		if expected, got := data.expected, buffer.BeginWith(data.prefix); expected != got {
			t.Errorf(`Buffer ("%s").BeginWith("%s"), expected (%v) got(%v)`, data.content, data.prefix, expected, got )
		}

	}

}
