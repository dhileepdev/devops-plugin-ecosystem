"use client"

import { useState } from "react"

interface ScanResult {
  scan: {
    kubectl_version: string
    api_version: string
    upgrade_policy: string
    rbac_roles: string[]
    has_port_forward: boolean
    port_forward_config: string
  }
  upgrade_path: {
    current_version: string
    target_version: string
    steps: string[]
    estimated_downtime: string
  }
  rbac: {
    gaps: { role: string; missing?: string[]; overly_broad?: string[] }[]
  }
}

export default function OperatorPage() {
  const [manifest, setManifest] = useState("")
  const [target, setTarget] = useState("latest")
  const [result, setResult] = useState<ScanResult | null>(null)
  const [loading, setLoading] = useState(false)

  async function analyze() {
    setLoading(true)
    const res = await fetch("/api/operator/analyze", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ manifest, target_version: target }),
    })
    const data = await res.json()
    setResult(data)
    setLoading(false)
  }

  return (
    <section className="max-w-6xl mx-auto p-6">
      <h1 className="text-2xl font-bold mb-4">Operator Lifecycle & RBAC Analyzer</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
        <div>
          <label className="block text-sm font-medium mb-1">Manifest (YAML)</label>
          <textarea
            value={manifest}
            onChange={(e) => setManifest(e.target.value)}
            className="w-full h-64 border border-slate-200 rounded p-2 font-mono text-xs"
            placeholder="paste operator manifest here..."
          />
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">Target version</label>
          <input
            value={target}
            onChange={(e) => setTarget(e.target.value)}
            className="w-full border border-slate-200 rounded p-2 text-sm mb-4"
          />
          <button
            onClick={analyze}
            disabled={loading || !manifest}
            className="px-4 py-2 bg-slate-900 text-white rounded text-sm disabled:opacity-50"
          >
            {loading ? "Analyzing..." : "Analyze"}
          </button>
        </div>
      </div>
      {result && (
        <div className="space-y-4">
          <div className="border border-slate-200 rounded p-4">
            <h2 className="font-semibold mb-2">Scan</h2>
            <pre className="text-xs bg-slate-50 p-2 rounded overflow-x-auto">{JSON.stringify(result.scan, null, 2)}</pre>
          </div>
          <div className="border border-slate-200 rounded p-4">
            <h2 className="font-semibold mb-2">Upgrade Path</h2>
            <ul className="list-disc pl-5 text-sm">
              {result.upgrade_path.steps.map((s, i) => (
                <li key={i}>{s}</li>
              ))}
            </ul>
            <p className="text-xs text-slate-500 mt-2">Estimated downtime: {result.upgrade_path.estimated_downtime}</p>
          </div>
          <div className="border border-slate-200 rounded p-4">
            <h2 className="font-semibold mb-2">RBAC Gaps</h2>
            {result.rbac.gaps.length === 0 ? (
              <p className="text-sm text-emerald-600">No gaps detected.</p>
            ) : (
              <ul className="space-y-2 text-sm">
                {result.rbac.gaps.map((g, i) => (
                  <li key={i} className="border-l-2 border-red-300 pl-2">
                    <strong>{g.role}</strong>
                    {g.missing && g.missing.length > 0 && <div className="text-red-600">Missing: {g.missing.join(", ")}</div>}
                    {g.overly_broad && g.overly_broad.length > 0 && <div className="text-amber-600">Overly broad: {g.overly_broad.join(", ")}</div>}
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>
      )}
    </section>
  )
}
