// frontend/src/stores/index.ts
import { writable, derived } from 'svelte/store';
import type { Service, Incident, Settings, StatusCounts } from '../types';

export const services = writable<Service[]>([]);
export const incidents = writable<Incident[]>([]);
export const resolvedIncidents = writable<Incident[]>([]);
export const selectedIncident = writable<Incident | null>(null);
export const searchQuery = writable<string>('');
export const activeTab = writable<'open' | 'resolved'>('open');
export const filterType = writable<string>('');

export const settings = writable<Settings>({
  refreshInterval: 3000,
  soundEnabled: true,
  desktopNotifications: true,
  theme: 'dark',
  compactView: false
});

export const statusCounts = derived(
  [incidents, resolvedIncidents, services],
  ([$incidents, $resolvedIncidents, $services]) => {
    const activeServices = $services.filter(s => s.active).map(s => s.name);
    const filtered = $incidents.filter(i => activeServices.includes(i.service));
    
    return {
      triggered: filtered.filter(i => i.status === 'triggered').length,
      acknowledged: filtered.filter(i => i.status === 'acknowledged').length,
      resolved: $resolvedIncidents.filter(i => activeServices.includes(i.service)).length
    };
  }
);
