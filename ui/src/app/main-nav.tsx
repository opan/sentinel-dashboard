'use client'

import * as React from "react"
import Link from "next/link"
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  navigationMenuTriggerStyle,
} from "@/components/ui/navigation-menu"

import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogClose,
} from "@/components/ui/dialog"

import { Label } from "@/components/ui/label"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"


export default function MainNav() {

  const [isOpen, setIsOpen] = React.useState(false)
  const [formData, setFormData] = React.useState({ name: "", hosts: "" })

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setFormData({ ...formData, [name]: value })
  }

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()

    const target = e.target as HTMLElement

    console.log('Form submitted: ', formData)

    try {
      const response = await fetch(`/api/sentinel`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData)
      })

      if (response.ok) {
        // reset form
        setFormData({ name: "", hosts: "" })
        setIsOpen(false)
        console.log('Sentinel Cluster created successfully', response)
      } else {
        throw new Error('Failed to create Sentinel Cluster')
      }
    } catch (error) {
      console.error(`Error while submitting form ${target.id}: `, error)
    }
  }

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <NavigationMenu>
        <NavigationMenuList>

          <NavigationMenuItem>
            <Link href="/" legacyBehavior passHref>
              <NavigationMenuLink>
                <span className="font-extrabold text-2xl">Sentinel Manager</span>
              </NavigationMenuLink>
            </Link>
          </NavigationMenuItem>

          <DialogTrigger>
            <NavigationMenuItem>
              <NavigationMenuLink className="">Add Sentinel Cluster</NavigationMenuLink>
            </NavigationMenuItem>
          </DialogTrigger>
        </NavigationMenuList>
      </NavigationMenu>

      <DialogContent>
        <DialogHeader>
          <DialogTitle>Add new Sentinel Cluster</DialogTitle>
        </DialogHeader>

        <DialogDescription>
          Hostname must be filled with comma separated values, e.g: 10.12.13.14:26379,10.12.14.15:26379
        </DialogDescription>

        <form onSubmit={handleSubmit} id="add-sentinel-cluster">
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="name" className="text-right">Name</Label>
              <Input id="name" name="name" className="col-span-3" onChange={handleChange} value={formData.name}></Input>
            </div>

            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="hosts" className="text-right">Hostnames</Label>
              <Input id="hosts" name="hosts" className="col-span-3" onChange={handleChange} value={formData.hosts}></Input>
            </div>
          </div>

          <div className="grid grid-cols-4 items-center gap-4">
            <Button type="submit" className="col-start-2">Create</Button>
            <DialogClose asChild>
              <Button type="button" className="bg-red-500">Cancel</Button>
            </DialogClose>
          </div>
        </form>

      </DialogContent>
    </Dialog>


  )
};