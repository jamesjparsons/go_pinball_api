export interface User {
    id: number;
    email: string;
    firstName: string;
    lastName: string;
}

export interface League {
    id: number;
    name: string;
    location: string;
    dateCreated: string;
    owner: User;
}

export interface Season {
    id: number;
    name: string;
    dateCreated: string;
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