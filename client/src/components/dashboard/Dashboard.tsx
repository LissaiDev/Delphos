"use client";

import React, { useCallback, useMemo } from "react";
import { Monitor } from "@/types/monitor";
import { useMonitorData } from "@/hooks/useMonitorData";
import { useNotifications } from "@/hooks/useNotifications";
import SystemOverview from "@/components/dashboard/SystemOverview";
import ResourceCards from "@/components/dashboard/ResourceCards";
import NetworkStats from "@/components/dashboard/NetworkStats";
import DiskUsage from "@/components/dashboard/DiskUsage";
import Notification from "@/components/ui/Notification";
import LoadingSpinner from "@/components/ui/LoadingSpinner";
import Icon from "@/components/ui/Icon";
import ConnectionStatus from "@/components/ui/ConnectionStatus";
import { NOTIFICATION_CONFIG } from "@/config/constants";

interface DashboardProps {
  endpoint?: string;
  data?: Monitor;
  isLoading?: boolean;
}

export default function Dashboard({
  endpoint,
  data,
  isLoading: externalLoading,
}: DashboardProps) {
  const {
    data: fetchedData,
    isLoading,
    error,
    isConnected,
    refresh,
    reconnect,
  } = useMonitorData(endpoint);
  const { notifications, addNotification, removeNotification } =
    useNotifications();

  // Memoize current data to prevent unnecessary re-renders
  const currentData = useMemo(() => {
    return data || fetchedData;
  }, [data, fetchedData]);

  const isCurrentlyLoading = useMemo(() => {
    return externalLoading !== undefined ? externalLoading : isLoading;
  }, [externalLoading, isLoading]);

  // Show notifications for connection status changes
  React.useEffect(() => {
    if (isConnected) {
      addNotification(
        "Connected to server successfully",
        "success",
        NOTIFICATION_CONFIG.SUCCESS_DURATION,
      );
    }
  }, [isConnected, addNotification]);

  React.useEffect(() => {
    if (error && !isConnected) {
      addNotification(error, "error", NOTIFICATION_CONFIG.ERROR_DURATION);
    }
  }, [error, isConnected, addNotification]);

  const handleRefresh = useCallback(() => {
    refresh();
  }, [refresh]);

  const handleReconnect = useCallback(() => {
    reconnect();
  }, [reconnect]);

  // Show loading state if no data is available
  if (isCurrentlyLoading || !currentData) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 p-4">
        <div className="max-w-7xl mx-auto">
          <div className="flex flex-col items-center justify-center min-h-[60vh]">
            <LoadingSpinner size="lg" className="mb-6" />
            <h2 className="text-2xl font-semibold text-white mb-2 animate-fade-in">
              Connecting to Server...
            </h2>
            <p className="text-slate-400 text-center max-w-md animate-fade-in-delay">
              Establishing real-time connection to monitor system resources
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 p-4">
      {/* Notifications */}
      {notifications.map((notification) => (
        <Notification
          key={notification.id}
          message={notification.message}
          type={notification.type}
          duration={notification.duration}
          onClose={() => removeNotification(notification.id)}
        />
      ))}

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
            {/* Connection Status */}
            <ConnectionStatus
              isConnected={isConnected}
              error={error}
              className="mt-2 animate-fade-in-delay"
            />
          </div>
          <div className="flex gap-2 mt-4 sm:mt-0">
            <button
              onClick={handleReconnect}
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-all duration-300 hover:scale-105 animate-fade-in-delay-2"
            >
              <span className="flex items-center gap-2">
                <Icon name="reconnect" />
                Reconnect
              </span>
            </button>
            <button
              onClick={handleRefresh}
              className="px-6 py-2 bg-purple-600 hover:bg-purple-700 text-white rounded-lg transition-all duration-300 hover:scale-105 animate-fade-in-delay-2"
            >
              <span className="flex items-center gap-2">
                <Icon name="refresh" />
                Refresh
              </span>
            </button>
          </div>
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
