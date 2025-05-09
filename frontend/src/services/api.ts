import axios from 'axios';
import { AuthResponse, CreateLeagueRequest, League, LoginCredentials, SignupCredentials, User, Season, CreateSeasonRequest } from '../types';

const API_URL = 'http://localhost:8080/api';

const api = axios.create({
    baseURL: API_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Add token to requests if it exists
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

export const authService = {
    login: async (credentials: LoginCredentials): Promise<AuthResponse> => {
        const response = await api.post<AuthResponse>('/auth/login', credentials);
        localStorage.setItem('token', response.data.token);
        return response.data;
    },

    signup: async (credentials: SignupCredentials): Promise<AuthResponse> => {
        const response = await api.post<AuthResponse>('/auth/signup', credentials);
        localStorage.setItem('token', response.data.token);
        return response.data;
    },

    logout: () => {
        localStorage.removeItem('token');
    },

    getCurrentUser: async (): Promise<User | null> => {
        try {
            const response = await api.get<User>('/auth/me');
            return response.data;
        } catch (err) {
            return null;
        }
    },
};

export const leagueService = {
    listLeagues: async (): Promise<League[]> => {
        const response = await api.get<{ data: League[] }>('/leagues');
        return response.data.data;
    },

    getLeague: async (leagueId: number): Promise<League> => {
        const response = await api.get<{ data: League }>(`/leagues/${leagueId}`);
        return response.data.data;
    },

    createLeague: async (league: CreateLeagueRequest): Promise<League> => {
        const response = await api.post<{ data: League }>('/leagues/create', league);
        return response.data.data;
    },
};

export const seasonService = {
    listSeasons: async (leagueId: number): Promise<Season[]> => {
        const response = await api.get<{ data: Season[] }>(`/leagues/${leagueId}/seasons`);
        return response.data.data;
    },

    getSeason: async (leagueId: number, seasonId: number): Promise<Season> => {
        const response = await api.get<{ data: Season }>(`/leagues/${leagueId}/seasons/${seasonId}`);
        return response.data.data;
    },

    createSeason: async (leagueId: number, season: CreateSeasonRequest): Promise<Season> => {
        const response = await api.post<{ data: Season }>(`/leagues/${leagueId}/seasons`, season);
        return response.data.data;
    },
}; 