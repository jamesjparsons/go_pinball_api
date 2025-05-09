import { AppBar, Toolbar, Typography, Button, Box } from '@mui/material';
import { Link as RouterLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

interface LayoutProps {
    children: React.ReactNode;
}

const Layout = ({ children }: LayoutProps) => {
    const { isAuthenticated, logout } = useAuth();
    const navigate = useNavigate();

    const handleLogout = () => {
        logout();
        navigate('/login');
    };

    return (
        <Box sx={{ flexGrow: 1 }}>
            <AppBar position="static">
                <Toolbar>
                    <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
                        Pinball League Manager
                    </Typography>
                    {isAuthenticated ? (
                        <>
                            <Button color="inherit" component={RouterLink} to="/leagues">
                                Leagues
                            </Button>
                            <Button color="inherit" component={RouterLink} to="/leagues/create">
                                Create League
                            </Button>
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
            <Box component="main" sx={{ p: 3 }}>
                {children}
            </Box>
        </Box>
    );
};

export default Layout; 