import { TaskService } from '@bindings/pleimann.com/camel-do/services'

import { createResource, ErrorBoundary } from 'solid-js';

import TitleBar from '@/components/TitleBar';
import Schedule from '@/components/Schedule';

const App = () => {
  const [tasks] = createResource(async () => {
    return await TaskService.GetTasks();
  });

  return (
    <div class="mt-[40px] pt-4 m-4">
      <div class="fixed top-0 left-0 px-[80px] h-[40px] w-full bg-primary text-primary-content shadow-md">
        <TitleBar title="🐪 Camel Do 🐫" />
      </div>

      <main class="h-[calc(100dvh-40px-2rem)] w-full overflow-y-auto">
        <ErrorBoundary fallback={(err) => <div>Error: {err.message}</div>}>
          <Schedule tasks={tasks()} />
        </ErrorBoundary>
      </main>
    </div>
  );
};

export default App;