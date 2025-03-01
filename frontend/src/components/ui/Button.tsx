import { mergeProps } from "solid-js";

interface Props {
  label: string;
  color?: "primary" | "accent" | "neutral";
  type?: "button" | "submit" | "reset";
  style?: "outline" | "soft" | "ghost" | "link"
  active?: boolean;
  disabled?: boolean;
  size?: "xs" | "sm" | "md" | "lg" | "xl";
  wide?: boolean;
  block?: boolean;
  square?: boolean;
  circle?: boolean;
  onClick?: () => void;
}

const defaults: Props = {
  label: "Button",
  color: "primary",
  type: "button",
  disabled: false,
  onClick: () => console.log("Button clicked"),
};

export default function Button(props: Props) {
  const p = mergeProps(defaults, props);

  return (
    <button type={p.type} disabled={p.disabled} onClick={p.onClick} 
      class={`
        btn 
        btn-${p.color} 
        btn-${p.style} 
        btn-${p.size} 
        ${p.wide ? "btn-wide" : ""} 
        ${p.block ? "btn-block" : ""} 
        ${p.square ? "btn-square" : ""} 
        ${p.circle ? "btn-circle" : ""}
      `}
    >{props.label}</button>
  );
}