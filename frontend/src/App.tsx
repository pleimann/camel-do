import { TaskService } from '@bindings/pleimann.com/camel-do/services'

import { createResource, ErrorBoundary, Suspense } from 'solid-js';

import TitleBar from '@/components/TitleBar';
import Backlog from '@/components/Backlog';
import Timeline from '@/components/Timeline';

const App = () => {
  const [tasks] = createResource(async () => await TaskService.GetTasks());

  return (
    <div class="mt-[64px]">
      <TitleBar />

      <main class="h-[calc(100dvh-64px)] max-h-[calc(100dvh-64px)] w-full overflow-hidden">
        <ErrorBoundary fallback={(err) => <div>Error: {err.message}</div>}>
          <Suspense ref={tasks()}>
            <div class="flex flex-row h-full w-full">
              <div class="py-4 pl-4 bg-primary-content w-96 h-full shadow-xl">
                <Backlog tasks={tasks() || []} />
              </div>
              <div class="grow place-items-center h-full overflow-y-auto">
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