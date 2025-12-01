
"use client";

import { ReactNode, useEffect, useState } from "react";
import { SocketProvider } from "../../contexts/SocketContext";
import { getJWT } from "../Auth/auth"; // Your client-side JWT getter (from localStorage or cookies)

export function ClientSocketProvider({ children }: { children: ReactNode }) {
  const [jwtToken, setJwtToken] = useState<string | null>(null);

  useEffect(() => {
    const token = getJWT();
    setJwtToken(token);
  }, []);

  if (!jwtToken) {
    // You can show a loader or just return children without socket connection until token loads
    return <>{children}</>;
  }

  return <SocketProvider jwtToken={jwtToken}>{children}</SocketProvider>;
}
