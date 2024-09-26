'use client'

import { Sentinel } from "./sentinel-columns"

export async function GetSentinel(): Promise<Sentinel[]> {
  const res = await fetch(`/api/sentinel`)

  if (!res.ok) {
    throw new Error("Failed to remove sentinel");
  }

  const data = await res.json()
  return data
}