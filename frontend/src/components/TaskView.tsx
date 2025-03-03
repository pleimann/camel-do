import { createSignal, mergeProps } from "solid-js";

import { 
  ChevronDownIcon, 
  ChevronUpIcon, 
  ScheduleIcon, 
  TrashIcon,
  MenuIcon,
  CheckCircleIcon,
} from '@/components/Icons';

import { Task } from "@bindings/pleimann.com/camel-do/services";
import { type TaskAction } from "@/components/Backlog";

const backgroundColors = [
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
  "bg-blue-100",
  "bg-indigo-100",
  "bg-violat-100",
  "bg-purple-100",
  "bg-fuchsio-100",
  "bg-pink-100",
  "bg-rose-100",
  "bg-slate-100",
  "bg-gray-100",
  "bg-zinc-100",
  "bg-neutral-100",
  "bg-stone-100",
]

const textColors = [
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
  "text-blue-900",
  "text-indigo-900",
  "text-violat-900",
  "text-purple-900",
  "text-fuchsio-900",
  "text-pink-900",
  "text-rose-900",
  "text-slate-900",
  "text-gray-900",
  "text-zinc-900",
  "text-neutral-900",
  "text-stone-900",
]

interface Props {
  task: Task;
  onTaskAction: (task: Task, action: TaskAction) => void
}

const defaults: Partial<Props> = {
  onTaskAction: (task, action) => {}
}

export default function TaskView(props: Props) {
  const p = mergeProps(defaults, props);

  const [show, toggleShow] = createSignal(false);

  const backgroundClass = backgroundColors[p.task.color];
  const textClass = textColors[p.task.color];
  
  let actionsMenu!: HTMLDivElement;

  const taskAction = (action: TaskAction) => {
    p.onTaskAction(p.task, action);
    
    close();
  };

  const close = () => (document.activeElement as HTMLElement)?.blur();

  return (
    <div class="card card-side card-sm bg-base-100 shadow-sm select-none">
      <figure class={`w-20 max-w-20 min-w-20 ${backgroundClass} ${textClass}`}>
        <CheckCircleIcon class="size-10" />
      </figure>
      <div class="card-body flex-col justify-start items-start">
        <p class="card-title">{p.task.title}</p>
        <p class="text-sm">{p.task.duration}</p>
      </div>
      <div class="card-actions rounded-e-box bg-base-200 p-2">
        <div class="flex flex-col gap-2">
          <div class="dropdown dropdown-hover dropdown-left dropdown-center rounded-2xl" ref={actionsMenu}>
            <button class="btn btn-circle btn-ghost" role="button" tabIndex={0}><MenuIcon class="size-6" /></button>
            <ul class="dropdown-content p-2 z-1 gap-2 flex flex-row-reverse rounded-s-full bg-base-200/90 bg-blend-overlay" tabIndex={0}>
              <li><button class="btn btn-circle tooltip tooltip-bottom shadow-sm" data-tip="Complete" onClick={(e) => taskAction('complete')}><CheckCircleIcon class="size-6" /></button></li>
              <li><button class="btn btn-circle tooltip tooltip-bottom shadow-sm" data-tip="Schedule" onClick={(e) => taskAction('schedule')}><ScheduleIcon class="size-6" /></button></li>
              <li><button class="btn btn-circle tooltip tooltip-bottom shadow-sm" data-tip="Delete" onClick={(e) => taskAction('delete')}><TrashIcon class="size-6" /></button></li>
            </ul>
          </div>
          <button class="btn btn-circle btn-ghost" onClick={() => toggleShow(!show())}>
            {show() ? (<ChevronUpIcon class="size-6" />) : (<ChevronDownIcon class="size-6" />)}
          </button>
        </div>
      </div>
    </div>
  );
}