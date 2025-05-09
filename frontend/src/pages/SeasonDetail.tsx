import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Box, Container, Paper, Typography } from '@mui/material';
import { seasonService } from '../services/api';
import { Season } from '../types';

const SeasonDetail = () => {
    const { leagueId, seasonId } = useParams<{ leagueId: string; seasonId: string }>();
    const navigate = useNavigate();
    const [season, setSeason] = useState<Season | null>(null);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchSeason = async () => {
            try {
                const seasons = await seasonService.listSeasons(Number(leagueId));
                const foundSeason = seasons.find(s => s.id === Number(seasonId));
                if (foundSeason) {
                    setSeason(foundSeason);
                } else {
                    setError('Season not found');
                }
            } catch (err) {
                setError('Failed to fetch season details');
            }
        };
        fetchSeason();
    }, [leagueId, seasonId]);

    if (error) {
        return (
            <Container maxWidth="lg">
                <Typography color="error" sx={{ mt: 4 }}>
                    {error}
                </Typography>
            </Container>
        );
    }

    if (!season) {
        return (
            <Container maxWidth="lg">
                <Typography sx={{ mt: 4 }}>Loading...</Typography>
            </Container>
        );
    }

    return (
        <Container maxWidth="lg">
            <Box sx={{ mt: 4, mb: 4 }}>
                <Typography variant="h4" component="h1" gutterBottom>
                    {season.name}
                </Typography>
                <Paper sx={{ p: 3, mb: 3 }}>
                    <Typography variant="h6" gutterBottom>
                        Season Details
                    </Typography>
                    <Typography>
                        Created: {new Date(season.dateCreated).toLocaleDateString()}
                    </Typography>
                    <Typography>
                        Events: {season.eventCount}
                    </Typography>
                    <Typography>
                        Counting Games: {season.countingGames}
                    </Typography>
                    <Typography>
                        Has Finals: {season.hasFinals ? 'Yes' : 'No'}
                    </Typography>
                </Paper>

                <Paper sx={{ p: 3 }}>
                    <Typography variant="h6" gutterBottom>
                        Point Distribution
                    </Typography>
                    {Object.entries(season.pointDistribution).map(([playerCount, points]) => (
                        <Box key={playerCount} sx={{ mb: 2 }}>
                            <Typography variant="subtitle1">
                                {playerCount} Players
                            </Typography>
                            <Typography>
                                Points: {points.join(', ')}
                            </Typography>
                        </Box>
                    ))}
                </Paper>
            </Box>
        </Container>
    );
};

export default SeasonDetail; 