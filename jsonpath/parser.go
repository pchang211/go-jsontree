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

// // NewJSONPath instantiates and returns a new JSONPath object
// func NewJSONPath(f func(interface{}) interface{}) *JSONPath {
// 	initialTraverser := &Traverser{Traverse: f}
// 	return &JSONPath{
// 		head: initialTraverser,
// 		tail: initialTraverser,
// 	}
// }

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
func (j *JSONPath) TraverseJSON(json interface{}) interface{} {
	for {
		if j.head == nil {
			return json
		}
		json = j.head.Traverse(json)
		j.head = j.head.child
	}
}

// Traverser is really a linked list wrapper over Traverse() function.
// Traverse() advances through an input json object and returns the result
type Traverser struct {
	child    *Traverser
	Traverse func(interface{}) interface{}
}

// NewTraverser returns a new Traverser object
func NewTraverser(f func(interface{}) interface{}) *Traverser {
	return &Traverser{Traverse: f}
}

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
	debugln("parseBracket")
	return nil, fmt.Errorf("not implemented")
}
