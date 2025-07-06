import React, { useMemo, useCallback } from 'react';
import { Disk } from '@/types/monitor';
import { motion } from 'framer-motion';

interface DiskUsageProps {
  disks: Disk[];
}

// Memoize individual disk item to prevent unnecessary re-renders
const DiskItem = React.memo(({ disk, index }: { disk: Disk; index: number }) => {
  const formatBytes = useCallback((bytes: number) => {
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    if (bytes === 0) return '0 B';
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${sizes[i]}`;
  }, []);

  const getUsageColor = useCallback((percent: number) => {
    if (percent >= 90) return 'from-red-500 to-red-600';
    if (percent >= 75) return 'from-yellow-500 to-yellow-600';
    if (percent >= 50) return 'from-orange-500 to-orange-600';
    return 'from-green-500 to-green-600';
  }, []);

  const getUsageStatus = useCallback((percent: number) => {
    if (percent >= 90) return 'Critical';
    if (percent >= 75) return 'Warning';
    if (percent >= 50) return 'Moderate';
    return 'Good';
  }, []);

  const statusClass = useMemo(() => {
    if (disk.usedPercent >= 90) return 'bg-red-500/20 text-red-400';
    if (disk.usedPercent >= 75) return 'bg-yellow-500/20 text-yellow-400';
    if (disk.usedPercent >= 50) return 'bg-orange-500/20 text-orange-400';
    return 'bg-green-500/20 text-green-400';
  }, [disk.usedPercent]);

  return (
    <motion.div
      className="bg-slate-700/30 rounded-lg p-4 border border-slate-600/30 hover:border-orange-500/50"
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{
        delay: index * 0.1,
        duration: 0.3,
        ease: "easeOut",
      }}
      whileHover={{
        y: -1,
        transition: { duration: 0.2 }
      }}
    >
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center space-x-3">
          <div className="p-2 bg-orange-600/20 rounded-lg">
            <svg className="w-4 h-4 text-orange-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4" />
            </svg>
          </div>
          <div>
            <div className="text-white font-medium">{disk.mountpoint}</div>
            <div className="text-slate-400 text-sm">{disk.type.toUpperCase()}</div>
          </div>
        </div>
        <div className="text-right">
          <div className={`text-sm font-medium px-2 py-1 rounded-full ${statusClass}`}>
            {getUsageStatus(disk.usedPercent)}
          </div>
        </div>
      </div>

      <div className="mb-3">
        <div className="flex justify-between text-sm text-slate-400 mb-1">
          <span>Usage</span>
          <span>{disk.usedPercent.toFixed(1)}%</span>
        </div>
        <div className="w-full bg-slate-600 rounded-full h-2">
          <div 
            className={`bg-gradient-to-r ${getUsageColor(disk.usedPercent)} h-2 rounded-full transition-all duration-500`}
            style={{ width: `${disk.usedPercent}%` }}
          ></div>
        </div>
      </div>

      <div className="grid grid-cols-3 gap-4 text-center">
        <div>
          <div className="text-lg font-semibold text-white">
            {formatBytes(disk.used)}
          </div>
          <div className="text-slate-400 text-xs">Used</div>
        </div>
        <div>
          <div className="text-lg font-semibold text-white">
            {formatBytes(disk.free)}
          </div>
          <div className="text-slate-400 text-xs">Free</div>
        </div>
        <div>
          <div className="text-lg font-semibold text-white">
            {formatBytes(disk.total)}
          </div>
          <div className="text-slate-400 text-xs">Total</div>
        </div>
      </div>
    </motion.div>
  );
});

DiskItem.displayName = 'DiskItem';

export default React.memo(function DiskUsage({ disks }: DiskUsageProps) {
  const formatBytes = useCallback((bytes: number) => {
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    if (bytes === 0) return '0 B';
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${sizes[i]}`;
  }, []);

  const totalStorage = useMemo(() => 
    disks.reduce((acc, disk) => acc + disk.total, 0), 
    [disks]
  );

  const totalFree = useMemo(() => 
    disks.reduce((acc, disk) => acc + disk.free, 0), 
    [disks]
  );

  const renderDiskItem = useCallback((disk: Disk, index: number) => (
    <DiskItem key={disk.mountpoint} disk={disk} index={index} />
  ), []);

  return (
    <motion.div
      className="bg-slate-800/50 rounded-xl p-6 border border-slate-700/50"
      initial={{ opacity: 0, y: 20, scale: 0.95 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      transition={{
        duration: 0.4,
        ease: "easeOut",
      }}
    >
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-xl font-semibold text-white">Disk Usage</h3>
        <div className="p-2 bg-orange-600/20 rounded-lg">
          <svg className="w-5 h-5 text-orange-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4" />
          </svg>
        </div>
      </div>

      <div className="space-y-4">
        {disks.map(renderDiskItem)}
      </div>

      {/* Summary */}
      <div className="mt-6 pt-4 border-t border-slate-700/50">
        <div className="grid grid-cols-2 gap-4 text-center">
          <div>
            <div className="text-2xl font-bold text-white">
              {formatBytes(totalStorage)}
            </div>
            <div className="text-slate-400 text-sm">Total Storage</div>
          </div>
          <div>
            <div className="text-2xl font-bold text-white">
              {formatBytes(totalFree)}
            </div>
            <div className="text-slate-400 text-sm">Available Space</div>
          </div>
        </div>
      </div>
    </motion.div>
  );
}); 