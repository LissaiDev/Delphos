import React from 'react';
import { motion } from 'framer-motion';

interface AnimatedCardProps {
  children: React.ReactNode;
  className?: string;
  index?: number;
  hover?: boolean;
  backdrop?: boolean;
}

export default React.memo(function AnimatedCard({ 
  children, 
  className = '', 
  index = 0,
  hover = true,
  backdrop = false
}: AnimatedCardProps) {
  const baseClasses = backdrop 
    ? 'bg-slate-800/50 backdrop-blur-sm rounded-xl p-6 border border-slate-700/50'
    : 'bg-slate-800/50 rounded-xl p-6 border border-slate-700/50';
  
  const hoverClasses = hover 
    ? 'hover:border-purple-500/50 cursor-pointer'
    : '';

  return (
    <motion.div
      initial={{ opacity: 0, y: 20, scale: 0.95 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      transition={{
        delay: index * 0.05,
        duration: 0.3,
        ease: "easeOut",
      }}
      whileHover={hover ? {
        y: -2,
        scale: 1.02,
        transition: { duration: 0.2 }
      } : undefined}
      className={`${baseClasses} ${hoverClasses} ${className}`}
      style={{
        willChange: 'transform, opacity',
      }}
    >
      {children}
    </motion.div>
  );
}); 