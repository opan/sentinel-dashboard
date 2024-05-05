import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { Sidebar } from "@/components/ui/sidebar"

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Sentinel Manager App",
  description: "Tools to manage multiple Redis Sentinels servers",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <div className="grid grid-cols-6 h-screen">
          <div className="col-span-1">
            <Sidebar />
          </div>
          <div className="col-start-2 col-end-7">
            {children}
          </div>
        </div>
      </body>
    </html>
  );
}
