
interface Props {
  title: string;
}

export default function TitleBar(props: Props) {
  return (
    <div class="h-full w-full flex items-center justify-center content-center">
      <div class="h-fit font-bold text-xl">{props.title}</div>
    </div>
  )
}
