package task

import (
	"cmp"
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"slices"
	"strconv"

	"github.com/guregu/null/v6/zero"
	"github.com/pleimann/camel-do/model"
	"github.com/pleimann/camel-do/services/oauth"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
)

// TaskService is a service for managing tasks.
type TaskSyncService struct {
	googleTasks *tasks.Service
	db          *sql.DB
}

func NewTaskSyncService(db *sql.DB) (*TaskSyncService, error) {
	http := oauth.NewGoogleAuth().GetClient()

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
	r, err := t.googleTasks.Tasklists.
		List().
		MaxResults(10).
		Do()

	if err != nil {
		log.Fatalf("Unable to retrieve task lists. %v", err)
	}

	modelTasks := []model.Task{}

	if len(r.Items) > 0 {
		taskList := r.Items[0]
		gtasks, err := t.googleTasks.Tasks.
			List(taskList.Id).
			Do()

		if err != nil {
			return []model.Task{}
		}

		for _, gtask := range gtasks.Items {
			order, _ := strconv.ParseInt(gtask.Position, 10, 32)

			modelTasks = append(modelTasks, model.Task{
				GTaskID:     zero.StringFrom(gtask.Id),
				Title:       zero.StringFrom(gtask.Title),
				Description: zero.StringFrom(gtask.Notes),
				Completed:   zero.BoolFrom(gtask.Completed != nil),
				Rank:        zero.Int32From(int32(order)),
			})
		}

		slices.SortStableFunc(modelTasks, func(a, b model.Task) int {
			return cmp.Compare(a.Rank.Int32, b.Rank.Int32)
		})

		return modelTasks

	} else {
		fmt.Print("No task lists found.")
	}

	return modelTasks
}
