import { clsx } from "clsx"
import { twMerge } from "tailwind-merge"

function cn(...classes: string[]) {
  return twMerge(clsx(classes))
}

interface HealthBadgeProps {
  status: "green" | "yellow" | "red" | string
}

const styles: Record<string, string> = {
  green: "bg-emerald-100 text-emerald-800",
  yellow: "bg-amber-100 text-amber-800",
  red: "bg-red-100 text-red-800",
}

export function HealthBadge({ status }: HealthBadgeProps) {
  return (
    <span className={cn("inline-block px-2 py-0.5 rounded text-xs font-medium uppercase", styles[status] || "bg-slate-100 text-slate-600")}>
      {status}
    </span>
  )
}
