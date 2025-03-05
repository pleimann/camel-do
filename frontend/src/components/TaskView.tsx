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
    <div class="card card-side card-sm bg-base-100 shadow-md select-none">
      <figure class={`w-20 max-w-20 min-w-20 ${backgroundClass} ${textClass}`}>
        <Icon name={p.task.icon} class="size-8" />
      </figure>
      <div class="card-body flex-col justify-start items-start">
        <p class="text-lg font-bold line-clamp-1">{p.task.title}</p>
        <time class="text-sm italic">{p.task.duration}</time>
      </div>
      <div class="card-actions rounded-e-box bg-base-200 p-2">
        <div class="flex flex-col gap-2">
          <div class="dropdown dropdown-hover dropdown-left dropdown-center rounded-2xl" ref={actionsMenu}>
            <button class="btn btn-circle btn-ghost" role="button" tabIndex={0}><Icon.Menu class="size-6" /></button>
            <ul class="dropdown-content p-2 z-1 gap-2 flex flex-row-reverse rounded-s-full bg-base-200/90 bg-blend-overlay" tabIndex={0}>
              <li><button class="btn btn-circle tooltip tooltip-bottom shadow-sm" data-tip="Complete" onClick={(e) => taskAction('complete')}><Icon.ChevronDown class="size-6" /></button></li>
              <li><button class="btn btn-circle tooltip tooltip-bottom shadow-sm" data-tip="Schedule" onClick={(e) => taskAction('schedule')}><Icon.Schedule class="size-6" /></button></li>
              <li><button class="btn btn-circle tooltip tooltip-bottom shadow-sm" data-tip="Delete" onClick={(e) => taskAction('delete')}><Icon.Trash class="size-6" /></button></li>
            </ul>
          </div>
          <button class="btn btn-circle btn-ghost" data-tip={p.task.completed ? "Uncomplete" : "Complete"} onClick={(e) => taskAction('complete')}>
            {p.task.completed ? (<Icon.CircleChecked class="size-6" />) : (<Icon.Circle class="size-6" />)}
          </button>
        </div>
      </div>
    </div>
  );
}
