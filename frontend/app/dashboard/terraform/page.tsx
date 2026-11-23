import { HealthBadge } from "@/components/HealthBadge"

interface PerfIssue {
  number?: number
  title?: string
}

interface HealthReport {
  provider: string
  status: "green" | "yellow" | "red"
  open_critical_bugs: number
  performance_regressions: string[]
  plan_stability: number
  open_issues: number
}

export default async function TerraformPage() {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/health/terraform`, {
    next: { revalidate: 60 },
  } as any)
  const report = (await res.json()) as HealthReport

  return (
    <section className="max-w-6xl mx-auto p-6">
      <h1 className="text-2xl font-bold mb-4">Terraform Provider Health Dashboard</h1>
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <div className="border border-slate-200 rounded p-4">
          <div className="text-sm text-slate-500">Provider Status</div>
          <div className="mt-2"><HealthBadge status={report.status} /></div>
        </div>
        <div className="border border-slate-200 rounded p-4">
          <div className="text-sm text-slate-500">Open Critical Bugs</div>
          <div className="text-2xl font-bold mt-2">{report.open_critical_bugs}</div>
        </div>
        <div className="border border-slate-200 rounded p-4">
          <div className="text-sm text-slate-500">Plan Stability</div>
          <div className="text-2xl font-bold mt-2">{report.plan_stability}%</div>
        </div>
        <div className="border border-slate-200 rounded p-4">
          <div className="text-sm text-slate-500">Open Issues</div>
          <div className="text-2xl font-bold mt-2">{report.open_issues}</div>
        </div>
      </div>
      {report.performance_regressions && report.performance_regressions.length > 0 && (
        <div className="border border-slate-200 rounded p-4">
          <h2 className="font-semibold mb-2">Performance Regressions</h2>
          <ul className="list-disc pl-5 text-sm">
            {report.performance_regressions.map((p, i) => (
              <li key={i}>{p}</li>
            ))}
          </ul>
        </div>
      )}
    </section>
  )
}
