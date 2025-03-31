package task

import (
	"cmp"
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"slices"
	"strconv"

	"github.com/guregu/null/v6/zero"
	"github.com/pleimann/camel-do/model"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

// TaskService is a service for managing tasks.
type TaskSyncService struct {
	googleTasks *tasks.Service
	db          *sql.DB
}

func NewTaskSyncService(http *http.Client, db *sql.DB) (*TaskSyncService, error) {
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
				GTaskId:     gtask.Id,
				Title:       gtask.Title,
				Description: zero.StringFrom(gtask.Notes),
				Completed:   gtask.Completed != nil,
				Rank:        int32(order),
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
