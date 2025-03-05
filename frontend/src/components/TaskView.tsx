import { mergeProps } from "solid-js";

import Icon from '@/components/Icon';

import { Task } from "@bindings/pleimann.com/camel-do/model";
import { type TaskAction } from "@/components/Backlog";

[
  "bg-red-100",
  "bg-orange-100",
  "bg-amber-100",
  "bg-yellow-100",
  "bg-lime-100",
  "bg-green-100",
  "bg-emerald-100",
  "bg-teal-100",
  "bg-cyan-100",
  "bg-sky-100",
  "bg-violet-100",
  "bg-purple-100",
  "bg-fuchsia-100",
  "bg-pink-100",
  "bg-rose-100",
];

[
  "text-red-900",
  "text-orange-900",
  "text-amber-900",
  "text-yellow-900",
  "text-lime-900",
  "text-green-900",
  "text-emerald-900",
  "text-teal-900",
  "text-cyan-900",
  "text-sky-900",
  "text-violet-900",
  "text-purple-900",
  "text-fuchsia-900",
  "text-pink-900",
  "text-rose-900",
];

interface Props {
  task: Task;
  onTaskAction: (task: Task, action: TaskAction) => void
}

const defaults: Partial<Props> = {
  onTaskAction: (task, action) => {}
}

export default function TaskView(props: Props) {
  const p = mergeProps(defaults, props);

  const backgroundClass = `bg-${p.task.color.toLowerCase()}-100`;
  const textClass = `text-${p.task.color.toLowerCase()}-900`;
  
  let actionsMenu!: HTMLDivElement;

  const taskAction = (action: TaskAction) => {
    p.onTaskAction(p.task, action);
  };

  return (
    <div class="card card-side card-xs bg-base-100 dark:bg-base-300 shadow-md text-sm select-none rounded-2xl">
      <figure class={`w-18 ${backgroundClass} ${textClass}`}>
        <Icon name={p.task.icon} class="size-7" strokeWidth={1.25} />
      </figure>
      <div class="card-body flex-col justify-center items-start">
        <div class="text-sm font-medium">{p.task.title}</div>
        <time class="italic">{p.task.duration}</time>
      </div>
      <div class="card-actions rounded-e-2xl bg-base-200 p-2">
        <div class="flex flex-col justify-between">
          <div class="dropdown dropdown-hover dropdown-left dropdown-center" ref={actionsMenu}>
            <button class="btn btn-circle btn-ghost" role="button" tabIndex={0}><Icon.Menu class="size-4" /></button>
            <ul class="dropdown-content p-2 z-1 gap-2 flex flex-row-reverse rounded-s-full bg-base-200/90 bg-blend-overlay" tabIndex={0}>
              <li><button class="btn btn-circle btn-ghost tooltip tooltip-bottom" data-tip="Edit" onClick={(e) => taskAction('edit')}><Icon.Edit class="size-4" /></button></li>
              <li>
                <button class="btn btn-circle btn-ghost tooltip tooltip-bottom" data-tip={p.task.completed ? "Uncomplete" : "Complete"} onClick={(e) => taskAction('complete')}>
                  {p.task.completed ? (<Icon.CircleChecked class="size-4" />) : (<Icon.Circle class="size-4" />)}
                </button>
              </li>
              <li><button class="btn btn-circle btn-ghost tooltip tooltip-bottom" data-tip="Delete" onClick={(e) => taskAction('delete')}><Icon.Trash class="size-4" /></button></li>
            </ul>
          </div>
          <button class="btn btn-circle btn-ghost tooltip tooltip-left" data-tip="Schedule" onClick={(e) => taskAction('schedule')}><Icon.Schedule class="size-4" /></button>
        </div>
      </div>
    </div>
  );
}
