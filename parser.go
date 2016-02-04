// updated by Philip Chang
// originally forked from Bryan Matsuo

package jsonpath

import (
	"errors"
	"fmt"
)

const parseDebug = false

// Parse takes an input string, instantiates a lexer with the input
// reads each token until EOF and returns (for now) Selector functions
func Parse(input string) (*JSONPath, error) {
	path := &JSONPath{}
	lex := NewPathLexer(input)
	for {
		switch item := lex.Next(); item.Type {
		case itemEOF:
			debug("EOF\n")
			return path, nil
		case itemError:
			debug("ERROR\n")
			return nil, errors.New(item.Value)
		case itemDollar:
			debug("DOLLAR ")
			next, _ := lex.lex.Peek()
			if next != '.' {
				return nil, fmt.Errorf("expected \".\" but got %q", next)
			}
		case itemDot:
			debug("DOT\n")
			switch next := lex.Next(); next.Type {
			case itemEOF:
				return nil, errors.New("unexpected EOF")
			// case itemStarStar:
			// 	debug("STAR STAR\n")
			case itemStar:
				debug("STAR\n")
				path.AddTraverser(NewTraverser(Star()))
			case itemPathKey:
				debugf("PATH KEY %s\n", next.Value)
				path.AddTraverser(NewTraverser(Key(next.Value)))
			default:
				return nil, fmt.Errorf("expected key but got %q", next.Value)
			}
		case itemLeftBracket:
			debug("LEFTBRACKET\n")
			switch next := lex.Next(); next.Type {
			case itemNumber:
				path.AddTraverser(NewTraverser(IndexKey(next.Value)))
				if next = lex.Next(); next.Type != itemRightBracket {
					return nil, fmt.Errorf("expected ']'")
				}
			case itemRightBracket:
				// no index or query
				debugf("right bracket")
			}
		case itemRightBracket:
			debug("RIGHTBRACKET\n")
			return nil, fmt.Errorf("close bracket ] seen without open bracket")
		case itemLess:
			switch next := lex.Next(); next.Type {
			case itemEOF:
				return nil, errors.New("unexpexted EOF")
			case itemPathKey:
				debugf("PATH KEY %s\n", next.Value)
				return nil, errors.New("less than comparison must be with a number")
			case itemNumber:
				debugf("NUMBER: %v\n", next.Value)
				path.AddTraverser(NewTraverser(LessThan(next.Value)))
			default:
				return nil, errors.New("don't know how to compare that")
			}
		case itemGreater:
			switch next := lex.Next(); next.Type {
			case itemEOF:
				return nil, errors.New("unexpexted EOF")
			case itemPathKey:
				debugf("PATH KEY %s\n", next.Value)
				return nil, errors.New("less than comparison must be with a number")
			case itemNumber:
				debugf("NUMBER: %v\n", next.Value)
				path.AddTraverser(NewTraverser(GreaterThan(next.Value)))
			default:
				return nil, errors.New("don't know how to compare that")
			}
		default:
			debugf("Not sure what to do with %v", item.Type)
		}
	}
}
