import { useState, useCallback, useMemo } from "react";
import { NOTIFICATION_CONFIG } from "@/config/constants";

interface Notification {
  id: string;
  message: string;
  type: "success" | "error" | "warning" | "info";
  duration?: number;
  timestamp: number;
}

export const useNotifications = () => {
  const [notifications, setNotifications] = useState<Notification[]>([]);

  const removeNotification = useCallback((id: string) => {
    setNotifications((prev) =>
      prev.filter((notification) => notification.id !== id),
    );
  }, []);

  const addNotification = useCallback(
    (
      message: string,
      type: Notification["type"] = "info",
      duration?: number,
    ) => {
      const id = Math.random().toString(36).substr(2, 9);
      const timestamp = Date.now();
      const notification: Notification = {
        id,
        message,
        type,
        duration,
        timestamp,
      };

      setNotifications((prev) => [...prev, notification]);

      // Auto-remove notification after duration
      if (duration !== 0) {
        setTimeout(() => {
          removeNotification(id);
        }, duration || NOTIFICATION_CONFIG.DEFAULT_DURATION);
      }
    },
    [removeNotification],
  );

  const clearAll = useCallback(() => {
    setNotifications([]);
  }, []);

  // Memoize notifications to prevent unnecessary re-renders
  const memoizedNotifications = useMemo(() => notifications, [notifications]);

  return {
    notifications: memoizedNotifications,
    addNotification,
    removeNotification,
    clearAll,
  };
};
