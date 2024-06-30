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
    accessorKey: "ID",
    header: "ID"
  },
  {
    accessorKey: "name",
    header: "Name",
  },
  {
    accessorKey: "hosts",
    header: "Hosts",
  },
]
