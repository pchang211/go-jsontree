// Copyright 2013, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// lexer.go [created: Sun, 23 Jun 2013]

package jsonpath

import (
	"fmt"
	"os"
	"unicode"

	"github.com/bmatsuo/go-lexer"
)

// EOF is the constant for EOF
const EOF = lexer.EOF

const (
	ItemEOF                    = lexer.ItemEOF
	ItemError                  = lexer.ItemError
	ItemPathKey lexer.ItemType = iota
	ItemDollar
	ItemDot
	ItemDotDot
	ItemStar
	ItemStarStar
	ItemLeftBracket
	ItemRightBracket
	ItemNumber
	ItemNull
	ItemTrue
	ItemFalse
	ItemEqual
	ItemGreater
	ItemGreaterEqual
	ItemLess
	ItemLessEqual
	ItemNotEqual
	ItemString
)

// JSONPathLexer is a wrapper around go-lexer lexer object
type JSONPathLexer struct {
	// Next() *lexer.Item
	lex *lexer.Lexer
}

// NewJSONPathLexer returns a new JSONPathLexer object
func NewJSONPathLexer(input string) *JSONPathLexer {
	return &JSONPathLexer{lex: lexer.New(Start, input)}
}

// Next calls next on the lexer object
func (l *JSONPathLexer) Next() *lexer.Item {
	return l.lex.Next()
}

// Start state for lexer
func Start(lex *lexer.Lexer) lexer.StateFn {
	// get rid of white space
	if lex.AcceptRunRange(unicode.Space) > 0 {
		lex.Ignore()
		return Start
	}

	// handle DOTs
	switch lex.AcceptRun(".") {
	case 0:
		break
	case 1:
		debugln("FOUND DOT")
		lex.Emit(ItemDot)
		return Start
	case 2:
		debugln("FOUND DOT DOT")
		lex.Emit(ItemDotDot)
		return Start
	default:
		return lex.Errorf("unexpected '.'")
	}

	// wildcard
	switch lex.AcceptRun("*") {
	case 0:
		break
	case 1:
		debugln("FOUND STAR")
		lex.Emit(ItemStar)
	case 2:
		debugln("FOUND STAR STAR")
		lex.Emit(ItemStarStar)
		return Start
	default:
		return lex.Errorf("unexpected '*'")
	}

	// everything else
	switch r, _ := lex.Peek(); {
	case r == lexer.EOF:
		return nil
	case unicode.IsLetter(r):
		return PathKey
	case unicode.IsDigit(r):
		return Number
	case r == '[':
		debugln("FOUND LEFT BRACKET")
		lex.Advance()
		lex.Emit(ItemLeftBracket)
		return Bracket
	case r == ']':
		debugln("FOUND RIGHT BRACKET")
		lex.Advance()
		lex.Emit(ItemRightBracket)
	case r == '$':
		debugln("FOUND DOLLAR")
		lex.Advance()
		lex.Emit(ItemDollar)
	case r == '>':
		lex.Advance()
		if r, _ := lex.Peek(); r == '=' {
			debugln("FOUND GREATER EQUAL")
			lex.Advance()
			lex.Emit(ItemGreaterEqual)
		} else {
			debugln("FOUND GREATER")
			lex.Emit(ItemGreater)
		}
	case r == '<':
		lex.Advance()
		if r, _ := lex.Peek(); r == '=' {
			debugln("FOUND LESS EQUAL")
			lex.Advance()
			lex.Emit(ItemLessEqual)
		} else {
			debugln("FOUND LESS")
			lex.Emit(ItemLess)
		}
	case r == '=':
		debugln("FOUND EQUAL")
		lex.Advance()
		lex.Emit(ItemEqual)
	case r == '!':
		lex.Advance()
		if r, _ := lex.Peek(); r == '=' {
			debugln("FOUND NOT EQUAL")
			lex.Advance()
			lex.Emit(ItemNotEqual)
		} else {
			return lex.Errorf("expected '=' got %c", r)
		}
	}

	// default
	return Start
}

