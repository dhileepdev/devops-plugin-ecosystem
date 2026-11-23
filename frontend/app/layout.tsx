import "./globals.css"

export const metadata = {
  title: "DevOps Plugin Ecosystem",
  description: "Unified plugin discovery & health management platform",
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-white text-slate-900">
        <nav className="border-b border-slate-200 px-6 py-3 flex gap-4 text-sm">
          <a href="/" className="font-semibold">DevOps Plugins</a>
          <a href="/dashboard/plugins">Plugins</a>
          <a href="/dashboard/plugins/available">Available Plugins</a>
          <a href="/dashboard/terraform">Terraform Health</a>
          <a href="/dashboard/operator">Operator RBAC</a>
          <a href="/dashboard/notifications">Notifications</a>
        </nav>
        {children}
      </body>
    </html>
  )
}
