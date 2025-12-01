import React, { useRef, useEffect } from 'react';
import { io } from 'socket.io-client';
import type { Socket } from 'socket.io-client';

export default function Test() {
  const socketRef = useRef<Socket | null>(null);

  useEffect(() => {
    const socket = io('http://localhost:3001');
    socketRef.current = socket;

    socket.on('connect', () => {
      console.log('Socket connected:', socket.id);
    });

    return () => {
      socket.disconnect();
    };
  }, []);

  return <div>Socket.io Client Test</div>;
}
