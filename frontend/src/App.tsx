import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider, CssBaseline } from '@mui/material';
import { AuthProvider } from './context/AuthContext';
import Layout from './components/Layout';
import Login from './pages/Login';
import Signup from './pages/Signup';
import Leagues from './pages/Leagues';
import LeagueDetail from './pages/LeagueDetail';
import CreateLeague from './pages/CreateLeague';
import SeasonDetail from './pages/SeasonDetail';
import CreateSeason from './pages/CreateSeason';
import ProtectedRoute from './components/ProtectedRoute';
import theme from './theme';

function App() {
    return (
        <ThemeProvider theme={theme}>
            <CssBaseline />
            <AuthProvider>
                <Router>
                    <Layout>
                        <Routes>
                            <Route path="/login" element={<Login />} />
                            <Route path="/signup" element={<Signup />} />
                            <Route
                                path="/leagues"
                                element={
                                    <ProtectedRoute>
                                        <Leagues />
                                    </ProtectedRoute>
                                }
                            />
                            <Route
                                path="/leagues/:leagueId"
                                element={
                                    <ProtectedRoute>
                                        <LeagueDetail />
                                    </ProtectedRoute>
                                }
                            />
                            <Route
                                path="/leagues/create"
                                element={
                                    <ProtectedRoute>
                                        <CreateLeague />
                                    </ProtectedRoute>
                                }
                            />
                            <Route
                                path="/leagues/:leagueId/seasons/:seasonId"
                                element={
                                    <ProtectedRoute>
                                        <SeasonDetail />
                                    </ProtectedRoute>
                                }
                            />
                            <Route
                                path="/leagues/:leagueId/seasons/create"
                                element={
                                    <ProtectedRoute>
                                        <CreateSeason />
                                    </ProtectedRoute>
                                }
                            />
                            <Route path="/" element={<Navigate to="/leagues" replace />} />
                        </Routes>
                    </Layout>
                </Router>
            </AuthProvider>
        </ThemeProvider>
    );
}

export default App;
