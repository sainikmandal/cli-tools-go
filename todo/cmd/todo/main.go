package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sainikmandal/cli-tools-go/todo"
)

const todoFileName = ".todo.json"

func main() {
	add := flag.Bool("add", false, "add the remaining arguments as a task")
	list := flag.Bool("list", false, "list tasks")
	complete := flag.Int("complete", 0, "complete task `number`")
	del := flag.Int("delete", 0, "delete task `number`")
	flag.Parse()

	l := &todo.List{}
	if err := l.Get(todoFileName); err != nil {
		fatal(err)
	}

	switch {
	case *add:
		task := strings.Join(flag.Args(), " ")
		if task == "" {
			fmt.Fprintln(os.Stderr, "no task given")
			os.Exit(2)
		}
		l.Add(task)
		save(l)
	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			fatal(err)
		}
		save(l)
	case *del > 0:
		if err := l.Delete(*del); err != nil {
			fatal(err)
		}
		save(l)
	case *list:
		for i, it := range *l {
			mark := " "
			if it.Done {
				mark = "x"
			}
			fmt.Printf("%d. [%s] %s\n", i+1, mark, it.Task)
		}
	default:
		fmt.Fprintln(os.Stderr, "specify one of -add, -list, -complete, -delete")
		os.Exit(2)
	}
}

func save(l *todo.List) {
	if err := l.Save(todoFileName); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
