import { TaskService } from '@bindings/pleimann.com/camel-do/services'

import { createResource, ErrorBoundary } from 'solid-js';

import TitleBar from '@/components/TitleBar';
import Backlog from '@/components/Backlog';

const App = () => {
  const [tasks] = createResource(async () => {
    return await TaskService.GetTasks();
  });

  return (
    <div class="mt-[64px]">
      <TitleBar />

      <main class="h-[calc(100dvh-64px)] w-full overflow-hidden">
        <ErrorBoundary fallback={(err) => <div>Error: {err.message}</div>}>
          <Backlog tasks={tasks()} />
        </ErrorBoundary>
      </main>
    </div>
  );
};

export default App;