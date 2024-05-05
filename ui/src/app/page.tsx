import { Button } from "@/components/ui/button"
import { Sidebar } from "@/components/ui/sidebar"

export default function Home() {
  return (
    <div className="grid grid-cols-6">
      <div className="col-span-1 bg-blue-300">
        <Sidebar />
      </div>
      <div className="col-start-2 col-end-7 bg-slate-600">
        <Button>Click me</Button>
      </div>
    </div>
  )
}

