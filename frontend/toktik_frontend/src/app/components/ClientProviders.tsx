// // src/app/components/ClientProviders.tsx
// 'use client';

// import React from 'react';
// import { useAuth } from '../contexts/AuthContext';
// import { SocketProvider } from '../contexts/SocketContext';

// export default function ClientProviders({ children }: { children: React.ReactNode }) {
//   const { jwtToken } = useAuth();

//   return (
//     <SocketProvider jwtToken={jwtToken}>
//       {children}
//     </SocketProvider>
//   );
// }



'use client';

import React from 'react';
import { useAuth } from '../contexts/AuthContext';
import { SocketProvider } from '../contexts/SocketContext';

export default function ClientProviders({ children }: { children: React.ReactNode }) {
  const { jwtToken, loaded } = useAuth();

  // Wait for token to load from localStorage
  if (!loaded) return null;

  return (
    <SocketProvider key={jwtToken ?? 'no-token'} jwtToken={jwtToken}>
      {children}
    </SocketProvider>
  );
}

