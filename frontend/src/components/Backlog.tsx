import { Task } from "@bindings/pleimann.com/camel-do/services";
import { For, mergeProps } from "solid-js";
import TaskView from "./TaskView";

interface Props {
  tasks?: Task[]
}

export default function Backlog(props: Props) {
  const p = mergeProps({tasks: []}, props);
 
  return (
    <div class="py-4 pl-4 bg-primary-content w-96 h-full">
      <div class="pr-4 mb-4 flex flex-col gap-2 h-full overflow-y-auto ">
        <For each={p.tasks} fallback={<div>Loading...</div>}>
          {(task) => <TaskView task={task} />}
        </For>
      </div>
    </div>
  )
}