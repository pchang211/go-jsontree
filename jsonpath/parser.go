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

// JSONPath is an object that can take in a json object, traverse according
// to the rules in the traverser, and return the resulting json. Underlying
// implementation is a singly linked list of Traverser objects
type JSONPath struct {
	head *Traverser
	tail *Traverser
}

// AddTraverser appends a new traverser to the JSONPath's linked list
// of Traverser objects
func (j *JSONPath) AddTraverser(traverser *Traverser) {
	// unitialized JSONPath
	if j.head == nil {
		j.head = traverser
	} else {
		j.tail.child = traverser
	}
	j.tail = traverser
}

// TraverseJSON takes in a json object and returns the subobject specified
// by the JSONPath
func (j *JSONPath) TraverseJSON(json interface{}) (interface{}, error) {
	for {
		if j.head == nil {
			return json, nil
		}
		var err error
		json, err = j.head.Traverse(json)
		if err != nil {
			return nil, err
		}
		j.head = j.head.child
	}
}

// TraverseFunc is a function that takes in json and returns json.
// Should traverse through the input json
type TraverseFunc func(interface{}) (interface{}, error)

// Traverser is really a linked list wrapper over Traverse() functions.
// Traverse() advances through an input json object and returns the result
type Traverser struct {
	child    *Traverser
	Traverse TraverseFunc
}

// NewTraverser returns a new Traverser object
func NewTraverser(f TraverseFunc) *Traverser {
	return &Traverser{Traverse: f}
}

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
			case itemStar:
				debug("STAR\n")
				path.AddTraverser(NewTraverser(TStar()))
			case itemPathKey:
				debugf("PATH KEY %s\n", next.Value)
				path.AddTraverser(NewTraverser(TKey(next.Value)))
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
		default:
			debugf("Not sure what to do with %v", item.Type)
		}
	}
}

func parseBracket(lex *PathLexer) (Selector, error) {
	debugln("parseBracket")
	return nil, fmt.Errorf("not implemented")
}
