import React from 'react';

interface ConnectionStatusProps {
  isConnected: boolean;
  error?: string | null;
  className?: string;
}

export default React.memo(function ConnectionStatus({ isConnected, error, className = '' }: ConnectionStatusProps) {
  return (
    <div className={`flex items-center gap-2 ${className}`}>
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
  );
}); 