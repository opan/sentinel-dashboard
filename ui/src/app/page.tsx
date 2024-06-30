import { Sentinel, columns } from "./sentinel-columns"
import { DataTable } from "./sentinel-data-table"

export async function getSentinel() {
  const apiUrl = process.env.API_URL
  const res = await fetch(apiUrl + "/sentinel")

  if (!res.ok) {
    throw new Error('Failed to fetch data')
  }

  return res.json()
}

export default async function Home() {
  const sentinel = await getSentinel()

  return (
    <div className="container mx-auto py-10">
      <DataTable columns={columns} data={sentinel.data}/>
    </div>
  )
}
