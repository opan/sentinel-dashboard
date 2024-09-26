'use client'
// import { SentinelTable }  from "./sentinel-table"

import { Sentinel, columns } from "./sentinel-columns";
import { DataTable } from "../components/ui/data-table";
import { useEffect, useState } from "react"

// async function getSentinel(): Promise<Sentinel[]> {
//   const res = await fetch(`/api/sentinel`)

//   if (!res.ok) {
//     throw new Error("Failed to remove sentinel");
//   }

//   const data = await res.json()
//   return data
// }

const Home = () => {
  const [data, setData] = useState<Sentinel[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      const response = await fetch(`/api/sentinel`)
      
      if (!response.ok) {
        throw new Error('Failed to load Sentinel')
      }

      const jsonRes = await response.json()
      setData(jsonRes.data)
    }

    fetchData()
  }, [])

  return (
    <DataTable columns={columns} data={data}/>
  )
}

export default Home;