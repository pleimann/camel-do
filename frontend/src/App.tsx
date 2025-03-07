import { createResource, ErrorBoundary, Suspense } from 'solid-js';
import { Portal } from "solid-js/web";

import { TaskService } from '@bindings/pleimann.com/camel-do/services/task';

import TitleBar, { TitleBarAction } from '@/components/TitleBar';
import Backlog from '@/components/Backlog';
import Timeline from '@/components/Timeline';
import Icon from '@/components/Icon';
import TaskDialog from './TaskDialog';

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

  let modal: HTMLDialogElement | undefined;

  const newTask = async () => {
    modal?.showModal()
  }

  return (
    <>
      <div>
        <TitleBar onAction={onTitleBarAction} />

        <main class="h-[calc(100dvh-(var(--spacing)*16))] max-h-[calc(100dvh-(var(--spacing)*16))] w-full overflow-hidden">
          <ErrorBoundary fallback={(err) => <div>Error: {err.message}</div>}>
            <Suspense ref={tasks()}>
              <div class="flex flex-row h-full">
                <div class="pr-1 bg-primary/40 dark:bg-primary-content h-full min-w-75">
                  <Backlog tasks={tasks() || []} />
                </div>
                <div class="grow overflow-y-auto m-1 p-8">
                  <Timeline tasks={tasks() || []} />
                </div>
              </div>
            </Suspense>
            <button class="btn btn-circle size-16 fixed bottom-4 right-4 shadow-xl bg-primary/75 text-primary-content">
              <Icon.Plus onClick={newTask} class="size-16" />
            </button>
          </ErrorBoundary>
        </main>
      </div>

      <Portal>
        <TaskDialog ref={modal} />
      </Portal>
    </>
  );
};

export default App;