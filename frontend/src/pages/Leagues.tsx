import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Box, Button, Container, Paper, Typography } from '@mui/material';
import { leagueService } from '../services/api';
import { League } from '../types';

const Leagues = () => {
    const navigate = useNavigate();
    const [leagues, setLeagues] = useState<League[]>([]);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchLeagues = async () => {
            try {
                const data = await leagueService.listLeagues();
                console.log('Fetched leagues:', data); // Debug log
                setLeagues(data);
            } catch (err) {
                console.error('Error fetching leagues:', err); // Debug log
                setError('Failed to fetch leagues');
            }
        };
        fetchLeagues();
    }, []);

    const handleLeagueClick = (league: League) => {
        console.log('Clicked league:', league); // Debug log
        if (league) {
            navigate(`/leagues/${league.id}`);
        } else {
            console.error('Invalid league data:', league); // Debug log
        }
    };

    return (
        <Container maxWidth="lg">
            <Box sx={{ mt: 4, mb: 4 }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
                    <Typography variant="h4" component="h1">
                        Pinball Leagues
                    </Typography>
                    <Button
                        variant="contained"
                        color="primary"
                        onClick={() => navigate('/leagues/create')}
                    >
                        Create League
                    </Button>
                </Box>
                {error && (
                    <Typography color="error" sx={{ mt: 2 }}>
                        {error}
                    </Typography>
                )}
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
                    {leagues.map((league) => (
                        <Paper
                            key={league.id}
                            sx={{
                                p: 2,
                                display: 'flex',
                                flexDirection: 'column',
                                cursor: 'pointer',
                                '&:hover': {
                                    bgcolor: 'action.hover',
                                },
                            }}
                            onClick={() => handleLeagueClick(league)}
                        >
                            <Typography variant="h6" component="h2">
                                {league.name}
                            </Typography>
                            <Typography color="text.secondary" sx={{ mt: 1 }}>
                                Location: {league.location}
                            </Typography>
                            <Typography color="text.secondary" sx={{ mt: 1 }}>
                                Created: {new Date(league.dateCreated).toLocaleDateString()}
                            </Typography>
                            <Typography color="text.secondary" sx={{ mt: 1 }}>
                                Owner: {league.owner.firstName} {league.owner.lastName}
                            </Typography>
                        </Paper>
                    ))}
                </Box>
            </Box>
        </Container>
    );
};

export default Leagues; 