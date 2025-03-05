import { TaskService } from '@bindings/pleimann.com/camel-do/services'

import { createResource, ErrorBoundary, Suspense } from 'solid-js';

import TitleBar, { TitleBarAction } from '@/components/TitleBar';
import Backlog from '@/components/Backlog';
import Timeline from '@/components/Timeline';

const App = () => {
  const [tasks, { refetch }] = createResource(async () => await TaskService.GetTasks());

  const onTitleBarAction = (action: TitleBarAction) => {
    switch(action) {
      case "refresh":
        refetch();
        break;
      case "search":
        break;
    }
  }

  return (
    <div>
      <TitleBar onAction={onTitleBarAction} />

      <main class="h-[calc(100dvh-(var(--spacing)*16))] max-h-[calc(100dvh-(var(--spacing)*16))] w-full overflow-hidden">
        <ErrorBoundary fallback={(err) => <div>Error: {err.message}</div>}>
          <Suspense ref={tasks()}>
            <div class="flex flex-row h-full">
              <div class="py-4 pr-1 bg-primary-content/50 dark:bg-primary-content h-full min-w-75">
                <Backlog tasks={tasks() || []} />
              </div>
              <div class="grow overflow-y-auto m-1 p-8">
                <Timeline tasks={tasks() || []} />
              </div>
            </div>
          </Suspense>
        </ErrorBoundary>
      </main>
    </div>
  );
};

export default App;