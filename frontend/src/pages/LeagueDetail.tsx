import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Box, Button, Container, Paper, Typography, Divider } from '@mui/material';
import { leagueService, seasonService } from '../services/api';
import { League, Season } from '../types';

const LeagueDetail = () => {
    const { leagueId } = useParams<{ leagueId: string }>();
    const navigate = useNavigate();
    const [league, setLeague] = useState<League | null>(null);
    const [seasons, setSeasons] = useState<Season[]>([]);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchData = async () => {
            if (!leagueId) return;

            try {
                // Fetch league details
                const league = await leagueService.getLeague(Number(leagueId));
                
                if (league) {
                    setLeague(league);
                } else {
                    setError('League not found');
                }

                // Fetch seasons
                const seasonsData = await seasonService.listSeasons(Number(leagueId));
                setSeasons(seasonsData);
            } catch (err) {
                setError('Failed to fetch league data');
            }
        };

        fetchData();
    }, [leagueId]);

    if (error) {
        return (
            <Container maxWidth="lg">
                <Typography color="error" sx={{ mt: 4 }}>
                    {error}
                </Typography>
            </Container>
        );
    }

    if (!league) {
        return (
            <Container maxWidth="lg">
                <Typography sx={{ mt: 4 }}>Loading...</Typography>
            </Container>
        );
    }

    // Sort seasons by date created, most recent first
    const sortedSeasons = [...seasons].sort((a, b) => 
        new Date(b.dateCreated).getTime() - new Date(a.dateCreated).getTime()
    );

    // Get the most recent season as active
    const activeSeason = sortedSeasons[0];
    const inactiveSeasons = sortedSeasons.slice(1);

    return (
        <Container maxWidth="lg">
            <Box sx={{ mt: 4, mb: 4 }}>
                {/* League Information */}
                <Paper sx={{ p: 3, mb: 4 }}>
                    <Typography variant="h4" component="h1" gutterBottom>
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

                {/* Active Season Section */}
                <Box sx={{ mb: 4 }}>
                    <Typography variant="h5" component="h2" gutterBottom>
                        Active Season
                    </Typography>
                    {activeSeason ? (
                        <Paper
                            sx={{
                                p: 2,
                                cursor: 'pointer',
                                '&:hover': {
                                    bgcolor: 'action.hover',
                                },
                            }}
                            onClick={() => navigate(`/leagues/${leagueId}/seasons/${activeSeason.id}`)}
                        >
                            <Typography variant="h6">
                                {activeSeason.name}
                            </Typography>
                            <Typography color="text.secondary" sx={{ mt: 1 }}>
                                Created: {new Date(activeSeason.dateCreated).toLocaleDateString()}
                            </Typography>
                            <Typography color="text.secondary">
                                Events: {activeSeason.eventCount}
                            </Typography>
                            <Typography color="text.secondary">
                                Counting Games: {activeSeason.countingGames}
                            </Typography>
                            {activeSeason.hasFinals && (
                                <Typography color="primary">
                                    Has Finals
                                </Typography>
                            )}
                        </Paper>
                    ) : (
                        <Box sx={{ textAlign: 'center', py: 4 }}>
                            <Typography color="text.secondary" sx={{ mb: 2 }}>
                                No active season
                            </Typography>
                            <Button
                                variant="contained"
                                color="primary"
                                onClick={() => navigate(`/leagues/${leagueId}/seasons/create`)}
                            >
                                Create New Season
                            </Button>
                        </Box>
                    )}
                </Box>

                {/* Inactive Seasons Section */}
                {inactiveSeasons.length > 0 && (
                    <Box>
                        <Typography variant="h5" component="h2" gutterBottom>
                            Past Seasons
                        </Typography>
                        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                            {inactiveSeasons.map((season) => (
                                <Paper
                                    key={season.id}
                                    sx={{
                                        p: 2,
                                        cursor: 'pointer',
                                        '&:hover': {
                                            bgcolor: 'action.hover',
                                        },
                                    }}
                                    onClick={() => navigate(`/leagues/${leagueId}/seasons/${season.id}`)}
                                >
                                    <Typography variant="h6">
                                        {season.name}
                                    </Typography>
                                    <Typography color="text.secondary" sx={{ mt: 1 }}>
                                        Created: {new Date(season.dateCreated).toLocaleDateString()}
                                    </Typography>
                                    <Typography color="text.secondary">
                                        Events: {season.eventCount}
                                    </Typography>
                                    <Typography color="text.secondary">
                                        Counting Games: {season.countingGames}
                                    </Typography>
                                    {season.hasFinals && (
                                        <Typography color="primary">
                                            Has Finals
                                        </Typography>
                                    )}
                                </Paper>
                            ))}
                        </Box>
                    </Box>
                )}
            </Box>
        </Container>
    );
};

export default LeagueDetail; 