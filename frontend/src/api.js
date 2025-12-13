import axios from 'axios';

const api = axios.create({
    baseURL: 'http://localhost:8080/api', // Звертаємось до Java
});

// Додаємо токен (якщо він є) до кожного запиту
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

export default api;