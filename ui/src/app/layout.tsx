import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import "./main-nav";
import MainNav from "./main-nav";

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
        <div className="grid grid-flow-col auto-cols-max m-10 gap-y-5">
          <div>
            <span className="text-2xl">Sentinel Manager</span>
          </div>
          <MainNav />
        </div>
        <div className="container mx-auto">
          {children}
        </div>
      </body>
    </html>
  );
}
