// Connection settings
export const CONNECTION_CONFIG = {
  MAX_RECONNECT_ATTEMPTS: 5,
  BASE_RECONNECT_DELAY: 1000, // 1 second
  EXPONENTIAL_BACKOFF: true,
} as const;

// Notification settings
export const NOTIFICATION_CONFIG = {
  DEFAULT_DURATION: 5000, // 5 seconds
  SUCCESS_DURATION: 3000, // 3 seconds
  ERROR_DURATION: 0, // Don't auto-dismiss
} as const;

// Animation delays
export const ANIMATION_DELAYS = {
  CARD_STAGGER: 50, // Reduced from 100ms for better performance
  FADE_IN_DELAY: 100, // Reduced from 200ms
  FADE_IN_DELAY_2: 200, // Reduced from 400ms
} as const;

// Performance thresholds
export const PERFORMANCE_THRESHOLDS = {
  HIGH_CPU_USAGE: 80,
  MEDIUM_CPU_USAGE: 50,
  HIGH_MEMORY_USAGE: 90,
  MEDIUM_MEMORY_USAGE: 75,
  HIGH_DISK_USAGE: 90,
  MEDIUM_DISK_USAGE: 75,
} as const;

// Performance optimization settings
export const PERFORMANCE_CONFIG = {
  DEBOUNCE_DELAY: 100, // ms for data updates
  ANIMATION_DURATION: 200, // ms for transitions
  MAX_ANIMATIONS: 10, // Maximum concurrent animations
  REDUCED_MOTION: false, // Respect user preferences
} as const;

// Default endpoints
export const ENDPOINTS = {
  DEFAULT_SSE: 'http://localhost:8080/api/stats/sse',
} as const; 