package google

import (
	"context"
	"fmt"

	"google.golang.org/api/option"
	tasks "google.golang.org/api/tasks/v1"
	"pleimann.com/camel-do/services/google/oauth"
)

type GoogleTasksService struct {
	tokenSourceProvider *oauth.TokenSourceProvider
}

// NewTasksService creates a new Google Tasks service client.
func NewGoogleTasksService(
	ctx context.Context,
	tokenSourceProvider *oauth.TokenSourceProvider,
) (*GoogleTasksService, error) {
	service := &GoogleTasksService{
		tokenSourceProvider: tokenSourceProvider,
	}

	return service, nil
}

func (s *GoogleTasksService) GetTasks(ctx context.Context) ([]*tasks.Task, error) {
	service, err := s.googleTasksService(ctx)
	if err != nil {
		return nil, err
	}

	taskLists, err := service.Tasklists.List().Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve task lists: %v", err)
	}

	if len(taskLists.Items) == 0 {
		return nil, fmt.Errorf("no task lists found")
	}

	taskList := taskLists.Items[0]

	tasks, err := service.Tasks.List(taskList.Id).Do()

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve tasks: %v", err)
	}

	return tasks.Items, nil
}

func (s *GoogleTasksService) googleTasksService(ctx context.Context) (*tasks.Service, error) {
	tokenSource, err := (*s.tokenSourceProvider)(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get token source: %v", err)
	}

	srv, err := tasks.NewService(ctx,
		option.WithTokenSource(tokenSource),
		option.WithScopes(tasks.TasksReadonlyScope, tasks.TasksScope),
	)

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Tasks client: %v", err)
	}

	return srv, nil
}
