import { useState, useEffect, useCallback, useRef } from "react";
import { Monitor } from "@/types/monitor";

interface UseMonitorDataReturn {
  data: Monitor | null;
  isLoading: boolean;
  error: string | null;
  isConnected: boolean;
  refresh: () => void;
  reconnect: () => void;
}

export const useMonitorData = (endpoint?: string): UseMonitorDataReturn => {
  const [data, setData] = useState<Monitor | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isConnected, setIsConnected] = useState(false);

  const eventSourceRef = useRef<EventSource | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const reconnectAttemptsRef = useRef(0);
  const maxReconnectAttempts = 5;
  const baseReconnectDelay = 1000; // 1 second

  const cleanup = useCallback(() => {
    if (eventSourceRef.current) {
      eventSourceRef.current.close();
      eventSourceRef.current = null;
    }
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }
  }, []);

  const connect = useCallback(() => {
    if (!endpoint) {
      setIsLoading(false);
      setIsConnected(false);
      return;
    }

    cleanup();

    try {
      setIsLoading(true);
      setError(null);

      const eventSource = new EventSource(endpoint);
      eventSourceRef.current = eventSource;

      eventSource.onopen = () => {
        setIsConnected(true);
        setIsLoading(false);
        setError(null);
        reconnectAttemptsRef.current = 0;
      };

      eventSource.onmessage = (event) => {
        console.log(`Received event: ${event}`);
        try {
          const monitorData: Monitor = JSON.parse(event.data);
          setData(monitorData);
          setError(null);
        } catch (err) {
          console.error("Failed to parse SSE data:", err);
          setError("Invalid data format received");
        }
      };

      eventSource.onerror = (event) => {
        console.error("SSE connection error:", event);
        setIsConnected(false);
        setIsLoading(false);

        if (reconnectAttemptsRef.current < maxReconnectAttempts) {
          const delay =
            baseReconnectDelay * Math.pow(2, reconnectAttemptsRef.current);
          reconnectAttemptsRef.current++;

          setError(
            `Connection lost. Reconnecting in ${delay / 1000}s... (Attempt ${reconnectAttemptsRef.current}/${maxReconnectAttempts})`,
          );

          reconnectTimeoutRef.current = setTimeout(() => {
            connect();
          }, delay);
        } else {
          setError(
            "Connection failed after multiple attempts. Please check your connection and try again.",
          );
        }
      };
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to establish connection",
      );
      setIsLoading(false);
      setIsConnected(false);
    }
  }, [endpoint, cleanup]);

  const reconnect = useCallback(() => {
    reconnectAttemptsRef.current = 0;
    connect();
  }, [connect]);

  const refresh = useCallback(() => {
    if (isConnected && eventSourceRef.current) {
      // For SSE, we can't manually refresh, but we can reconnect
      reconnect();
    } else {
      connect();
    }
  }, [isConnected, connect, reconnect]);

  useEffect(() => {
    connect();

    return () => {
      cleanup();
    };
  }, [connect, cleanup]);

  return {
    data,
    isLoading,
    error,
    isConnected,
    refresh,
    reconnect,
  };
};
