import { TaskService } from '@bindings/pleimann.com/camel-do/services'

import { createResource, For, ErrorBoundary } from 'solid-js';

import TitleBar from '@/components/TitleBar';
import Button from '@/components/ui/Button';

const App = () => {
  const [todos] = createResource(async () => {
    return await TaskService.GetTasks();
  });

  return (
    <div class="mt-[40px] pt-4 m-4">
      <div class="fixed top-0 left-0 pl-[80px] h-[40px] w-full bg-primary text-primary-content shadow-md">
        <TitleBar title="🐪 Camel Do 🐫" />
      </div>

      <main class="h-[calc(100dvh-40px-2rem)] w-full overflow-y-auto">
        <ErrorBoundary fallback={(err) => <div>Error: {err.message}</div>}>
          <Button label="Open" />
          <For each={todos()} fallback={<div>Loading...</div>}>
            {(todo) => <div>{todo.title}</div>}
          </For>
        </ErrorBoundary>
      </main>
    </div>
  );
};

export default App;