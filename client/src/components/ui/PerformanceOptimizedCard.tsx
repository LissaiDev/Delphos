import React from 'react';

interface PerformanceOptimizedCardProps {
  children: React.ReactNode;
  className?: string;
  animate?: boolean;
  index?: number;
}

export default React.memo(function PerformanceOptimizedCard({ 
  children, 
  className = '', 
  animate = true,
  index = 0 
}: PerformanceOptimizedCardProps) {
  const animationClass = animate ? 'animate-slide-up' : '';
  const animationDelay = animate ? { animationDelay: `${index * 50}ms` } : {};

  return (
    <div 
      className={`bg-slate-800/50 rounded-xl p-6 border border-slate-700/50 hover:border-purple-500/50 transition-all duration-200 ${animationClass} ${className}`}
      style={animationDelay}
    >
      {children}
    </div>
  );
}); 