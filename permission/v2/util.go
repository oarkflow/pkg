package v2

import (
	"strings"
	"sync"
)

type Pool[T any] struct {
	syncPool sync.Pool
}

func NewPool[T any](size int) *Pool[T] {
	return &Pool[T]{
		syncPool: sync.Pool{New: func() any {
			return make([]T, 0, size)
		}},
	}
}

func (p *Pool[T]) Get() []T {
	return p.syncPool.Get().([]T)
}

func (p *Pool[T]) Put(s []T) {
	p.syncPool.Put(s[:0])
}

var (
	stringSlice   = NewPool[string](100)
	userRoleSlice = NewPool[*UserRole](100)
)

func MatchResource(value, pattern string) bool {
	vIndex, pIndex := 0, 0
	vLen, pLen := len(value), len(pattern)

	for pIndex < pLen {
		if pattern[pIndex] == '*' {
			// If '*' is the last character in the pattern, it matches everything
			if pIndex == pLen-1 {
				return true
			}

			// Find the next character in pattern after '*'
			nextChar := pattern[pIndex+1]

			// If the next character is '*', skip it
			if nextChar == '*' {
				pIndex++
				continue
			}

			// Find the next occurrence of the character after '*' in the value
			nextIndex := strings.IndexByte(value[vIndex:], nextChar)

			// If the character is not found, no match
			if nextIndex == -1 {
				return false
			}

			// Move the value index to the next occurrence of the character
			vIndex += nextIndex
		} else if pIndex < pLen && vIndex < vLen && (pattern[pIndex] == value[vIndex] || pattern[pIndex] == ':') {
			// If pattern part matches value part or is a parameter, move to the next parts
			vIndex++
			pIndex++
			// If pattern part is a parameter, skip it in the value
			if pattern[pIndex-1] == ':' {
				// Find the end of the parameter segment
				endIndex := pIndex
				for endIndex < pLen && pattern[endIndex] != '/' {
					endIndex++
				}
				// Skip the parameter segment in the value
				for vIndex < vLen && value[vIndex] != '/' {
					vIndex++
				}
				// Move pattern index to the end of the parameter segment
				pIndex = endIndex
			}
		} else {
			return false
		}
	}

	// If both value and pattern are exhausted, return true
	return vIndex == vLen && pIndex == pLen
}
