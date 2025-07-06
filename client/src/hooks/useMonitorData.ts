import { useState, useEffect, useCallback } from 'react';
import { Monitor } from '@/types/monitor';

interface UseMonitorDataReturn {
  data: Monitor | null;
  isLoading: boolean;
  error: string | null;
  refresh: () => void;
}

export const useMonitorData = (endpoint?: string): UseMonitorDataReturn => {
  const [data, setData] = useState<Monitor | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchData = useCallback(async () => {
    if (!endpoint) {
      // Use mock data if no endpoint is provided
      setIsLoading(false);
      return;
    }

    try {
      setIsLoading(true);
      setError(null);
      
      const response = await fetch(endpoint);
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      const monitorData: Monitor = await response.json();
      setData(monitorData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch data');
    } finally {
      setIsLoading(false);
    }
  }, [endpoint]);

  const refresh = useCallback(() => {
    fetchData();
  }, [fetchData]);

  useEffect(() => {
    fetchData();
    
    // Set up polling if endpoint is provided
    if (endpoint) {
      const interval = setInterval(fetchData, 5000); // Refresh every 5 seconds
      return () => clearInterval(interval);
    }
  }, [fetchData, endpoint]);

  return {
    data,
    isLoading,
    error,
    refresh
  };
}; 