import { PluginTable } from "@/components/PluginTable"

export const dynamic = "force-dynamic"

export default async function PluginsPage() {
  let plugins: Array<{
    id: string
    name: string
    platform: string
    version: string
    latest_version: string
    last_updated: string
  }> = []
  try {
    const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/plugins`, {
      next: { revalidate: 30 },
    } as any)
    if (res.ok) plugins = (await res.json()) as typeof plugins
  } catch {
    plugins = []
  }

  return (
    <section className="max-w-6xl mx-auto p-6">
      <h1 className="text-2xl font-bold mb-4">Plugin Discovery & Health Matrix</h1>
      <PluginTable plugins={plugins} />
    </section>
  )
}
