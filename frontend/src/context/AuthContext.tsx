import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { User } from '../types';
import { authService } from '../services/api';

interface AuthContextType {
    user: User | null;
    login: (email: string, password: string) => Promise<void>;
    signup: (email: string, password: string, firstName: string, lastName: string) => Promise<void>;
    logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};

interface AuthProviderProps {
    children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
    const [user, setUser] = useState<User | null>(null);

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (token) {
            authService.getCurrentUser()
                .then(response => {
                    setUser(response);
                })
                .catch(() => {
                    localStorage.removeItem('token');
                });
        }
    }, []);

    const login = async (email: string, password: string) => {
        const response = await authService.login(email, password);
        localStorage.setItem('token', response.token);
        setUser(response.user);
    };

    const signup = async (email: string, password: string, firstName: string, lastName: string) => {
        const response = await authService.signup(email, password, firstName, lastName);
        localStorage.setItem('token', response.token);
        setUser(response.user);
    };

    const logout = () => {
        localStorage.removeItem('token');
        setUser(null);
    };

    return (
        <AuthContext.Provider value={{ user, login, signup, logout }}>
            {children}
        </AuthContext.Provider>
    );
}; 