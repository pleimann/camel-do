import { SearchIcon, BellIcon, MoonIcon, SunIcon, RefreshIcon } from '@/components/Icons';

export type TitleBarAction = "refresh" | "search";

interface Props {
  onAction: (action: TitleBarAction) => void
}

export default function TitleBar(props: Props) {
  const p = props;

  const refresh = (e: MouseEvent) => {
    p.onAction("refresh");
  }

  const search = (e: MouseEvent) => {
    p.onAction("search");
  }

  return (
    <div class="navbar bg-neutral-100 shadow-md fixed top-0 left-0 w-full">
      <div class="navbar-start">
      </div>
      <div class="navbar-center">
        <a class="btn btn-ghost text-xl">Camel Do</a>
      </div>
      <div class="navbar-end flex flex-row gap-2">
        <button class="btn btn-ghost btn-circle" onClick={search}>
          <SearchIcon class="size-6" />
        </button>
        <button class="btn btn-ghost btn-circle">
          <BellIcon class="size-6" />
        </button>
        <button class="btn btn-ghost btn-circle" onClick={refresh}>
          <RefreshIcon class="size-6" />
        </button>
        <button class="btn btn-ghost btn-circle swap swap-rotate">
          <SunIcon class="size-6 swap-on" />
          <MoonIcon class="size-6 swap-off" />
        </button>
      </div>
    </div>
  )
}
