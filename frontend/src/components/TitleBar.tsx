import { SearchIcon, BellIcon, MoonIcon, SunIcon } from '@/components/Icons';

export default function TitleBar() {
  return (
    <div class="navbar bg-neutral-100 shadow-md fixed top-0 left-0 w-full">
      <div class="navbar-start">
      </div>
      <div class="navbar-center">
        <a class="btn btn-ghost text-xl">Camel Do</a>
      </div>
      <div class="navbar-end flex flex-row gap-2">
        <button class="btn btn-ghost btn-circle">
          <SearchIcon class="size-6" />
        </button>
        <button class="btn btn-ghost btn-circle">
          <div class="indicator">
            <BellIcon class="size-6" />
          </div>
        </button>
        <button class="btn btn-ghost btn-circle">
          <div class="swap swap-rotate">
            <SunIcon class="size-6 swap-on" />
            <MoonIcon class="size-6 swap-off" />
          </div>
        </button>
      </div>
    </div>
  )
}
