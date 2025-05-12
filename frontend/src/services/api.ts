import axios from 'axios';
import { User, League, Season, Event } from '../types';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

// Create axios instance with default config
const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add request interceptor to add auth token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Auth service
export const authService = {
  async login(email: string, password: string) {
    const response = await api.post('/auth/login', { email, password });
    return response.data;
  },

  async signup(email: string, password: string, firstName: string, lastName: string) {
    const response = await api.post('/auth/signup', { email, password, firstName, lastName });
    return response.data;
  },

  async getCurrentUser() {
    const response = await api.get('/auth/me');
    return response.data;
  },
};

// League service
export const leagueService = {
  async createLeague(name: string, location: string) {
    const response = await api.post('/leagues/create', { name, location });
    return response.data;
  },

  async listLeagues() {
    const response = await api.get('/leagues');
    return response.data.data;
  },

  async getLeague(id: number) {
    const response = await api.get(`/leagues/${id}`);
    return response.data.data;
  },

  async listPlayers(leagueId: number) {
    const response = await api.get(`/leagues/${leagueId}/players`);
    return response.data.data;
  },

  async addPlayersByIFPA(leagueId: number, ifpaNumbers: number[]) {
    const response = await api.post(`/leagues/${leagueId}/players/ifpa`, { ifpaNumbers });
    return response.data.data;
  },
};

// Season service
export const seasonService = {
  async createSeason(leagueId: number, name: string, countingGames: number, hasFinals: boolean) {
    const response = await api.post(`/leagues/${leagueId}/seasons/create`, {
      name,
      countingGames,
      hasFinals,
    });
    return response.data;
  },

  async listSeasons(leagueId: number) {
    const response = await api.get(`/leagues/${leagueId}/seasons`);
    return response.data.data;
  },

  async getSeason(seasonId: number) {
    const response = await api.get(`/seasons/${seasonId}`);
    return response.data.data;
  },
};

// Event service
export const eventService = {
  async createEvent(seasonId: number, name: string, date: string, location: string) {
    const response = await api.post(`/seasons/${seasonId}/events/create`, {
      name,
      date,
      location,
    });
    return response.data;
  },

  async listEvents(seasonId: number) {
    const response = await api.get(`/seasons/${seasonId}/events`);
    return response.data.data;
  },

  async getEvent(eventId: number) {
    const response = await api.get(`/events/${eventId}`);
    return response.data.data;
  },
}; 