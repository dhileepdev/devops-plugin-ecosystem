import { clsx } from "clsx"
import { twMerge } from "tailwind-merge"

function cn(...classes: string[]) {
  return twMerge(clsx(classes))
}

interface Plugin {
  id: string
  name: string
  platform: string
  version: string
  latest_version: string
  last_updated: string
}

interface PluginTableProps {
  plugins: Plugin[]
}

export function PluginTable({ plugins }: PluginTableProps) {
  return (
    <div className="overflow-x-auto">
      <table className="w-full text-left border-collapse">
        <thead>
          <tr className="border-b border-slate-200">
            <th className="pb-3 font-medium text-slate-600">ID</th>
            <th className="pb-3 font-medium text-slate-600">Name</th>
            <th className="pb-3 font-medium text-slate-600">Platform</th>
            <th className="pb-3 font-medium text-slate-600">Version</th>
            <th className="pb-3 font-medium text-slate-600">Latest</th>
            <th className="pb-3 font-medium text-slate-600">Last updated</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-slate-100">
          {plugins.map((p) => (
            <tr key={p.id} className="hover:bg-slate-50">
              <td className="py-3 font-mono text-sm">{p.id}</td>
              <td className="py-3 text-sm">{p.name}</td>
              <td className="py-3 text-sm">
                <span className={cn("inline-block px-2 py-0.5 rounded text-xs font-medium", p.platform === "terraform" ? "bg-amber-100 text-amber-800" : "bg-slate-100 text-slate-600")}>
                  {p.platform}
                </span>
              </td>
              <td className="py-3 text-sm font-mono">{p.version}</td>
              <td className="py-3 text-sm">{p.latest_version}</td>
              <td className="py-3 text-sm text-slate-500">{p.last_updated}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
