export default function HomePage() {
  return (
    <section className="max-w-6xl mx-auto p-6">
      <h1 className="text-3xl font-bold mb-2">DevOps Plugin Ecosystem</h1>
      <p className="text-slate-600 mb-6">
        Unified plugin discovery & health management platform. Choose a module:
      </p>
      <ul className="space-y-2">
        <li><a className="text-blue-600 hover:underline" href="/dashboard/plugins">Plugin Discovery & Health Matrix</a></li>
        <li><a className="text-blue-600 hover:underline" href="/dashboard/terraform">Terraform Provider Health Dashboard</a></li>
        <li><a className="text-blue-600 hover:underline" href="/dashboard/operator">Operator Lifecycle & RBAC Analyzer</a></li>
        <li><a className="text-blue-600 hover:underline" href="/dashboard/notifications">Smart Notification Deduplicator</a></li>
      </ul>
    </section>
  )
}
