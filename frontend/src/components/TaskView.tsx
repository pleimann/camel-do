import { createSignal } from "solid-js";
import { BsCheck2Circle as CheckCircleIcon } from 'solid-icons/bs';
import { TbMenu as MenuIcon } from 'solid-icons/tb'
import { RiArrowsArrowDownSLine as ChevronDownIcon, RiArrowsArrowUpSLine as ChevronUpIcon } from 'solid-icons/ri'
import { AiOutlineSchedule as ScheduleIcon } from 'solid-icons/ai'
import { TbTrashX as TrashIcon } from 'solid-icons/tb'

import { Task } from "@bindings/pleimann.com/camel-do/services";

interface Props {
  task: Task;
}

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

export default function TaskView(props: Props) {
  const p = props;

  const [show, toggleShow] = createSignal(false);

  const backgroundClass = backgroundColors[p.task.color];
  const textClass = textColors[p.task.color];
  
  let actionsMenu!: HTMLDivElement;

  const completeTask = () => {
    console.log("complete")
    actionsMenu.removeAttribute("open");
  };

  const scheduleTask = () => {
    console.log("schedule")
    actionsMenu.removeAttribute("open");
  };

  const deleteTask = () => {
    console.log("delete")
    actionsMenu.removeAttribute("open");
  };

  return (
    <div class="card card-side card-sm bg-base-100 shadow-sm select-none">
      <figure class={`w-20 max-w-20 min-w-20 ${backgroundClass} ${textClass}`}>
        <CheckCircleIcon class="size-10" />
      </figure>
      <div class="card-body flex-col justify-start items-start">
        <p class="card-title">{p.task.title}</p>
        <p class="text-sm">{p.task.duration}</p>
      </div>
      <div class="card-actions rounded-e-box p-2 bg-base-200">
        <div class="flex flex-col">
          <div class="dropdown dropdown-hover dropdown-left dropdown-center" ref={actionsMenu}>
            <button tabindex="0" class="btn btn-circle btn-ghost"><MenuIcon class="size-6" /></button>
            <ul tabindex="0" class="dropdown-content pr-4 z-1 gap-2 flex flex-row-reverse">
              <li><button class="btn btn-circle tooltip tooltip-bottom shadow-sm" data-tip="Complete" onClick={completeTask}><CheckCircleIcon class="size-6" /></button></li>
              <li><button class="btn btn-circle tooltip tooltip-bottom shadow-sm" data-tip="Schedule" onClick={scheduleTask}><ScheduleIcon class="size-6" /></button></li>
              <li><button class="btn btn-circle tooltip tooltip-bottom shadow-sm" data-tip="Delete" onClick={deleteTask}><TrashIcon class="size-6" /></button></li>
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