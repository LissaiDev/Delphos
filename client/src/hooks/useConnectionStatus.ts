import { useState, useCallback, useMemo } from 'react';

interface ConnectionStatus {
  isConnected: boolean;
  error: string | null;
  isLoading: boolean;
  reconnectAttempts: number;
  maxAttempts: number;
}

interface UseConnectionStatusReturn extends ConnectionStatus {
  setConnected: (connected: boolean) => void;
  setError: (error: string | null) => void;
  setLoading: (loading: boolean) => void;
  incrementAttempts: () => void;
  resetAttempts: () => void;
  getConnectionMessage: () => string;
}

export const useConnectionStatus = (maxAttempts: number = 5): UseConnectionStatusReturn => {
  const [isConnected, setIsConnected] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [reconnectAttempts, setReconnectAttempts] = useState(0);

  const setConnected = useCallback((connected: boolean) => {
    setIsConnected(connected);
  }, []);

  const setErrorState = useCallback((errorMessage: string | null) => {
    setError(errorMessage);
  }, []);

  const setLoading = useCallback((loading: boolean) => {
    setIsLoading(loading);
  }, []);

  const incrementAttempts = useCallback(() => {
    setReconnectAttempts(prev => prev + 1);
  }, []);

  const resetAttempts = useCallback(() => {
    setReconnectAttempts(0);
  }, []);

  const getConnectionMessage = useCallback(() => {
    if (isConnected) return 'Connected';
    if (error) return 'Connection Error';
    if (isLoading) return 'Connecting...';
    return 'Disconnected';
  }, [isConnected, error, isLoading]);

  const status = useMemo(() => ({
    isConnected,
    error,
    isLoading,
    reconnectAttempts,
    maxAttempts
  }), [isConnected, error, isLoading, reconnectAttempts, maxAttempts]);

  return {
    ...status,
    setConnected,
    setError: setErrorState,
    setLoading,
    incrementAttempts,
    resetAttempts,
    getConnectionMessage
  };
}; 