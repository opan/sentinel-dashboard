// import { RedirectStatusCode } from "next/dist/client/components/redirect-status-code"
import { Sentinel, columns } from "./sentinel-columns"
import { DataTable } from "../components/ui/data-table"

async function getSentinel(): Promise<Sentinel[]> {
  const apiUrl = process.env.API_URL
  const res = await fetch(apiUrl + "/sentinel")

  if (!res.ok) {
    throw new Error('Failed to fetch data')
  }

  const resData = await res.json()

  return resData.data
}

export default async function Home() {
  const data = await getSentinel()

  return (
    <div className="container mx-auto py-10">
      <DataTable columns={columns} data={data}/>
    </div>
  )
}
