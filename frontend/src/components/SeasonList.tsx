import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Box, Button, Paper, Typography } from '@mui/material';
import { seasonService } from '../services/api';
import { Season } from '../types';

interface SeasonListProps {
    leagueId: number;
}

const SeasonList = ({ leagueId }: SeasonListProps) => {
    const navigate = useNavigate();
    const [seasons, setSeasons] = useState<Season[]>([]);
    const [error, setError] = useState('');

    useEffect(() => {
        const fetchSeasons = async () => {
            if (!leagueId) return;
            
            try {
                const data = await seasonService.listSeasons(leagueId);
                setSeasons(data);
            } catch (err) {
                setError('Failed to fetch seasons');
            }
        };
        fetchSeasons();
    }, [leagueId]);

    if (!leagueId) {
        return null;
    }

    return (
        <Box sx={{ mt: 4 }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                <Typography variant="h5" component="h2">
                    Seasons
                </Typography>
                <Button
                    variant="contained"
                    color="primary"
                    onClick={() => navigate(`/leagues/${leagueId}/seasons/create`)}
                >
                    Create Season
                </Button>
            </Box>
            {error && (
                <Typography color="error" sx={{ mt: 2 }}>
                    {error}
                </Typography>
            )}
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                {seasons.map((season) => (
                    <Paper
                        key={season.ID}
                        sx={{
                            p: 2,
                            display: 'flex',
                            flexDirection: 'column',
                            cursor: 'pointer',
                            '&:hover': {
                                bgcolor: 'action.hover',
                            },
                        }}
                        onClick={() => navigate(`/leagues/${leagueId}/seasons/${season.ID}`)}
                    >
                        <Typography variant="h6" component="h3">
                            {season.name}
                        </Typography>
                        <Typography color="text.secondary" sx={{ mt: 1 }}>
                            Created: {new Date(season.CreatedAt).toLocaleDateString()}
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
    );
};

export default SeasonList; 