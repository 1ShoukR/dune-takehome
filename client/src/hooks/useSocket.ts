import { useEffect, useRef } from 'react';
import { io, Socket } from 'socket.io-client';

export function useSocket(serverPath: string) {
  const socket = useRef<Socket | null>(null);

  useEffect(() => {
    socket.current = io('http://localhost:8080', {
      path: '/socket.io/',
      transports: ['websocket'],
    });

    socket.current.on('connect', () => {
      console.log('🔌 Connected to WebSocket server');
    });

    socket.current.on('disconnect', () => {
      console.log('🔌❌ Disconnected from WebSocket server');
    });

    socket.current.on('connect_error', (error) => {
      console.error('❌ WebSocket connection error:', error);
    });

    return () => {
      if (socket.current) {
        socket.current.disconnect();
      }
    };
  }, [serverPath]);

  return socket.current;
}
