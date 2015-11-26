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
	itemEOF                    = lexer.ItemEOF
	itemError                  = lexer.ItemError
	itemPathKey lexer.ItemType = iota
	itemDollar
	itemDot
	itemDotDot
	itemStar
	itemStarStar
	itemLeftBracket
	itemRightBracket
	itemNumber
	itemNull
	itemTrue
	itemFalse
	itemEqual
	itemGreater
	itemGreaterEqual
	itemLess
	itemLessEqual
	itemNotEqual
	ittemString
)

// PathLexer is a wrapper around go-lexer lexer object
type PathLexer struct {
	// Next() *lexer.Item
	lex *lexer.Lexer
}

// NewPathLexer returns a new PathLexer object
func NewPathLexer(input string) *PathLexer {
	return &PathLexer{lex: lexer.New(Start, input)}
}

// Next calls next on the lexer object
func (l *PathLexer) Next() *lexer.Item {
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
		// debugln("FOUND DOT")
		lex.Emit(itemDot)
		return Start
	case 2:
		// debugln("FOUND DOT DOT")
		lex.Emit(itemDotDot)
		return Start
	default:
		return lex.Errorf("unexpected '.'")
	}

	// wildcard
	switch lex.AcceptRun("*") {
	case 0:
		break
	case 1:
		// debugln("FOUND STAR")
		lex.Emit(itemStar)
		return Start
	// why do we even need STAR STAR?
	// case 2:
	// 	// debugln("FOUND STAR STAR")
	// 	lex.Emit(itemStarStar)
	// 	return Start
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
		// debugln("FOUND LEFT BRACKET")
		lex.Advance()
		lex.Emit(itemLeftBracket)
		// return Bracket
	case r == ']':
		// debugln("FOUND RIGHT BRACKET")
		lex.Advance()
		lex.Emit(itemRightBracket)
	case r == '$':
		// debugln("FOUND DOLLAR")
		lex.Advance()
		lex.Emit(itemDollar)
	case r == '>':
		lex.Advance()
		if r, _ := lex.Peek(); r == '=' {
			// debugln("FOUND GREATER EQUAL")
			lex.Advance()
			lex.Emit(itemGreaterEqual)
		} else {
			// debugln("FOUND GREATER")
			lex.Emit(itemGreater)
		}
	case r == '<':
		lex.Advance()
		if r, _ := lex.Peek(); r == '=' {
			// debugln("FOUND LESS EQUAL")
			lex.Advance()
			lex.Emit(itemLessEqual)
		} else {
			// debugln("FOUND LESS")
			lex.Emit(itemLess)
		}
	case r == '=':
		// debugln("FOUND EQUAL")
		lex.Advance()
		lex.Emit(itemEqual)
	case r == '!':
		lex.Advance()
		if r, _ := lex.Peek(); r == '=' {
			// debugln("FOUND NOT EQUAL")
			lex.Advance()
			lex.Emit(itemNotEqual)
		} else {
			return lex.Errorf("expected '=' got %c", r)
		}
	}

	// default
	return Start
}

// PathKey state. Keep going while we hit letters
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
		// debugln("FOUND PATH KEY")
		lex.Emit(itemPathKey)
	}
	return Start
}

// Number state. keep going until we can parse the number
func Number(lex *lexer.Lexer) lexer.StateFn {
	// keep advancing as long as we see digits
	if lex.AcceptRunRange(unicode.Digit) == 0 {
		r, _ := lex.Peek()
		// (PC) I don't believe this catches the `.12345` decimal case
		if r == lexer.EOF {
			return nil
		}
		return lex.Errorf("expected digit got %c", r)
	}
	// debugln("FOUND NUMBER")
	if lex.Accept(".") {
		lex.AcceptRunRange(unicode.Digit)
	}
	if lex.Accept("eE") {
		lex.Accept("-")
		lex.AcceptRunRange(unicode.Digit)
	}
	lex.Emit(itemNumber)
	return Start
}

// func StartBracket(lex *lexer.Lexer) lexer.StateFn {
// 	if lex.Accept("[") {
// 		lex.Emit(itemLeftBracket)
// 		return Bracket
// 	} else {
// 		lex.Errorf("missing '['")
// 		return nil
// 	}
// }
//

//
// // Bracket state
// func Bracket(lex *lexer.Lexer) lexer.StateFn {
// 	switch r, _ := lex.Peek(); {
// 	case r == lexer.EOF:
// 		return nil
// 	case unicode.IsDigit(r):
// 		return BracketNumber
// 	case unicode.IsLetter(r):
// 		return BracketKey
// 	case r == '>':
// 		lex.Advance()
// 		if r, _ := lex.Peek(); r == '=' {
// 			lex.Advance()
// 			lex.Emit(itemGreaterEqual)
// 		} else {
// 			lex.Emit(itemGreater)
// 		}
// 	case r == '<':
// 		lex.Advance()
// 		if r, _ := lex.Peek(); r == '=' {
// 			lex.Advance()
// 			lex.Emit(itemLessEqual)
// 		} else {
// 			lex.Emit(itemLess)
// 		}
// 	case r == '=':
// 		lex.Advance()
// 		lex.Emit(itemEqual)
// 	case r == '!':
// 		lex.Advance()
// 		if r, _ := lex.Peek(); r == '=' {
// 			lex.Advance()
// 			lex.Emit(itemNotEqual)
// 		} else {
// 			lex.Errorf("expected '=' got %c", r)
// 			return nil
// 		}
// 	}
// 	return Bracket
// }
//
// func BracketNumber(lex *lexer.Lexer) lexer.StateFn {
// 	if lex.AcceptRunRange(unicode.Digit) == 0 {
// 		r, _ := lex.Peek()
// 		if r == lexer.EOF {
// 			return nil
// 		}
// 		lex.Errorf("expected digit got %c", r)
// 	}
// 	if lex.Accept(".") {
// 		lex.AcceptRunRange(unicode.Digit)
// 	}
// 	if lex.Accept("e") {
// 		lex.Accept("-")
// 		lex.AcceptRunRange(unicode.Digit)
// 	}
// 	lex.Emit(itemNumber)
// 	return Bracket
// }
//
// func BracketKey(lex *lexer.Lexer) lexer.StateFn {
// 	if lex.AcceptRunRange(unicode.Letter) > 0 {
// 		lex.Emit(itemPathKey)
// 		switch r, _ := lex.Peek(); r {
// 		case lexer.EOF:
// 			return nil
// 		case '.':
// 			lex.Advance()
// 			lex.Emit(itemDot)
// 			return BracketKey
// 		}
// 	}
// 	r, _ := lex.Peek()
// 	if r != lexer.EOF {
// 		lex.Errorf("expected key got EOF")
// 	}
// 	return nil
// }

func debug(v ...interface{}) {
	if parseDebug {
		fmt.Fprint(os.Stderr, v...)
	}
}

func debugln(v ...interface{}) {
	if parseDebug {
		fmt.Fprintln(os.Stderr, v...)
	}
}

func debugf(format string, v ...interface{}) {
	if parseDebug {
		fmt.Fprintf(os.Stderr, format, v...)
	}
}
