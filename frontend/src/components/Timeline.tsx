import { Task } from "@bindings/pleimann.com/camel-do/model";
import { For, Show } from "solid-js";
import Icon from "@/components/Icon";

interface Props {
  tasks: Task[]
}

export default function Timeline(props: Props) {
  const p = props;

  return (
    <ul class="timeline timeline-snap-icon max-md:timeline-compact timeline-vertical">
      <For each={p.tasks} fallback={<div>Loading...</div>}>
        {(task, i) => (
          <li>
            <Show when={i() < 0}>
              <hr />
            </Show>
            <div class="timeline-middle mb-2">
              {task.completed ? (<Icon.CircleChecked class="size-6" />) : (<Icon.Circle class="size-6" />)}
            </div>
            <div class={`pt-1 mb-4 ${i() % 2 ? "md:text-end timeline-start" : "md:text-start timeline-end"} `} >
              <div class="font-bold">{task.title}</div>
              <div class="font-light italic">{task.duration}</div>
              {task.description}
            </div>
            <Show when={i() < p.tasks.length - 1}>
              <hr />
            </Show>
          </li>
        )}
      </For>
    </ul>
  );
}
