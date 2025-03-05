import Icon from "@/components/Icon";

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
    <div class="navbar shadow-lg bg-auto w-full p-0 pr-6">
      <div class="navbar-start">
      </div>
      <div class="navbar-center">
        <a class="text-lg font-black font-stretch-125%">Camel Do</a>
      </div>
      <div class="navbar-end flex flex-row gap-2">
        <button class="btn btn-lg btn-ghost btn-circle" onClick={search}>
          <Icon.Search class="size-6" />
        </button>
        <button class="btn btn-lg btn-ghost btn-circle">
          <Icon.Bell class="size-6" />
        </button>
        <button class="btn btn-lg btn-ghost btn-circle" onClick={refresh}>
          <Icon.Refresh class="size-6" />
        </button>
        <label class="btn btn-lg btn-ghost btn-circle swap swap-rotate">
          <input type="checkbox" value="dark" class="theme-controller" />
          <Icon.Sun class="size-6 swap-on" />
          <Icon.Moon class="size-6 swap-off" />
        </label>
      </div>
    </div>
  )
}
