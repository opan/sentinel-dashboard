'use client'

import React, { createContext, useState } from "react";
import { Sentinel } from "../app/sentinel-columns";

interface SharedContextProps {
  sentinelContext: Sentinel[];
  setSentinelContext: React.Dispatch<React.SetStateAction<Sentinel[]>>;
}

export const SentinelContext = createContext<SharedContextProps>({
  sentinelContext: [],
  setSentinelContext: () => {},
})

export const SentinelContextProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [sentinelContext, setSentinelContext] = useState<Sentinel[]>([]);

  return (
    <SentinelContext.Provider value={{ sentinelContext, setSentinelContext }}>
      {children}
    </SentinelContext.Provider>
  );
};
