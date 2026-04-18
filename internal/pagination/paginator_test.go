package pagination_test

import (
	"testing"

	"github.com/yourusername/vaultpulse/internal/pagination"
)

func makeItems(n int) []string {
	items := make([]string, n)
	for i := range items {
		items[i] = fmt.Sprintf("secret/path/%d", i)
	}
	return items
}

import "fmt"

func TestNew_InvalidPageSize(t *testing.T) {
	_, err := pagination.New([]string{"a"}, 0)
	if err == nil {
		t.Fatal("expected error for page size 0")
	}
}

func TestNew_ValidPageSize(t *testing.T) {
	p, err := pagination.New([]string{"a", "b"}, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Len() != 2 {
		t.Fatalf("expected len 2, got %d", p.Len())
	}
}

func TestNext_FirstPage(t *testing.T) {
	items := []string{"a", "b", "c", "d", "e"}
	p, _ := pagination.New(items, 2)
	page := p.Next(0)
	if len(page.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(page.Items))
	}
	if !page.HasMore {
		t.Fatal("expected HasMore=true")
	}
	if page.Cursor != 2 {
		t.Fatalf("expected cursor 2, got %d", page.Cursor)
	}
}

func TestNext_LastPage(t *testing.T) {
	items := []string{"a", "b", "c"}
	p, _ := pagination.New(items, 2)
	page := p.Next(2)
	if len(page.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(page.Items))
	}
	if page.HasMore {
		t.Fatal("expected HasMore=false on last page")
	}
}

func TestNext_BeyondEnd(t *testing.T) {
	p, _ := pagination.New([]string{"a"}, 5)
	page := p.Next(10)
	if len(page.Items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(page.Items))
	}
	if page.HasMore {
		t.Fatal("expected HasMore=false")
	}
}

func TestNext_NegativeCursor(t *testing.T) {
	p, _ := pagination.New([]string{"x", "y"}, 1)
	page := p.Next(-5)
	if page.Items[0] != "x" {
		t.Fatalf("expected first item 'x', got %s", page.Items[0])
	}
}
