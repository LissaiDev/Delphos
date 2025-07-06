'use client';

import { Monitor } from '@/types/monitor';
import { useMonitorData } from '@/hooks/useMonitorData';
import SystemOverview from '@/components/dashboard/SystemOverview';
import ResourceCards from '@/components/dashboard/ResourceCards';
import NetworkStats from '@/components/dashboard/NetworkStats';
import DiskUsage from '@/components/dashboard/DiskUsage';

interface DashboardProps {
  endpoint?: string;
  data?: Monitor;
  isLoading?: boolean;
}

export default function Dashboard({ endpoint, data, isLoading: externalLoading }: DashboardProps) {
  const { data: fetchedData, isLoading, error, refresh } = useMonitorData(endpoint);

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
          <div className="animate-pulse">
            <div className="h-8 bg-slate-700 rounded-lg mb-8 w-1/3"></div>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
              {[...Array(4)].map((_, i) => (
                <div key={i} className="h-32 bg-slate-800 rounded-xl"></div>
              ))}
            </div>
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <div className="h-96 bg-slate-800 rounded-xl"></div>
              <div className="h-96 bg-slate-800 rounded-xl"></div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 p-4">
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
          </div>
          <button
            onClick={handleRefresh}
            className="mt-4 sm:mt-0 px-6 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-lg transition-all duration-300 hover:scale-105 animate-fade-in-delay-2"
          >
            <span className="flex items-center gap-2">
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              Refresh
            </span>
          </button>
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