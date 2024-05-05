import { Button, buttonVariants } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import Link from 'next/link'

export function Sidebar() {
  return (
    <div>
      <h2 className="text-2xl font-semibold p-5">Sentinel Manager</h2>
      <ScrollArea className="border">
        <div className="">
          <Button asChild className="justify-start w-full text-lg" variant="ghost">
            <Link href="/" className="">Dashboard</Link>
          </Button>

          <Button asChild className="justify-start w-full text-lg" variant="ghost">
            <Link href="/">Setting</Link>
          </Button>
        </div>
      </ScrollArea>
    </div>
  )
}