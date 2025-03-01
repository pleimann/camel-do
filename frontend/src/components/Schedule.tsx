import { Task } from "@bindings/pleimann.com/camel-do/services";
import { For, mergeProps } from "solid-js";

interface Props {
  tasks?: Task[]
}

export default function Schedule(props: Props) {
  const p = mergeProps({tasks: []}, props);
 
  return (
    <div class={"flex flex-col gap-2"}>
      <For each={p.tasks} fallback={<div>Loading...</div>}>
        {(task) => (
          <label style="text-wrap: nowrap;">
            <input type="checkbox" class="checkbox checkbox-primary" checked={task.completed} />
            <span class="ml-2">{task.title}</span>
          </label>
        )}
      </For>
    </div>
  )
}