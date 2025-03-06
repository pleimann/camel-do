import { Task } from "@bindings/pleimann.com/camel-do/model";
import { For, Show } from "solid-js";
import { format } from 'date-fns';
import Icon from "@/components/Icon";

interface Props {
  tasks: Task[]
}

export default function Timeline(props: Props) {
  const p = props;

  return (
    <ul class="timeline timeline-vertical timeline-compact timeline-snap-icon">
      <For each={p.tasks} fallback={<div>Loading...</div>}>
        {(task, i) => (
          <li>
            <Show when={i() < 0}><hr /></Show>
            <div class="timeline-start">
              <time class="font-medium" dateTime={task.startTime}>{format(task.startTime, "HH:mm")}</time>
              <div class="font-light italic">{task.duration}</div>
            </div>
            <div class="timeline-middle mb-2">
              {task.completed ? (<Icon.CircleChecked class="size-8" />) : (<Icon.Circle class="size-8" />)}
            </div>
            <div class="pt-1 mb-4 timeline-end">
              <div class="font-bold">{task.title}</div>
              {task.description}
            </div>
            <Show when={i() < p.tasks.length - 1}><hr /></Show>
          </li>
        )}
      </For>
    </ul>
  );
}
