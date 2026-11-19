"use client";

import { useId } from "react";

export interface SegmentOption<T extends string> {
  value: T;
  label: string;
  title?: string;
}

interface SegmentedControlProps<T extends string> {
  label: string;
  showLabel?: boolean;
  options: readonly SegmentOption<T>[];
  value: T;
  onChange: (value: T) => void;
  variant?: "default" | "large" | "tabs";
}

export function SegmentedControl<T extends string>({
  label,
  showLabel = true,
  options,
  value,
  onChange,
  variant = "default",
}: SegmentedControlProps<T>) {
  const id = useId();
  const labelId = `${id}-label`;

  return (
    <div className="field">
      <span className={showLabel ? "label" : "visually-hidden"} id={labelId}>
        {label}
      </span>
      <div className={`segmented ${variant}`} role="radiogroup" aria-labelledby={labelId}>
        {options.map((option) => (
          <label key={option.value} className="segment" title={option.title}>
            <input
              type="radio"
              className="visually-hidden"
              name={id}
              value={option.value}
              checked={value === option.value}
              onChange={() => onChange(option.value)}
            />
            <span>{option.label}</span>
          </label>
        ))}
      </div>
    </div>
  );
}