func StartBracket(lex *lexer.Lexer) lexer.StateFn {
	if lex.Accept("[") {
		lex.Emit(ItemLeftBracket)
		return Bracket
	} else {
		lex.Errorf("missing '['")
		return nil
	}
}

func PathKey(lex *lexer.Lexer) lexer.StateFn {
	found := false
	for {
		if lex.AcceptRun("_") > 0 {
			found = true
			continue
		}
		if lex.AcceptRunRange(unicode.Letter) > 0 {
			found = true
			continue
		}
		break
	}
	if found {
		debugln("FOUND PATH KEY")
		lex.Emit(ItemPathKey)
	}
	return Start
}

func Number(lex *lexer.Lexer) lexer.StateFn {
	if lex.AcceptRunRange(unicode.Digit) == 0 {
		r, _ := lex.Peek()
		if r == lexer.EOF {
			return nil
		}
		return lex.Errorf("expected digit got %c", r)
	}
	debugln("FOUND NUMBER")
	if lex.Accept(".") {
		lex.AcceptRunRange(unicode.Digit)
	}
	if lex.Accept("eE") {
		lex.Accept("-")
		lex.AcceptRunRange(unicode.Digit)
	}
	lex.Emit(ItemNumber)
	return Start
}

// Bracket state
func Bracket(lex *lexer.Lexer) lexer.StateFn {
	switch r, _ := lex.Peek(); {
	case r == lexer.EOF:
		return nil
	case unicode.IsDigit(r):
		return BracketNumber
	case unicode.IsLetter(r):
		return BracketKey
	case r == '>':
		lex.Advance()
		if r, _ := lex.Peek(); r == '=' {
			lex.Advance()
			lex.Emit(ItemGreaterEqual)
		} else {
			lex.Emit(ItemGreater)
		}
	case r == '<':
		lex.Advance()
		if r, _ := lex.Peek(); r == '=' {
			lex.Advance()
			lex.Emit(ItemLessEqual)
		} else {
			lex.Emit(ItemLess)
		}
	case r == '=':
		lex.Advance()
		lex.Emit(ItemEqual)
	case r == '!':
		lex.Advance()
		if r, _ := lex.Peek(); r == '=' {
			lex.Advance()
			lex.Emit(ItemNotEqual)
		} else {
			lex.Errorf("expected '=' got %c", r)
			return nil
		}
	}
	return Bracket
}

func BracketNumber(lex *lexer.Lexer) lexer.StateFn {
	if lex.AcceptRunRange(unicode.Digit) == 0 {
		r, _ := lex.Peek()
		if r == lexer.EOF {
			return nil
		}
		lex.Errorf("expected digit got %c", r)
	}
	if lex.Accept(".") {
		lex.AcceptRunRange(unicode.Digit)
	}
	if lex.Accept("e") {
		lex.Accept("-")
		lex.AcceptRunRange(unicode.Digit)
	}
	lex.Emit(ItemNumber)
	return Bracket
}

func BracketKey(lex *lexer.Lexer) lexer.StateFn {
	if lex.AcceptRunRange(unicode.Letter) > 0 {
		lex.Emit(ItemPathKey)
		switch r, _ := lex.Peek(); r {
		case lexer.EOF:
			return nil
		case '.':
			lex.Advance()
			lex.Emit(ItemDot)
			return BracketKey
		}
	}
	r, _ := lex.Peek()
	if r != lexer.EOF {
		lex.Errorf("expected key got EOF")
	}
	return nil
}

var PARSE_DEBUG = false

func debug(v ...interface{}) {
	if PARSE_DEBUG {
		fmt.Fprint(os.Stderr, v...)
	}
}

func debugln(v ...interface{}) {
	if PARSE_DEBUG {
		fmt.Fprintln(os.Stderr, v...)
	}
}

func debugf(format string, v ...interface{}) {
	if PARSE_DEBUG {
		fmt.Fprintf(os.Stderr, format, v...)
	}
}
