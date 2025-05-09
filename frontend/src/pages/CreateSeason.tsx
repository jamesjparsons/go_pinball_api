import { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Box, Button, Container, TextField, Typography } from '@mui/material';
import { seasonService } from '../services/api';

const CreateSeason = () => {
    const { leagueId } = useParams<{ leagueId: string }>();
    const navigate = useNavigate();
    const [formData, setFormData] = useState({
        name: '',
    });
    const [error, setError] = useState('');

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!leagueId) return;

        try {
            await seasonService.createSeason(Number(leagueId), formData);
            navigate(`/leagues/${leagueId}`);
        } catch (err) {
            setError('Failed to create season');
        }
    };

    return (
        <Container maxWidth="sm">
            <Box sx={{ mt: 4, mb: 4 }}>
                <Typography variant="h4" component="h1" gutterBottom>
                    Create New Season
                </Typography>
                {error && (
                    <Typography color="error" sx={{ mt: 2 }}>
                        {error}
                    </Typography>
                )}
                <Box component="form" onSubmit={handleSubmit} sx={{ mt: 3 }}>
                    <TextField
                        required
                        fullWidth
                        label="Season Name"
                        name="name"
                        value={formData.name}
                        onChange={handleChange}
                        margin="normal"
                    />
                    <Box sx={{ display: 'flex', gap: 2, mt: 3 }}>
                        <Button
                            type="submit"
                            variant="contained"
                            sx={{ flex: 1 }}
                        >
                            Create Season
                        </Button>
                        <Button
                            variant="outlined"
                            onClick={() => navigate(`/leagues/${leagueId}`)}
                            sx={{ flex: 1 }}
                        >
                            Cancel
                        </Button>
                    </Box>
                </Box>
            </Box>
        </Container>
    );
};

export default CreateSeason; 