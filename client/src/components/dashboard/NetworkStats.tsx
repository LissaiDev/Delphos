import { Network } from '@/types/monitor';

interface NetworkStatsProps {
  networks: Network[];
}

export default function NetworkStats({ networks }: NetworkStatsProps) {
  const formatBytes = (bytes: number) => {
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    if (bytes === 0) return '0 B';
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${sizes[i]}`;
  };

  const formatBytesPerSec = (bytes: number) => {
    return `${formatBytes(bytes)}/s`;
  };

  return (
    <div className="bg-slate-800/50 backdrop-blur-sm rounded-xl p-6 border border-slate-700/50 animate-slide-up">
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-xl font-semibold text-white">Network Statistics</h3>
        <div className="p-2 bg-cyan-600/20 rounded-lg">
          <svg className="w-5 h-5 text-cyan-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8.111 16.404a5.5 5.5 0 017.778 0M12 20h.01m-7.08-7.071c3.904-3.905 10.236-3.905 14.141 0M1.394 9.393c5.857-5.857 15.355-5.857 21.213 0" />
          </svg>
        </div>
      </div>

      <div className="space-y-4">
        {networks.map((network, index) => (
          <div 
            key={network.interfaceName}
            className="bg-slate-700/30 rounded-lg p-4 border border-slate-600/30 hover:border-cyan-500/50 transition-all duration-300 animate-fade-in"
            style={{ animationDelay: `${index * 100}ms` }}
          >
            <div className="flex items-center justify-between mb-3">
              <div className="flex items-center space-x-2">
                <div className="w-3 h-3 bg-cyan-500 rounded-full animate-pulse"></div>
                <span className="text-white font-medium">{network.interfaceName}</span>
              </div>
              <span className="text-slate-400 text-sm">Active</span>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="text-center">
                <div className="flex items-center justify-center space-x-2 mb-1">
                  <svg className="w-4 h-4 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4" />
                  </svg>
                  <span className="text-slate-400 text-sm">Upload</span>
                </div>
                <div className="text-lg font-semibold text-white">
                  {formatBytes(network.totalBytesSent)}
                </div>
                <div className="text-xs text-slate-400">
                  {formatBytesPerSec(network.totalBytesSent / 3600)}
                </div>
              </div>

              <div className="text-center">
                <div className="flex items-center justify-center space-x-2 mb-1">
                  <svg className="w-4 h-4 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 8V4m0 0l-4 4m4-4l4 4M7 20v-4m0 0l4-4m-4 4l-4-4" />
                  </svg>
                  <span className="text-slate-400 text-sm">Download</span>
                </div>
                <div className="text-lg font-semibold text-white">
                  {formatBytes(network.totalBytesRecv)}
                </div>
                <div className="text-xs text-slate-400">
                  {formatBytesPerSec(network.totalBytesRecv / 3600)}
                </div>
              </div>
            </div>

            {/* Network Activity Bar */}
            <div className="mt-3">
              <div className="flex space-x-1">
                <div className="flex-1 bg-slate-600 rounded-full h-1">
                  <div 
                    className="bg-gradient-to-r from-green-500 to-green-600 h-1 rounded-full transition-all duration-500"
                    style={{ 
                      width: `${Math.min((network.totalBytesSent / (network.totalBytesSent + network.totalBytesRecv)) * 100, 100)}%` 
                    }}
                  ></div>
                </div>
                <div className="flex-1 bg-slate-600 rounded-full h-1">
                  <div 
                    className="bg-gradient-to-r from-blue-500 to-blue-600 h-1 rounded-full transition-all duration-500"
                    style={{ 
                      width: `${Math.min((network.totalBytesRecv / (network.totalBytesSent + network.totalBytesRecv)) * 100, 100)}%` 
                    }}
                  ></div>
                </div>
              </div>
              <div className="flex justify-between text-xs text-slate-400 mt-1">
                <span>Upload</span>
                <span>Download</span>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Summary */}
      <div className="mt-6 pt-4 border-t border-slate-700/50">
        <div className="grid grid-cols-2 gap-4 text-center">
          <div>
            <div className="text-2xl font-bold text-white">
              {formatBytes(networks.reduce((acc, net) => acc + net.totalBytesSent, 0))}
            </div>
            <div className="text-slate-400 text-sm">Total Upload</div>
          </div>
          <div>
            <div className="text-2xl font-bold text-white">
              {formatBytes(networks.reduce((acc, net) => acc + net.totalBytesRecv, 0))}
            </div>
            <div className="text-slate-400 text-sm">Total Download</div>
          </div>
        </div>
      </div>
    </div>
  );
} 