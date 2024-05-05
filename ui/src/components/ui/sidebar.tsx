import { Button } from "@/components/ui/button";

export function Sidebar() {
  return (
    <div className="grid grid-rows-2 auto-rows-auto">
      <div className="bg-gray-100 p-5">
        <h2 className="text-xl font-bold">Sentinel Manager</h2>
      </div>
      <div className="bg-gray-300 grid grid-flow-row gap-1">
        <Button variant="ghost">Dashboard</Button>
        <Button variant="ghost">Rules</Button>
        <Button variant="ghost">Alerts</Button>
        <Button variant="ghost">Settings</Button>
      </div>
    </div>
  )
}