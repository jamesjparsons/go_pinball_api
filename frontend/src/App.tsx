import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import theme from './theme';
import { AuthProvider } from './context/AuthContext';

// Pages
import Leagues from './pages/Leagues';
import LeagueDetail from './pages/LeagueDetail';
import CreateLeague from './pages/CreateLeague';
import CreateSeason from './pages/CreateSeason';
import SeasonDetail from './pages/SeasonDetail';
import Login from './pages/Login';
import Signup from './pages/Signup';

// Components
import PrivateRoute from './components/PrivateRoute';
import Layout from './components/Layout';

const App: React.FC = () => {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <AuthProvider>
        <Router>
          <Layout>
            <Routes>
              <Route path="/login" element={<Login />} />
              <Route path="/signup" element={<Signup />} />
              <Route path="/" element={<Navigate to="/leagues" replace />} />
              <Route
                path="/leagues"
                element={
                  <PrivateRoute>
                    <Leagues />
                  </PrivateRoute>
                }
              />
              <Route
                path="/leagues/create"
                element={
                  <PrivateRoute>
                    <CreateLeague />
                  </PrivateRoute>
                }
              />
              <Route
                path="/leagues/:leagueId"
                element={
                  <PrivateRoute>
                    <LeagueDetail />
                  </PrivateRoute>
                }
              />
              <Route
                path="/leagues/:leagueId/seasons/create"
                element={
                  <PrivateRoute>
                    <CreateSeason />
                  </PrivateRoute>
                }
              />
              <Route
                path="/leagues/:leagueId/seasons/:seasonId"
                element={
                  <PrivateRoute>
                    <SeasonDetail />
                  </PrivateRoute>
                }
              />
            </Routes>
          </Layout>
        </Router>
      </AuthProvider>
    </ThemeProvider>
  );
};

export default App; 