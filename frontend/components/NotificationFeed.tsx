interface Notification {
  id: string
  event_type: string
  source: string
  severity: string
  message: string
  timestamp: string
  deduplicated?: boolean
  routed_to?: string[]
}

interface NotificationFeedProps {
  notifications: Notification[]
}

export function NotificationFeed({ notifications }: NotificationFeedProps) {
  if (!notifications || notifications.length === 0) {
    return <p className="text-slate-500 text-sm">No notifications yet.</p>
  }
  return (
    <ul className="space-y-3">
      {notifications.map((n) => (
        <li key={n.id} className="border border-slate-200 rounded p-3">
          <div className="flex items-center justify-between">
            <span className="font-medium text-sm">{n.event_type}</span>
            <span className="text-xs text-slate-500">{n.severity}</span>
          </div>
          <p className="text-sm text-slate-700 mt-1">{n.message}</p>
          <div className="text-xs text-slate-400 mt-1">
            {n.source} · {n.timestamp}
            {n.deduplicated && <span className="ml-2 text-amber-600">(deduplicated)</span>}
            {n.routed_to && n.routed_to.length > 0 && (
              <span className="ml-2">→ {n.routed_to.join(", ")}</span>
            )}
          </div>
        </li>
      ))}
    </ul>
  )
}
