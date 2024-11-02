'use client'

import { useContext } from "react";
import { Sentinel, columns } from "./sentinel-columns";
import { DataTable } from "../components/ui/data-table";
import { useEffect, useState } from "react"
import { useToast } from "@/hooks/use-toast"
import { SentinelContext } from '../contexts/sentinel-context';

const Home = () => {
  // const [data, setData] = useState<Sentinel[]>([])
  const [error, setError] = useState<string | null>(null)
  const { toast } = useToast()
  const { sentinelContext, setSentinelContext } = useContext(SentinelContext)

  const fetchData = async () => {
    try {
      const response = await fetch(`/api/sentinel`)

      if (!response.ok) {
        throw new Error('Failed to load Sentinel')
      }

      const jsonRes = await response.json()
      // setData(jsonRes.data)
      setSentinelContext(jsonRes.data)

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
    setSentinelContext((prevData) => prevData.filter((item) => item.id != sentinel.id))
  }

  return (
    <DataTable columns={columns(deleteSentinel)} data={sentinelContext} />
  )
}

export default Home;