import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Box, Button, Container, TextField, Typography } from '@mui/material';
import { seasonService } from '../services/api';

const CreateSeason = () => {
    const { leagueId } = useParams<{ leagueId: string }>();
    const navigate = useNavigate();
    const [formData, setFormData] = useState({
        name: '',
        countingGames: 0,
        hasFinals: false,
    });
    const [error, setError] = useState('');

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value, type } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: type === 'number' ? Number(value) : value,
        }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!leagueId) return;

        try {
            await seasonService.createSeason(
                Number(leagueId),
                formData.name,
                formData.countingGames,
                formData.hasFinals
            );
            navigate(`/leagues/${leagueId}`);
        } catch (err) {
            setError('Failed to create season');
        }
    };

    return (
        <Container maxWidth="sm">
            <Box sx={{ mt: 4 }}>
                <Typography variant="h4" component="h1" gutterBottom>
                    Create New Season
                </Typography>
                {error && (
                    <Typography color="error" sx={{ mb: 2 }}>
                        {error}
                    </Typography>
                )}
                <form onSubmit={handleSubmit}>
                    <TextField
                        fullWidth
                        label="Season Name"
                        name="name"
                        value={formData.name}
                        onChange={handleChange}
                        margin="normal"
                        required
                    />
                    <TextField
                        fullWidth
                        label="Counting Games"
                        name="countingGames"
                        type="number"
                        value={formData.countingGames}
                        onChange={handleChange}
                        margin="normal"
                        required
                    />
                    <Button
                        type="submit"
                        variant="contained"
                        color="primary"
                        fullWidth
                        sx={{ mt: 3 }}
                    >
                        Create Season
                    </Button>
                </form>
            </Box>
        </Container>
    );
};

export default CreateSeason; 