// Copyright 2025 Abdulrahman Abdulhamid. All rights reserved.
// Use of this source code is governed by Apache-2.0 
// license that can be found in the LICENSE file.

package peruse

import(
	"fmt"
)

type Location interface {
	Origin() string
	Line() int
	Column() int
	String() string
}

func NewLocation(origin string, line, column int) Location {
	return location{origin, line, column}
	
}

type location struct {
	origin string
	line   int
	column int	
}

func (l location) Origin() string {
	return l.origin
}

func (l location) Line() int {
	return l.line
}

func (l location) Column() int {
	return l.column
}

func (l location) String() string {
	return fmt.Sprintf("%s:%d:%d", l.origin, l.line, l.column)
}

