package todo_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/sainikmandal/cli-tools-go/todo"
)

func TestAdd(t *testing.T) {
	l := todo.List{}

	taskname := "New Task"
	l.Add(taskname)

	if l[0].Task != taskname {
		t.Errorf("Expected %q, got %q instead", taskname, l[0].Task)
	}
	if l[0].Done {
		t.Errorf("New task shouldn't be completed")
	}
	if l[0].CreatedAt.IsZero() {
		t.Errorf("CreatedAt should be set")
	}
}

func TestComplete(t *testing.T) {
	l := todo.List{}
	l.Add("New Task")

	if err := l.Complete(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !l[0].Done {
		t.Errorf("Task should have been completed")
	}
	if l[0].CompletedAt.IsZero() {
		t.Errorf("CompletedAt should be set")
	}
}

func TestCompleteNotFound(t *testing.T) {
	l := todo.List{}
	l.Add("New Task")

	for _, i := range []int{0, -1, 2} {
		if err := l.Complete(i); !errors.Is(err, todo.ErrNotFound) {
			t.Errorf("Complete(%d): expected ErrNotFound, got %v", i, err)
		}
	}
}

func TestDelete(t *testing.T) {
	l := todo.List{}
	for _, task := range []string{"a", "b", "c"} {
		l.Add(task)
	}

	if err := l.Delete(2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(l) != 2 {
		t.Fatalf("Expected 2 items, got %d", len(l))
	}
	if l[0].Task != "a" || l[1].Task != "c" {
		t.Errorf("Unexpected remaining tasks: %q, %q", l[0].Task, l[1].Task)
	}
}

func TestDeleteNotFound(t *testing.T) {
	l := todo.List{}
	l.Add("a")

	for _, i := range []int{0, -1, 2} {
		if err := l.Delete(i); !errors.Is(err, todo.ErrNotFound) {
			t.Errorf("Delete(%d): expected ErrNotFound, got %v", i, err)
		}
	}
	if len(l) != 1 {
		t.Errorf("Expected list unchanged, got %d items", len(l))
	}
}

func TestSaveGet(t *testing.T) {
	l1 := todo.List{}
	l1.Add("New Task")
	l1.Add("Second Task")
	if err := l1.Complete(2); err != nil {
		t.Fatal(err)
	}

	filename := filepath.Join(t.TempDir(), "todo.json")
	if err := l1.Save(filename); err != nil {
		t.Fatalf("Save: %v", err)
	}

	l2 := todo.List{}
	if err := l2.Get(filename); err != nil {
		t.Fatalf("Get: %v", err)
	}

	if len(l2) != len(l1) {
		t.Fatalf("Expected %d items, got %d", len(l1), len(l2))
	}
	for i := range l1 {
		if l1[i].Task != l2[i].Task || l1[i].Done != l2[i].Done {
			t.Errorf("Item %d mismatch: %+v vs %+v", i, l1[i], l2[i])
		}
		if !l1[i].CreatedAt.Equal(l2[i].CreatedAt) || !l1[i].CompletedAt.Equal(l2[i].CompletedAt) {
			t.Errorf("Item %d timestamps mismatch: %+v vs %+v", i, l1[i], l2[i])
		}
	}
}

func TestGetMissingFile(t *testing.T) {
	l := todo.List{}
	if err := l.Get(filepath.Join(t.TempDir(), "missing.json")); err != nil {
		t.Errorf("Expected no error for missing file, got %v", err)
	}
	if len(l) != 0 {
		t.Errorf("Expected empty list, got %d items", len(l))
	}
}

func TestGetEmptyFile(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "empty.json")
	if err := os.WriteFile(filename, nil, 0o644); err != nil {
		t.Fatal(err)
	}

	l := todo.List{}
	if err := l.Get(filename); err != nil {
		t.Errorf("Expected no error for empty file, got %v", err)
	}
	if len(l) != 0 {
		t.Errorf("Expected empty list, got %d items", len(l))
	}
}

func TestGetInvalidJSON(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(filename, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}

	l := todo.List{}
	if err := l.Get(filename); err == nil {
		t.Errorf("Expected error for invalid JSON")
	}
}
