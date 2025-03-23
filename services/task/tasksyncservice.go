package task

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"slices"
	"strconv"

	"github.com/pleimann/camel-do/model"
	"github.com/pleimann/camel-do/services/db"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

// TaskService is a service for managing tasks.
type TaskSyncService struct {
	googleTasks *tasks.Service
	db          *db.DatabaseService
}

func NewTaskSyncService(http *http.Client, db *db.DatabaseService) (*TaskSyncService, error) {
	service, err := tasks.NewService(context.Background(), option.WithHTTPClient(http))

	if err != nil {
		slog.Error("error creating google tasks service", "error", err.Error())
		return nil, err
	}

	syncService := &TaskSyncService{
		googleTasks: service,
		db:          db,
	}

	return syncService, nil
}

func (t *TaskSyncService) GetGoogleTasks() []model.Task {
	r, err := t.googleTasks.Tasklists.List().MaxResults(10).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve task lists. %v", err)
	}

	tasks := []model.Task{}

	if len(r.Items) > 0 {
		taskList := r.Items[0]
		gtasks, err := t.googleTasks.Tasks.List(taskList.Id).Do()
		if err != nil {
			return []model.Task{}
		}

		for i, gtask := range gtasks.Items {
			b, _ := gtask.MarshalJSON()

			fmt.Printf("%d gtask: %v\n", i, string(b))

			order, _ := strconv.Atoi(gtask.Position)

			tasks = append(tasks, model.Task{
				ID:          gtask.Id,
				Title:       gtask.Title,
				Description: gtask.Notes,
				Completed:   gtask.Completed != nil,
				Rank:        order,
			})
		}

		slices.SortStableFunc(tasks, func(a, b model.Task) int {
			return cmp.Compare(a.Rank, b.Rank)
		})

		return tasks

	} else {
		fmt.Print("No task lists found.")
	}

	return tasks
}
