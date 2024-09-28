'use client'
// import { SentinelTable }  from "./sentinel-table"

import { Sentinel, columns } from "./sentinel-columns";
import { DataTable } from "../components/ui/data-table";
import { useEffect, useState } from "react"

const Home = () => {
  const [data, setData] = useState<Sentinel[]>([]);
  const fetchData = async () => {
    const response = await fetch(`/api/sentinel`)
    
    if (!response.ok) {
      throw new Error('Failed to load Sentinel')
    }

    const jsonRes = await response.json()
    setData(jsonRes.data)
  }

  useEffect(() => {
    fetchData()
  }, [])

  return (
    <DataTable columns={columns} data={data}/>
  )
}

export default Home;