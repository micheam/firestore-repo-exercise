package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/urfave/cli/v2"
	"micheam.com/todolist"
)

const (
	version = "0.1.0"
)

func mustNewClient(ctx context.Context) *firestore.Client {
	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		panic("PROJECT_ID is empty")
	}
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}
	return client
}

func newTaskRepository(ctx context.Context) todolist.TaskRepository {
	client := mustNewClient(ctx)
	middlewares := []todolist.TaskRepositoryMiddleware{}
	return todolist.NewTaskRepository(client, middlewares...)
}

func main() {
	if err := newApp().Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "todolist"
	app.Usage = "sample implementation of https://github.com/go-generalize/firestore-repo"
	app.Version = version
	app.Authors = []*cli.Author{
		{
			Name:  "Michto Maeda",
			Email: "michito.maeda@gmail.com",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:      "add",
			Action:    doAdd,
			Usage:     "add new task",
			ArgsUsage: "[task text]",
		},
		{
			Name:   "list",
			Action: doList,
			Usage:  "list existing tasks",
		},
	}
	return app
}

func doAdd(ctx *cli.Context) error {

	if ctx.Args().Len() == 0 {
		return errors.New("task must specified")
	}

	childCtx, cancel := context.WithTimeout(ctx.Context, 5*time.Second)
	defer cancel()

	task := todolist.NewTask(ctx.Args().First())
	id, err := newTaskRepository(ctx.Context).Insert(childCtx, task)
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}

	fmt.Print(id)
	return nil
}

func doList(ctx *cli.Context) error {

	childCtx, cancel := context.WithTimeout(ctx.Context, 5*time.Second)
	defer cancel()

	found, err := newTaskRepository(ctx.Context).List(childCtx, &todolist.TaskListReq{}, nil)
	if err != nil {
		return fmt.Errorf("failed to list tasks: %w", err)
	}

	for _, task := range found {
		b, _ := json.Marshal(task)
		fmt.Println(string(b))
	}

	return nil
}
