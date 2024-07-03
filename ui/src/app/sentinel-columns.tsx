'use client'

import { ColumnDef } from "@tanstack/react-table"
 
// This type is used to define the shape of our data.
// You can use a Zod schema here if you want.
export type Sentinel = {
  id: number
  name: string
  hosts: string
  created_at: string
}
 
export const columns: ColumnDef<Sentinel>[] = [
  {
    accessorKey: "id",
    header: "ID"
  },
  {
    accessorKey: "name",
    header: "Name",
  },
  {
    accessorKey: "hosts",
    header: () => <div className="text-left">Sentinel Hosts</div>,
    cell: ({row}) => {
      const hosts = (row.getValue("hosts") as string).split(',')
      
      return <div>
        {hosts.map((host, i) => (
          <span className="font-bold mr-1 border p-1" key={i}>{host}</span>
        ))}
      </div>
    }
  },
  {
    accessorKey: "created_at",
    header: "Created At"
  }
]
