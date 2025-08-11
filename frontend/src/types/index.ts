// frontend/src/types/index.ts
export interface Service {
  id: string;
  name: string;
  abbreviation: string;
  incidentCount: number;
  active: boolean;
}

export interface Incident {
  id: string;
  title: string;
  service: string;
  status: 'triggered' | 'acknowledged' | 'resolved';
  urgency: 'high' | 'low';
  createdAt: Date;
  alertCount: number;
  escalatesIn?: number;
  description?: string;
  notes?: Note[];
  pinned?: boolean; // Added for pin functionality
}

export interface Note {
  author: string;
  content: string;
  timestamp: Date;
}

export interface StatusCounts {
  triggered: number;
  acknowledged: number;
  resolved: number;
}

export interface Settings {
  refreshInterval: number;
  soundEnabled: boolean;
  desktopNotifications: boolean;
  theme: 'dark' | 'light';
  compactView: boolean;
}