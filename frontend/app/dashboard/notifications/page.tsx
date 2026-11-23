import { NotificationFeed } from "@/components/NotificationFeed"

export const dynamic = "force-dynamic"

export default async function NotificationsPage() {
  let notifications: any[] = []
  try {
    const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/notifications`, {
      next: { revalidate: 10 },
    } as any)
    if (res.ok) notifications = (await res.json()) as any[]
  } catch {
    notifications = []
  }

  return (
    <section className="max-w-6xl mx-auto p-6">
      <h1 className="text-2xl font-bold mb-4">Smart Notification Deduplicator</h1>
      <NotificationFeed notifications={notifications} />
    </section>
  )
}
