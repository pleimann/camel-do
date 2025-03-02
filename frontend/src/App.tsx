import { TaskService } from '@bindings/pleimann.com/camel-do/services'

import { createResource, ErrorBoundary, Suspense } from 'solid-js';

import TitleBar from '@/components/TitleBar';
import Backlog from '@/components/Backlog';

const App = () => {
  const [tasks] = createResource(async () => await TaskService.GetTasks());

  return (
    <div class="mt-[64px]">
      <TitleBar />

      <main class="h-[calc(100dvh-64px)] w-full overflow-hidden">
        <ErrorBoundary fallback={(err) => <div>Error: {err.message}</div>}>
          <Suspense ref={tasks()}>
            <Backlog tasks={tasks() || []} />
          </Suspense>
        </ErrorBoundary>
      </main>
    </div>
  );
};

export default App;