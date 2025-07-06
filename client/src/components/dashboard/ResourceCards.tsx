import React, { useMemo, useCallback } from 'react';
import { Memory, CPU } from '@/types/monitor';
import { formatBytes, getUsageColor, getLoadStatus } from '@/utils/formatters';
import Icon from '@/components/ui/Icon';
import { motion } from 'framer-motion';

interface ResourceCardsProps {
  memory: Memory;
  cpu: CPU[];
  className?: string;
}

// Memoize individual card component to prevent unnecessary re-renders
const ResourceCard = React.memo(({ card, index }: { card: any; index: number }) => (
  <motion.div 
    className="bg-slate-800/50 rounded-xl p-6 border border-slate-700/50 hover:border-purple-500/50"
    initial={{ opacity: 0, y: 20, scale: 0.95 }}
    animate={{ opacity: 1, y: 0, scale: 1 }}
    transition={{
      delay: index * 0.05,
      duration: 0.3,
      ease: "easeOut",
    }}
    whileHover={{
      y: -2,
      scale: 1.02,
      transition: { duration: 0.2 }
    }}
  >
    <div className="flex items-center justify-between mb-4">
      <div className={`p-2 bg-${card.color}-600/20 rounded-lg`}>
        <Icon name={card.icon} className={`text-${card.color}-400`} />
      </div>
      <span className="text-slate-400 text-sm">{card.subtitle}</span>
    </div>
    <h3 className="text-white font-semibold mb-2">{card.title}</h3>
    <div className="text-3xl font-bold text-white mb-2">{card.value}</div>
    <div className="w-full bg-slate-700 rounded-full h-2 mb-2">
      <div 
        className={`bg-gradient-to-r ${getUsageColor(card.percentage)} h-2 rounded-full transition-all duration-500`}
        style={{ width: `${card.percentage}%` }}
      ></div>
    </div>
    <p className="text-slate-400 text-xs truncate">{card.description}</p>
  </motion.div>
));

ResourceCard.displayName = 'ResourceCard';

export default React.memo(function ResourceCards({ memory, cpu, className = '' }: ResourceCardsProps) {
  // Memoize calculations to prevent unnecessary re-renders
  const calculations = useMemo(() => {
    const memoryUsagePercent = ((memory.used / memory.total) * 100).toFixed(1);
    const swapUsagePercent = memory.swapTotal > 0 ? ((memory.swapUsed / memory.swapTotal) * 100).toFixed(1) : '0.0';
    const avgCpuUsage = cpu.reduce((acc, core) => acc + core.usage, 0) / cpu.length;
    const systemLoad = (avgCpuUsage / 100).toFixed(2);
    
    return {
      memoryUsagePercent: parseFloat(memoryUsagePercent),
      swapUsagePercent: parseFloat(swapUsagePercent),
      avgCpuUsage,
      systemLoad: parseFloat(systemLoad),
      loadStatus: getLoadStatus(avgCpuUsage / 100)
    };
  }, [memory.used, memory.total, memory.swapUsed, memory.swapTotal, cpu]);

  const cards = useMemo(() => [
    {
      title: 'CPU Usage',
      value: `${calculations.avgCpuUsage.toFixed(1)}%`,
      percentage: calculations.avgCpuUsage,
      color: 'red',
      icon: 'cpu',
      subtitle: `${cpu.length} cores`,
      description: cpu[0]?.model || 'CPU Model'
    },
    {
      title: 'Memory Usage',
      value: `${calculations.memoryUsagePercent}%`,
      percentage: calculations.memoryUsagePercent,
      color: 'blue',
      icon: 'memory',
      subtitle: 'RAM',
      description: `${formatBytes(memory.used)} / ${formatBytes(memory.total)}`
    },
    {
      title: 'Swap Usage',
      value: `${calculations.swapUsagePercent}%`,
      percentage: calculations.swapUsagePercent,
      color: 'yellow',
      icon: 'swap',
      subtitle: 'Swap',
      description: `${formatBytes(memory.swapUsed)} / ${formatBytes(memory.swapTotal)}`
    },
    {
      title: 'System Load',
      value: calculations.systemLoad,
      percentage: Math.min(calculations.avgCpuUsage, 100),
      color: 'green',
      icon: 'lightning',
      subtitle: 'Load',
      description: `${calculations.loadStatus} load`
    }
  ], [calculations, cpu.length, cpu[0]?.model, memory.used, memory.total, memory.swapUsed, memory.swapTotal]);

  const renderCard = useCallback((card: any, index: number) => (
    <ResourceCard key={card.title} card={card} index={index} />
  ), []);

  return (
    <div className={`grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 ${className}`}>
      {cards.map(renderCard)}
    </div>
  );
}); 