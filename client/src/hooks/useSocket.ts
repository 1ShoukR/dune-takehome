/* eslint-disable @typescript-eslint/no-unused-vars */
/* eslint-disable @typescript-eslint/no-explicit-any */
import { useEffect, useRef, useCallback } from 'react';

interface WebSocketMessage {
  type: string;
  [key: string]: any;
}

export function useSocket(_path: string) {
  const ws = useRef<WebSocket | null>(null);
  const listeners = useRef<Map<string, (data: any) => void>>(new Map());

  useEffect(() => {
    ws.current = new WebSocket(process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/ws');

    ws.current.onopen = () => {
      console.log('🔌 Connected to WebSocket server');
    };

    ws.current.onmessage = (event) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data);
        console.log('📨 WebSocket message received:', message);
        
        // Call appropriate listener
        const listener = listeners.current.get(message.type);
        if (listener) {
          listener(message);
        }
      } catch (error) {
        console.error('❌ Error parsing WebSocket message:', error);
      }
    };

    ws.current.onclose = () => {
      console.log('🔌❌ Disconnected from WebSocket server');
    };

    ws.current.onerror = (error) => {
      console.error('❌ WebSocket error:', error);
    };

    return () => {
      if (ws.current) {
        ws.current.close();
      }
    };
  }, []);

  const emit = useCallback((type: string, data?: any) => {
    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      const message = { type, ...data };
      ws.current.send(JSON.stringify(message));
    }
  }, []);

  const on = useCallback((event: string, callback: (data: any) => void) => {
    listeners.current.set(event, callback);
  }, []);

  const off = useCallback((event: string) => {
    listeners.current.delete(event);
  }, []);

  return { emit, on, off };
}