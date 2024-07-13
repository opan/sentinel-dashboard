import { NextRequest, NextResponse } from "next/server"
import { URL } from "url"

export async function DELETE(request: NextRequest, { params }: { params: { id: number }}) {
  const apiUrl = process.env.API_URL
  // const { searchParams } = new URL(request.url)
  const id = params.id
  const res = await fetch(apiUrl + `/sentinel/${id}`,{
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json'
    }
  })

  const jsonRes = await res.json()
  return NextResponse.json(jsonRes)
}