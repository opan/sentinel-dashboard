'use client'

import { useContext } from "react";
import { Sentinel, columns } from "./sentinel-columns";
import { DataTable } from "../components/ui/data-table";
import { useEffect, useState } from "react"
import { useToast } from "@/hooks/use-toast"
import { SharedContext } from './shared-context';

const Home = () => {
  const [data, setData] = useState<Sentinel[]>([])
  const [error, setError] = useState<string | null>(null)
  const { toast } = useToast()
  const { sharedData, setSharedData } = useContext(SharedContext)

  const fetchData = async () => {
    try {
      const response = await fetch(`/api/sentinel`)

      if (!response.ok) {
        throw new Error('Failed to load Sentinel')
      }

      const jsonRes = await response.json()
      setData(jsonRes.data)
      setSharedData(jsonRes.data)
    } catch (error: any) {
      const errMsg = "Error when fetching sentinels"
      console.error(`${errMsg}: `, error)
      setError(error.message || 'Something went wrong')
      toast({
        title: 'Error',
        description: errMsg,
        variant: 'destructive'
      })
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  const deleteSentinel = (sentinel: Sentinel) => {
    setData((prevData) => prevData.filter((item) => item.id != sentinel.id))
    setSharedData((prevData) => prevData.filter((item) => item.id != sentinel.id))
  }

  return (
    <DataTable columns={columns(deleteSentinel)} data={data} />
  )
}

export default Home;