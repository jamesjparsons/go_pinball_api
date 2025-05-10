export interface User {
    id: number;
    email: string;
    firstName: string;
    lastName: string;
}

export interface League {
    ID: number;
    name: string;
    location: string;
    CreatedAt: string;
    UpdatedAt: string;
    DeletedAt: string | null;
    owner: User;
}

export interface Season {
    ID: number;
    name: string;
    CreatedAt: string;
    leagueID: number;
    league: League;
    countingGames: number;
    eventCount: number;
    hasFinals: boolean;
    pointDistribution: {
        [key: string]: number[];
    };
}

export interface AuthResponse {
    user: User;
    token: string;
}

export interface LoginCredentials {
    email: string;
    password: string;
}

export interface SignupCredentials extends LoginCredentials {
    firstName: string;
    lastName: string;
}

export interface CreateLeagueRequest {
    name: string;
    location: string;
}

export interface CreateSeasonRequest {
    name: string;
} 