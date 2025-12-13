import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { UserPlus } from 'lucide-react';
import api from '../api';

export default function Register() {
    const [formData, setFormData] = useState({ username: '', password: '', email: '' });
    const [error, setError] = useState('');
    const navigate = useNavigate();

    const handleRegister = async (e) => {
        e.preventDefault();
        try {
            await api.post('/auth/register', formData);
            alert('Registration successful! Please login.');
            navigate('/');
        } catch (err) {
            setError('Registration failed. Username might be taken.');
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-slate-100">
            <div className="bg-white p-8 rounded-lg shadow-lg w-96 border border-slate-200">
                <div className="flex justify-center mb-6 text-emerald-600">
                    <UserPlus size={48} />
                </div>
                <h2 className="text-2xl font-bold text-center mb-6 text-slate-700">Create Account</h2>
                {error && <div className="bg-red-50 text-red-600 p-2 rounded mb-4 text-sm">{error}</div>}

                <form onSubmit={handleRegister} className="space-y-4">
                    <input
                        type="text" placeholder="Username"
                        className="w-full p-2 border rounded"
                        onChange={(e) => setFormData({...formData, username: e.target.value})}
                    />
                    <input
                        type="email" placeholder="Email (optional)"
                        className="w-full p-2 border rounded"
                        onChange={(e) => setFormData({...formData, email: e.target.value})}
                    />
                    <input
                        type="password" placeholder="Password"
                        className="w-full p-2 border rounded"
                        onChange={(e) => setFormData({...formData, password: e.target.value})}
                    />
                    <button type="submit" className="w-full bg-emerald-600 text-white py-2 rounded hover:bg-emerald-700">
                        Register
                    </button>
                </form>
                <div className="mt-4 text-center text-sm">
                    <Link to="/" className="text-blue-600 hover:underline">Back to Login</Link>
                </div>
            </div>
        </div>
    );
}