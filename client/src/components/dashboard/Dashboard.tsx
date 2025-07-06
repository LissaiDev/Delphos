'use client';

import React from 'react';
import { Monitor } from '@/types/monitor';
import { useMonitorData } from '@/hooks/useMonitorData';
import { useNotifications } from '@/hooks/useNotifications';
import SystemOverview from '@/components/dashboard/SystemOverview';
import ResourceCards from '@/components/dashboard/ResourceCards';
import NetworkStats from '@/components/dashboard/NetworkStats';
import DiskUsage from '@/components/dashboard/DiskUsage';
import Notification from '@/components/ui/Notification';
import LoadingSpinner from '@/components/ui/LoadingSpinner';

interface DashboardProps {
  endpoint?: string;
  data?: Monitor;
  isLoading?: boolean;
}

export default function Dashboard({ endpoint, data, isLoading: externalLoading }: DashboardProps) {
  const { data: fetchedData, isLoading, error, isConnected, refresh, reconnect } = useMonitorData(endpoint);
  const { notifications, addNotification, removeNotification } = useNotifications();

  // Show notifications for connection status changes
  React.useEffect(() => {
    if (isConnected) {
      addNotification('Connected to server successfully', 'success', 3000);
    }
  }, [isConnected, addNotification]);

  React.useEffect(() => {
    if (error && !isConnected) {
      addNotification(error, 'error', 0); // Don't auto-dismiss error notifications
    }
  }, [error, isConnected, addNotification]);

  // Simular dados para demonstração (remover em produção)
  const mockData: Monitor = {
    host: {
      hostname: 'delphos-server',
      os: 'Linux 6.12.33+deb13-amd64',
      uptime: 86400
    },
    memory: {
      total: 16777216,
      used: 8388608,
      free: 8388608,
      swapTotal: 4194304,
      swapUsed: 1048576,
      swapFree: 3145728
    },
    cpu: [
      {
        usage: 45.2,
        model: 'Intel(R) Core(TM) i7-10700K CPU @ 3.80GHz',
        cores: 8
      }
    ],
    disk: [
      {
        mountpoint: '/',
        type: 'ext4',
        total: 107374182400,
        used: 64424509440,
        free: 42949672960,
        usedPercent: 60.0
      },
      {
        mountpoint: '/home',
        type: 'ext4',
        total: 214748364800,
        used: 107374182400,
        free: 107374182400,
        usedPercent: 50.0
      }
    ],
    network: [
      {
        interfaceName: 'eth0',
        totalBytesSent: 1073741824,
        totalBytesRecv: 2147483648
      },
      {
        interfaceName: 'wlan0',
        totalBytesSent: 536870912,
        totalBytesRecv: 1073741824
      }
    ]
  };

  const currentData = data || fetchedData || mockData;
  const isCurrentlyLoading = externalLoading !== undefined ? externalLoading : isLoading;

  const handleRefresh = () => {
    refresh();
  };

  if (isCurrentlyLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 p-4">
        <div className="max-w-7xl mx-auto">
          <div className="flex flex-col items-center justify-center min-h-[60vh]">
            <LoadingSpinner size="lg" className="mb-6" />
            <h2 className="text-2xl font-semibold text-white mb-2 animate-fade-in">
              Connecting to Server...
            </h2>
            <p className="text-slate-400 text-center max-w-md animate-fade-in-delay">
              Establishing real-time connection to monitor system resources
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 p-4">
      {/* Notifications */}
      {notifications.map(notification => (
        <Notification
          key={notification.id}
          message={notification.message}
          type={notification.type}
          duration={notification.duration}
          onClose={() => removeNotification(notification.id)}
        />
      ))}
      
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-8">
          <div>
            <h1 className="text-3xl font-bold text-white mb-2 animate-fade-in">
              System Monitor
            </h1>
            <p className="text-slate-300 animate-fade-in-delay">
              {currentData.host.hostname} • {currentData.host.os}
            </p>
            {/* Connection Status */}
            <div className="flex items-center gap-2 mt-2 animate-fade-in-delay">
              <div className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-500 animate-pulse' : 'bg-red-500'}`}></div>
              <span className={`text-sm ${isConnected ? 'text-green-400' : 'text-red-400'}`}>
                {isConnected ? 'Connected' : 'Disconnected'}
              </span>
              {error && (
                <span className="text-sm text-yellow-400 ml-2">
                  {error}
                </span>
              )}
            </div>
          </div>
          <div className="flex gap-2 mt-4 sm:mt-0">
            <button
              onClick={reconnect}
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-all duration-300 hover:scale-105 animate-fade-in-delay-2"
            >
              <span className="flex items-center gap-2">
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                Reconnect
              </span>
            </button>
            <button
              onClick={handleRefresh}
              className="px-6 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-lg transition-all duration-300 hover:scale-105 animate-fade-in-delay-2"
            >
              <span className="flex items-center gap-2">
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
                Refresh
              </span>
            </button>
          </div>
        </div>

        {/* System Overview */}
        <SystemOverview host={currentData.host} />

        {/* Resource Cards */}
        <ResourceCards 
          memory={currentData.memory} 
          cpu={currentData.cpu} 
          className="mb-8" 
        />

        {/* Charts Section */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <NetworkStats networks={currentData.network} />
          <DiskUsage disks={currentData.disk} />
        </div>
      </div>
    </div>
  );
} 