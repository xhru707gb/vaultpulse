// Package pagination provides cursor-based pagination for listing Vault secret paths.
package pagination

import (
	"errors"
)

// ErrInvalidPageSize is returned when page size is less than 1.
var ErrInvalidPageSize = errors.New("pagination: page size must be at least 1")

// Page holds a single page of results and the cursor for the next page.
type Page[T any] struct {
	Items  []T
	Cursor int
	HasMore bool
}

// Paginator pages through a slice of items.
type Paginator[T any] struct {
	items    []T
	pageSize int
}

// New creates a Paginator for the given items and page size.
func New[T any](items []T, pageSize int) (*Paginator[T], error) {
	if pageSize < 1 {
		return nil, ErrInvalidPageSize
	}
	return &Paginator[T]{items: items, pageSize: pageSize}, nil
}

// Next returns the page starting at cursor. cursor=0 begins from the start.
func (p *Paginator[T]) Next(cursor int) Page[T] {
	if cursor < 0 {
		cursor = 0
	}
	if cursor >= len(p.items) {
		return Page[T]{Items: []T{}, Cursor: cursor, HasMore: false}
	}
	end := cursor + p.pageSize
	hasMore := end < len(p.items)
	if end > len(p.items) {
		end = len(p.items)
	}
	return Page[T]{
		Items:   p.items[cursor:end],
		Cursor:  end,
		HasMore: hasMore,
	}
}

// Len returns the total number of items.
func (p *Paginator[T]) Len() int { return len(p.items) }
