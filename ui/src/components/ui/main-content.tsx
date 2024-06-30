'use client'
import { TooltipProvider } from "@/components/ui/tooltip"
import { cn } from "@/lib/utils"

import * as React from "react"

import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable"

interface MainProps {
  defaultCollapsed?: boolean
}

export function MainContent({
  defaultCollapsed = false
}: MainProps) {
  const [isCollapsed, setIsCollapsed] = React.useState(defaultCollapsed)
  return (
    <TooltipProvider delayDuration={0}>
      <ResizablePanelGroup
        direction="horizontal"
        onLayout={(sizes: number[]) => {
          document.cookie = `react-resizable-panels:layout=${JSON.stringify(
            sizes
          )}`
        }}
        className="h-full max-h-[800px] items-stretch"
      >
        <ResizablePanel
          defaultSize={150}
          collapsedSize={4}
          collapsible={true}
          minSize={15}
          maxSize={20}
          onCollapse={(collapsed) => {
            setIsCollapsed(collapsed)
            document.cookie = `react-resizable-panels:collapsed=${JSON.stringify(
              collapsed
            )}`
          }}
          className={cn(
            isCollapsed &&
              "min-w-[50px] transition-all duration-300 ease-in-out"
          )}
        >
          Sidebar panel
        </ResizablePanel>
        <ResizableHandle withHandle/>
        <ResizablePanel>
          Main Panel
        </ResizablePanel>
      </ResizablePanelGroup>
    </TooltipProvider>
  )
}