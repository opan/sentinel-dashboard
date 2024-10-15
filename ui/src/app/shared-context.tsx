'use client'

import React, { createContext, useState } from "react";
import { Sentinel } from "./sentinel-columns";
import { create } from "domain";

interface SharedContextProps {
  sharedData: Sentinel[];
  setSharedData: React.Dispatch<React.SetStateAction<Sentinel[]>>;
}

export const SharedContext = createContext<SharedContextProps>({
  sharedData: [],
  setSharedData: () => {},
})

export const SharedContextProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [sharedData, setSharedData] = useState<Sentinel[]>([]);

  return (
    <SharedContext.Provider value={{ sharedData, setSharedData }}>
      {children}
    </SharedContext.Provider>
  );
};
