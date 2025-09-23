// Copyright 2025 Abdulrahman Abdulhamid. All rights reserved.
// Use of this source code is governed by Apache-2.0 
// license that can be found in the LICENSE file.

package numeric

func Numbers(many ...Number) Number {
	return numbers{many}
}

type numbers struct {
	many []Number
}

func (n numbers) Scan(content []rune) []rune {
	for _, number := range n.many {
		if result := number.Scan(content); len(result) != 0 {
			return result
		}
	}
	return []rune{}
}

