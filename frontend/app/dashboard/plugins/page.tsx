import { PluginTable } from "@/components/PluginTable"

export default async function PluginsPage() {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/plugins`, {
    next: { revalidate: 30 },
  } as any)
  const plugins = (await res.json()) as Array<{
    id: string
    name: string
    platform: string
    version: string
    latest_version: string
    last_updated: string
  }>

  return (
    <section className="max-w-6xl mx-auto p-6">
      <h1 className="text-2xl font-bold mb-4">Plugin Discovery & Health Matrix</h1>
      <PluginTable plugins={plugins} />
    </section>
  )
}
