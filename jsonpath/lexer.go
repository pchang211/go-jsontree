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

// PathLexer is a wrapper around the go-lexer lexer object
type PathLexer struct {
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
		lex.Emit(itemDot)
		return Start
	case 2:
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
		lex.Emit(itemStar)
		return Start
	case 2:
		lex.Emit(itemStarStar)
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
		lex.Advance()
		lex.Emit(itemLeftBracket)
		return Bracket
	case r == ']':
		lex.Advance()
		lex.Emit(itemRightBracket)
		return Start
	case r == '$':
		lex.Advance()
		lex.Emit(itemDollar)

	// (currently) unused characters
	case r == '>':
		lex.Advance()
		if r, _ := lex.Peek(); r == '=' {
			lex.Advance()
			lex.Emit(itemGreaterEqual)
		} else {
			lex.Emit(itemGreater)
		}
	case r == '<':
		lex.Advance()
		if r, _ := lex.Peek(); r == '=' {
			lex.Advance()
			lex.Emit(itemLessEqual)
		} else {
			lex.Emit(itemLess)
		}
	case r == '=':
		lex.Advance()
		lex.Emit(itemEqual)
	case r == '!':
		lex.Advance()
		if r, _ := lex.Peek(); r == '=' {
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

// Bracket state
func Bracket(lex *lexer.Lexer) lexer.StateFn {
	switch r, _ := lex.Peek(); {
	case r == lexer.EOF:
		lex.Emit(itemError)
		return nil
	case unicode.IsDigit(r):
		return Start
	}
	return nil
}

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
