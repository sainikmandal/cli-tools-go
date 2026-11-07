package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"slices"
	"time"
)

var ErrNotFound = errors.New("item not found")

type item struct {
	Task        string    `json:"task"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at,omitzero"`
}

type List []item

func (l *List) Add(task string) {
	*l = append(*l, item{Task: task, CreatedAt: time.Now()})
}

func (l *List) Complete(i int) error {
	if err := l.check(i); err != nil {
		return err
	}
	(*l)[i-1].Done = true
	(*l)[i-1].CompletedAt = time.Now()
	return nil
}

func (l *List) Delete(i int) error {
	if err := l.check(i); err != nil {
		return err
	}
	*l = slices.Delete(*l, i-1, i)
	return nil
}

func (l *List) Save(filename string) error {
	data, err := json.Marshal(l)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0o644)
}

func (l *List) Get(filename string) error {
	data, err := os.ReadFile(filename)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, l)
}

func (l *List) check(i int) error {
	if i < 1 || i > len(*l) {
		return fmt.Errorf("%w: %d", ErrNotFound, i)
	}
	return nil
}
