import { NextRequest, NextResponse } from "next/server"

export async function GET(
  request: NextRequest
) {
  const apiUrl = process.env.API_URL
  const res = await fetch(apiUrl + `/sentinel`,{
    method: 'GET',
    next: { revalidate: 360 },
    headers: {
      'Content-Type': 'application/json'
    }
  })

  const jsonRes = await res.json()
  return NextResponse.json(jsonRes)
}