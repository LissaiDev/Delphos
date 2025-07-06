import { Memory, CPU } from '@/types/monitor';

interface ResourceCardsProps {
  memory: Memory;
  cpu: CPU[];
  className?: string;
}

export default function ResourceCards({ memory, cpu, className = '' }: ResourceCardsProps) {
  const formatBytes = (bytes: number) => {
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    if (bytes === 0) return '0 B';
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${sizes[i]}`;
  };

  const memoryUsagePercent = ((memory.used / memory.total) * 100).toFixed(1);
  const swapUsagePercent = memory.swapTotal > 0 ? ((memory.swapUsed / memory.swapTotal) * 100).toFixed(1) : '0.0';
  const avgCpuUsage = cpu.reduce((acc, core) => acc + core.usage, 0) / cpu.length;

  return (
    <div className={`grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 ${className}`}>
      {/* CPU Usage */}
      <div className="bg-slate-800/50 backdrop-blur-sm rounded-xl p-6 border border-slate-700/50 hover:border-purple-500/50 transition-all duration-300 animate-slide-up">
        <div className="flex items-center justify-between mb-4">
          <div className="p-2 bg-red-600/20 rounded-lg">
            <svg className="w-5 h-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
            </svg>
          </div>
          <span className="text-slate-400 text-sm">{cpu.length} cores</span>
        </div>
        <h3 className="text-white font-semibold mb-2">CPU Usage</h3>
        <div className="text-3xl font-bold text-white mb-2">{avgCpuUsage.toFixed(1)}%</div>
        <div className="w-full bg-slate-700 rounded-full h-2 mb-2">
          <div 
            className="bg-gradient-to-r from-red-500 to-red-600 h-2 rounded-full transition-all duration-500"
            style={{ width: `${avgCpuUsage}%` }}
          ></div>
        </div>
        <p className="text-slate-400 text-xs truncate">{cpu[0]?.model}</p>
      </div>

      {/* Memory Usage */}
      <div className="bg-slate-800/50 backdrop-blur-sm rounded-xl p-6 border border-slate-700/50 hover:border-blue-500/50 transition-all duration-300 animate-slide-up">
        <div className="flex items-center justify-between mb-4">
          <div className="p-2 bg-blue-600/20 rounded-lg">
            <svg className="w-5 h-5 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <span className="text-slate-400 text-sm">RAM</span>
        </div>
        <h3 className="text-white font-semibold mb-2">Memory Usage</h3>
        <div className="text-3xl font-bold text-white mb-2">{memoryUsagePercent}%</div>
        <div className="w-full bg-slate-700 rounded-full h-2 mb-2">
          <div 
            className="bg-gradient-to-r from-blue-500 to-blue-600 h-2 rounded-full transition-all duration-500"
            style={{ width: `${memoryUsagePercent}%` }}
          ></div>
        </div>
        <p className="text-slate-400 text-xs">
          {formatBytes(memory.used)} / {formatBytes(memory.total)}
        </p>
      </div>

      {/* Swap Usage */}
      <div className="bg-slate-800/50 backdrop-blur-sm rounded-xl p-6 border border-slate-700/50 hover:border-yellow-500/50 transition-all duration-300 animate-slide-up">
        <div className="flex items-center justify-between mb-4">
          <div className="p-2 bg-yellow-600/20 rounded-lg">
            <svg className="w-5 h-5 text-yellow-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4" />
            </svg>
          </div>
          <span className="text-slate-400 text-sm">Swap</span>
        </div>
        <h3 className="text-white font-semibold mb-2">Swap Usage</h3>
        <div className="text-3xl font-bold text-white mb-2">{swapUsagePercent}%</div>
        <div className="w-full bg-slate-700 rounded-full h-2 mb-2">
          <div 
            className="bg-gradient-to-r from-yellow-500 to-yellow-600 h-2 rounded-full transition-all duration-500"
            style={{ width: `${swapUsagePercent}%` }}
          ></div>
        </div>
        <p className="text-slate-400 text-xs">
          {formatBytes(memory.swapUsed)} / {formatBytes(memory.swapTotal)}
        </p>
      </div>

      {/* System Load */}
      <div className="bg-slate-800/50 backdrop-blur-sm rounded-xl p-6 border border-slate-700/50 hover:border-green-500/50 transition-all duration-300 animate-slide-up">
        <div className="flex items-center justify-between mb-4">
          <div className="p-2 bg-green-600/20 rounded-lg">
            <svg className="w-5 h-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <span className="text-slate-400 text-sm">Load</span>
        </div>
        <h3 className="text-white font-semibold mb-2">System Load</h3>
        <div className="text-3xl font-bold text-white mb-2">
          {(avgCpuUsage / 100).toFixed(2)}
        </div>
        <div className="w-full bg-slate-700 rounded-full h-2 mb-2">
          <div 
            className="bg-gradient-to-r from-green-500 to-green-600 h-2 rounded-full transition-all duration-500"
            style={{ width: `${Math.min(avgCpuUsage, 100)}%` }}
          ></div>
        </div>
        <p className="text-slate-400 text-xs">
          {avgCpuUsage > 80 ? 'High' : avgCpuUsage > 50 ? 'Medium' : 'Low'} load
        </p>
      </div>
    </div>
  );
} 