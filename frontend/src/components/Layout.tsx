import React from 'react';
import { AppBar, Box, Button, Container, Toolbar, Typography } from '@mui/material';
import { Link as RouterLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

interface LayoutProps {
    children: React.ReactNode;
}

const Layout = ({ children }: LayoutProps) => {
    const { user, logout } = useAuth();
    const navigate = useNavigate();

    const handleLogout = () => {
        logout();
        navigate('/login');
    };

    return (
        <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
            <AppBar position="static">
                <Toolbar>
                    <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
                        <RouterLink to="/" style={{ color: 'white', textDecoration: 'none' }}>
                            Pinball League Manager
                        </RouterLink>
                    </Typography>
                    {user ? (
                        <>
                            <Button
                                color="inherit"
                                component={RouterLink}
                                to="/leagues"
                                sx={{ mr: 2 }}
                            >
                                Leagues
                            </Button>
                            <Typography sx={{ mr: 2 }}>
                                {user.email}
                            </Typography>
                            <Button color="inherit" onClick={handleLogout}>
                                Logout
                            </Button>
                        </>
                    ) : (
                        <>
                            <Button color="inherit" component={RouterLink} to="/login">
                                Login
                            </Button>
                            <Button color="inherit" component={RouterLink} to="/signup">
                                Sign Up
                            </Button>
                        </>
                    )}
                </Toolbar>
            </AppBar>
            <Container component="main" sx={{ mt: 4, mb: 4, flex: 1 }}>
                {children}
            </Container>
        </Box>
    );
};

export default Layout; 