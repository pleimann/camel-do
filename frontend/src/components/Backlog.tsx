import { Task } from "@bindings/pleimann.com/camel-do/model";
import { For, mergeProps } from "solid-js";
import TaskView from "./TaskView";

interface Props {
  tasks: Task[]
}

export type TaskAction = "complete" | "schedule" | "delete" | "edit";

export default function Backlog(props: Props) {
  const p = mergeProps({tasks: []}, props);
 
  const onTaskAction = (task: Task, action: TaskAction) => {
    console.log(`${action} task ${task.id}`)
  }

  return ( // padding being defined in Backlog component rather than parent (App) prevents drop shadows from being clipped
    <div class="pl-4 pb-2 pr-1 flex flex-col gap-2 h-full w-full overflow-y-scroll">
      <For each={p.tasks} fallback={<div>Loading...</div>}>
        {(task) => <TaskView task={task} onTaskAction={onTaskAction} />}
      </For>
    </div>
  )
}