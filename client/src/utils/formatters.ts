/**
 * Utility functions for formatting data in the dashboard
 */

export const formatBytes = (bytes: number): string => {
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  if (bytes === 0) return '0 B';
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${sizes[i]}`;
};

// Format memory values that come from backend in MB
export const formatMemoryMB = (mb: number): string => {
  if (mb === 0) return '0 MB';
  if (mb < 1024) return `${mb.toFixed(1)} MB`;
  const gb = mb / 1024;
  if (gb < 1024) return `${gb.toFixed(1)} GB`;
  const tb = gb / 1024;
  return `${tb.toFixed(1)} TB`;
};

export const formatBytesPerSec = (bytes: number): string => {
  return `${formatBytes(bytes)}/s`;
};

export const formatUptime = (seconds: number): string => {
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

export const formatPercentage = (value: number, total: number): string => {
  return ((value / total) * 100).toFixed(1);
};

export const getUsageColor = (percent: number): string => {
  if (percent >= 90) return 'from-red-500 to-red-600';
  if (percent >= 75) return 'from-yellow-500 to-yellow-600';
  if (percent >= 50) return 'from-orange-500 to-orange-600';
  return 'from-green-500 to-green-600';
};

export const getUsageStatus = (percent: number): string => {
  if (percent >= 90) return 'Critical';
  if (percent >= 75) return 'Warning';
  if (percent >= 50) return 'Moderate';
  return 'Good';
};

export const getLoadStatus = (load: number): string => {
  if (load > 0.8) return 'High';
  if (load > 0.5) return 'Medium';
  return 'Low';
}; 