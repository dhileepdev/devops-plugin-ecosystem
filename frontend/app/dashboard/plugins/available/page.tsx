export default async function AvailablePluginsPage() {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/manifests`, {
    next: { revalidate: 60 },
  } as any)
  const manifests = (await res.json()) as Array<{
    name: string
    version: string
    description: string
    category: string
    platforms: string[]
    cli_commands: string[]
    standalone_usage: string
  }>

  return (
    <section className="max-w-6xl mx-auto p-6">
      <h1 className="text-2xl font-bold mb-4">Available Plugins</h1>
      <p className="text-slate-600 mb-6">Each plugin can be used individually via CLI or API.</p>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {manifests.map((m) => (
          <div key={m.name} className="border border-slate-200 rounded p-4">
            <div className="flex items-center justify-between mb-2">
              <h2 className="font-semibold">{m.name}</h2>
              <span className="text-xs bg-slate-100 px-2 py-0.5 rounded">{m.version}</span>
            </div>
            <p className="text-sm text-slate-600 mb-2">{m.description}</p>
            <div className="text-xs text-slate-500 mb-1">
              Category: {m.category} · Platforms: {m.platforms.join(", ")}
            </div>
            {m.standalone_usage && (
              <code className="block text-xs bg-slate-50 p-2 rounded mt-2">{m.standalone_usage}</code>
            )}
          </div>
        ))}
      </div>
    </section>
  )
}
