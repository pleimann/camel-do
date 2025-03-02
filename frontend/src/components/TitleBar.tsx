import { TbSearch as SearchIcon, TbBell as BellIcon, TbMoon as MoonIcon, TbSun as SunIcon } from 'solid-icons/tb'

export default function TitleBar() {
  return (
    <div class="navbar bg-neutral-100 shadow-md fixed top-0 left-0 w-full">
      <div class="navbar-start">
      </div>
      <div class="navbar-center">
        <a class="btn btn-ghost text-xl">Camel Do</a>
      </div>
      <div class="navbar-end">
        <button class="btn btn-ghost btn-circle">
          <SearchIcon class="size-6" />
        </button>
        <button class="btn btn-ghost btn-circle">
          <div class="indicator">
            <BellIcon class="size-6" />
            <span class="badge badge-xs badge-primary indicator-item"></span>
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
