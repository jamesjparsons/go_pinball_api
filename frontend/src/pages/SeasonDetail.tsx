import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Box, Container, Paper, Typography } from '@mui/material';
import { seasonService } from '../services/api';
import { Season } from '../types';

const SeasonDetail = () => {
    const { leagueId, seasonId } = useParams<{ leagueId: string; seasonId: string }>();
    const [season, setSeason] = useState<Season | null>(null);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchSeason = async () => {
            if (!leagueId || !seasonId) return;

            try {
                const seasons = await seasonService.listSeasons(Number(leagueId));
                const foundSeason = seasons.find((s: Season) => s.ID === Number(seasonId));
                if (foundSeason) {
                    setSeason(foundSeason);
                } else {
                    setError('Season not found');
                }
            } catch (err) {
                setError('Failed to fetch season');
            }
        };

        fetchSeason();
    }, [leagueId, seasonId]);

    if (error) {
        return (
            <Container>
                <Typography color="error">{error}</Typography>
            </Container>
        );
    }

    if (!season) {
        return (
            <Container>
                <Typography>Loading...</Typography>
            </Container>
        );
    }

    return (
        <Container>
            <Box sx={{ mt: 4 }}>
                <Paper sx={{ p: 3 }}>
                    <Typography variant="h4" component="h1" gutterBottom>
                        {season.name}
                    </Typography>
                    <Typography variant="body1">
                        Counting Games: {season.countingGames}
                    </Typography>
                    <Typography variant="body1">
                        Has Finals: {season.hasFinals ? 'Yes' : 'No'}
                    </Typography>
                    <Typography variant="body1">
                        Event Count: {season.eventCount}
                    </Typography>
                </Paper>
            </Box>
        </Container>
    );
};

export default SeasonDetail; 