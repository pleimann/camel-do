import { Task } from "@bindings/pleimann.com/camel-do/model";
import { For, mergeProps } from "solid-js";
import TaskView from "./TaskView";

interface Props {
  tasks: Task[]
}

export type TaskAction = "complete" | "schedule" | "delete";

export default function Backlog(props: Props) {
  const p = mergeProps({tasks: []}, props);
 
  const onTaskAction = (task: Task, action: TaskAction) => {
    console.log(`${action} task ${task.id}`)
  }

  return (
    <div class="w-96 min-w-96 mb-4 mr-1 pr-1 flex flex-col gap-2 h-full overflow-y-auto">
      <For each={p.tasks} fallback={<div>Loading...</div>}>
        {(task) => <TaskView task={task} onTaskAction={onTaskAction} />}
      </For>
    </div>
  )
}