import { Host } from '@/types/monitor';

interface SystemOverviewProps {
  host: Host;
}

export default function SystemOverview({ host }: SystemOverviewProps) {
  const formatUptime = (seconds: number) => {
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    
    if (days > 0) {
      return `${days}d ${hours}h ${minutes}m`;
    } else if (hours > 0) {
      return `${hours}h ${minutes}m`;
    } else {
      return `${minutes}m`;
    }
  };

  return (
    <div className="bg-slate-800/50 backdrop-blur-sm rounded-xl p-6 mb-8 border border-slate-700/50 animate-slide-up">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="flex items-center space-x-4">
          <div className="p-3 bg-purple-600/20 rounded-lg">
            <svg className="w-6 h-6 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
          </div>
          <div>
            <p className="text-slate-400 text-sm">Hostname</p>
            <p className="text-white font-semibold">{host.hostname}</p>
          </div>
        </div>

        <div className="flex items-center space-x-4">
          <div className="p-3 bg-blue-600/20 rounded-lg">
            <svg className="w-6 h-6 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
            </svg>
          </div>
          <div>
            <p className="text-slate-400 text-sm">Operating System</p>
            <p className="text-white font-semibold text-sm">{host.os}</p>
          </div>
        </div>

        <div className="flex items-center space-x-4">
          <div className="p-3 bg-green-600/20 rounded-lg">
            <svg className="w-6 h-6 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div>
            <p className="text-slate-400 text-sm">Uptime</p>
            <p className="text-white font-semibold">{formatUptime(host.uptime)}</p>
          </div>
        </div>
      </div>
    </div>
  );
} 