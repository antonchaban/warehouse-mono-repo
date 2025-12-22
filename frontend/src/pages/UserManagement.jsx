import { useEffect, useState } from 'react';
import { Users, Save, Shield } from 'lucide-react';
import api from '../api';

export default function UserManagement() {
    const [users, setUsers] = useState([]);
    const [roles] = useState(["ROLE_STOREKEEPER", "ROLE_LOGISTICIAN", "ROLE_ADMIN"]);

    useEffect(() => {
        fetchUsers();
    }, []);

    const fetchUsers = async () => {
        try {
            const res = await api.get('/admin/users');
            setUsers(res.data);
        } catch (err) {
            console.error("Failed to fetch users", err);
        }
    };

    const handleRoleChange = async (userId, newRole) => {
        try {
            await api.put(`/admin/users/${userId}/role`, { roleName: newRole });
            alert(`Role updated to ${newRole}`);
            fetchUsers(); // Refresh list
        } catch (err) {
            console.error(err);
            alert("Failed to update role");
        }
    };

    return (
        <div className="min-h-screen bg-slate-50 p-8">
            <div className="max-w-6xl mx-auto bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
                <div className="px-6 py-4 border-b border-slate-100 bg-slate-50 flex items-center gap-2">
                    <Users className="text-blue-600" />
                    <h2 className="text-lg font-bold text-slate-700">User Management Console</h2>
                </div>

                <div className="overflow-x-auto">
                    <table className="w-full text-left text-sm">
                        <thead className="bg-slate-50 text-slate-500 uppercase font-bold text-xs">
                        <tr>
                            <th className="px-6 py-3">ID</th>
                            <th className="px-6 py-3">Username</th>
                            <th className="px-6 py-3">Email</th>
                            <th className="px-6 py-3">Current Role</th>
                            <th className="px-6 py-3">Action</th>
                        </tr>
                        </thead>
                        <tbody className="divide-y divide-slate-100">
                        {users.map(user => (
                            <tr key={user.id} className="hover:bg-slate-50 transition">
                                <td className="px-6 py-4 font-mono text-slate-500">#{user.id}</td>
                                <td className="px-6 py-4 font-bold text-slate-700">{user.username}</td>
                                <td className="px-6 py-4 text-slate-600">{user.email}</td>
                                <td className="px-6 py-4">
                                    {/* Show current role as badge */}
                                    <span className={`px-2 py-1 rounded text-xs font-bold border ${
                                        user.roles[0]?.name === 'ROLE_ADMIN' ? 'bg-purple-100 text-purple-700 border-purple-200' :
                                            user.roles[0]?.name === 'ROLE_LOGISTICIAN' ? 'bg-blue-100 text-blue-700 border-blue-200' :
                                                'bg-green-100 text-green-700 border-green-200'
                                    }`}>
                                            {user.roles[0]?.name || 'NO_ROLE'}
                                        </span>
                                </td>
                                <td className="px-6 py-4">
                                    <select
                                        className="border border-slate-300 rounded px-2 py-1 text-xs focus:ring-2 focus:ring-blue-500 outline-none"
                                        defaultValue={user.roles[0]?.name}
                                        onChange={(e) => {
                                            if(window.confirm(`Change role for ${user.username} to ${e.target.value}?`)) {
                                                handleRoleChange(user.id, e.target.value);
                                            } else {
                                                e.target.value = user.roles[0]?.name;
                                            }
                                        }}
                                    >
                                        {roles.map(r => (
                                            <option key={r} value={r}>{r.replace('ROLE_', '')}</option>
                                        ))}
                                    </select>
                                </td>
                            </tr>
                        ))}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
}