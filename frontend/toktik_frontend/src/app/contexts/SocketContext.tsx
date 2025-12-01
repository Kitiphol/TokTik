// // src/app/context/SocketContext.tsx
// 'use client';

// import React, { createContext, useContext, useEffect, useState } from 'react';
// import { io, Socket } from 'socket.io-client';

// type SocketContextType = {
//   socket: Socket | null;
//   isConnected: boolean;
// };

// const SocketContext = createContext<SocketContextType>({
//   socket: null,
//   isConnected: false,
// });

// export const SocketProvider = ({
//   children,
//   jwtToken,
// }: {
//   children: React.ReactNode;
//   jwtToken: string | null;
// }) => {
//   const [socket, setSocket] = useState<Socket | null>(null);
//   const [isConnected, setIsConnected] = useState(false);

//   useEffect(() => {
//     if (!jwtToken) return;

//     // Connect socket with auth token
//     const newSocket = io(undefined, {
//       auth: { token: jwtToken },
//       transports: ['websocket'],
//     });

//     // Save to state
//     setSocket(newSocket);

//     // Setup event listeners
//     const handleConnect = () => {
//       console.log('[Socket] Connected');
//       setIsConnected(true);
//     };

//     const handleDisconnect = () => {
//       console.log('[Socket] Disconnected');
//       setIsConnected(false);
//     };

//     const handleConnectError = (err: unknown) => {
//       console.error('[Socket] Connection Error:', err);
//       setIsConnected(false);
//     };

//     newSocket.on('connect', () => {
//     console.log('[Socket] Connected id=', newSocket.id);
//     setIsConnected(true);
//     });
//     newSocket.on('disconnect', handleDisconnect);
//     newSocket.on('connect_error', handleConnectError);

//     // Cleanup on unmount
//     return () => {
//       newSocket.off('connect', handleConnect);
//       newSocket.off('disconnect', handleDisconnect);
//       newSocket.off('connect_error', handleConnectError);
//       newSocket.disconnect();
//     };
//   }, [jwtToken]);

//   return (
//     <SocketContext.Provider value={{ socket, isConnected }}>
//       {children}
//     </SocketContext.Provider>
//   );
// };

// export const useSocket = () => useContext(SocketContext);



// src/app/contexts/SocketContext.tsx
'use client';

import React, { createContext, useContext, useEffect, useState, useRef } from 'react';
import { io, Socket } from 'socket.io-client';

type SocketContextType = {
  socket: Socket | null;
  isConnected: boolean;
};

const SocketContext = createContext<SocketContextType>({
  socket: null,
  isConnected: false,
});

export const SocketProvider = ({
  children,
  jwtToken,
}: {
  children: React.ReactNode;
  jwtToken: string | null;
}) => {
  const [socket, setSocket] = useState<Socket | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const socketRef = useRef<Socket | null>(null);

  useEffect(() => {
    if (!jwtToken) {
      if (socketRef.current) {
        socketRef.current.disconnect();
        socketRef.current = null;
        setSocket(null);
        setIsConnected(false);
      }
      return;
    }

    const newSocket = io(undefined, {
      auth: { token: jwtToken },
      transports: ['websocket'],
    });

    socketRef.current = newSocket;
    setSocket(newSocket);

    newSocket.on('connect', () => {
      console.log('[Socket] Connected:', newSocket.id);
      setIsConnected(true);
    });

    newSocket.on('disconnect', () => {
      console.log('[Socket] Disconnected');
      setIsConnected(false);
    });

    newSocket.on('connect_error', (err) => {
      console.error('[Socket] Connection Error:', err);
      setIsConnected(false);
    });

    return () => {
      newSocket.disconnect();
      socketRef.current = null;
      setIsConnected(false);
      setSocket(null);
    };
  }, [jwtToken]);

  return (
    <SocketContext.Provider value={{ socket, isConnected }}>
      {children}
    </SocketContext.Provider>
  );
};

export const useSocket = () => useContext(SocketContext);
