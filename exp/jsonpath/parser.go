// Copyright 2013, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// parser.go [created: Fri, 21 Jun 2013]

package jsonpath

import (
	"errors"
	"fmt"
)

const parseDebug = true

// Parse takes an input string, instantiates a lexer with the input
// reads each token until EOF and returns (for now) Selector functions
func Parse(input string) (Selector, error) {
	selectors := make([]Selector, 0, 1)
	lex := NewPathLexer(input)
	for {
		switch item := lex.Next(); item.Type {
		case itemEOF:
			debug("EOF\n")
			debugf("%d selectors\n", len(selectors))
			switch len(selectors) {
			case 0:
				return nil, fmt.Errorf("empty")
			case 1:
				return selectors[0], nil
			default:
				return Chain(selectors...), nil
			}
		case itemError:
			debug("ERROR\n")
			return nil, errors.New(item.Value)
		case itemDollar:
			debug("DOLLAR ")
			next, _ := lex.lex.Peek()
			if next != '.' {
				return nil, fmt.Errorf("expected \".\" but got %q", next)
			}
		case itemDotDot:
			debug("DOTDOT ")
			fallthrough // FIXME
		case itemDot:
			debug("DOT\n")
			switch next := lex.Next(); next.Type {
			case itemEOF:
				return nil, errors.New("unexpected EOF")
			case itemStarStar:
				debug("STAR STAR\n")
				selectors = append(selectors, RecursiveDescent)
			case itemStar:
				debug("STAR\n")
				selectors = append(selectors, All)
			case itemPathKey:
				debugf("PATH KEY %s\n", next.Value)
				selectors = append(selectors, Key(next.Value))
			default:
				return nil, fmt.Errorf("expected key but got %q", next.Value)
			}
		case itemLeftBracket:
			debug("LEFTBRACKET\n")
			sel, err := parseBracket(lex)
			if err != nil {
				return nil, err
			}
			selectors = append(selectors, sel)
		}
	}
}

func parseBracket(lex *PathLexer) (Selector, error) {
	debug("parseBracket")
	return nil, fmt.Errorf("not implemented")
}
