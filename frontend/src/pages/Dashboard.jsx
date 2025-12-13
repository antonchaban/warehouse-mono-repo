import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
// 👇 1. Додав іконку PieChart для краси
import { LogOut, Play, RefreshCw, Box, ArrowRight, Truck, Plus, Warehouse, Package, PieChart } from 'lucide-react';
import api from '../api';

// 👇 2. Імпорт вашого нового компонента (перевірте шлях, якщо він в іншій папці)
import WarehouseChart from './WarehouseChart';

export default function Dashboard() {
    const [shipments, setShipments] = useState([]);
    const [supplyId, setSupplyId] = useState('555');
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();
    const [newSupplyWhId, setNewSupplyWhId] = useState('');
    const [newSupplyProdId, setNewSupplyProdId] = useState('');
    const [newSupplyQty, setNewSupplyQty] = useState('');

    // Нові стейти для створення
    const [showForms, setShowForms] = useState(false);
    const [newWhCapacity, setNewWhCapacity] = useState('');
    const [newProdVolume, setNewProdVolume] = useState('');

    const createSupply = async () => {
        try {
            const res = await api.post('/admin/supplies', {
                warehouseId: parseInt(newSupplyWhId),
                productId: parseInt(newSupplyProdId),
                quantity: parseInt(newSupplyQty)
            });

            console.log("Server response:", res.data);

            let newId = 'unknown';
            if (typeof res.data === 'string' && res.data.includes('ID:')) {
                const parts = res.data.split('ID:');
                if (parts.length > 1) {
                    newId = parts[1].trim();
                }
            } else {
                newId = '1';
            }

            alert(`Supply Created! ID from server: ${newId}`);

            if (newId !== 'unknown') {
                setSupplyId(newId);
            }

            setNewSupplyWhId('');
            setNewSupplyProdId('');
            setNewSupplyQty('');

        } catch (e) {
            console.error(e);
            alert('Supply might be created, but frontend failed to parse response. Check console.');
        }
    };

    const fetchShipments = async () => {
        try {
            const res = await api.get('/v1/distribution/shipments');
            setShipments(res.data);
        } catch (err) {
            console.error("Failed to fetch shipments", err);
        }
    };

    useEffect(() => {
        fetchShipments();
        const interval = setInterval(fetchShipments, 2000);
        return () => clearInterval(interval);
    }, []);

    const handleLogout = () => {
        localStorage.removeItem('token');
        navigate('/');
    };

    const handleCalculate = async () => {
        setLoading(true);
        try {
            await api.post('/v1/distribution/calculate', { supplyId: parseInt(supplyId) });
        } catch (err) {
            alert('Error triggering calculation. Check permissions or supply ID.');
        } finally {
            setTimeout(() => setLoading(false), 1000);
        }
    };

    const createWarehouse = async () => {
        try {
            await api.post('/admin/warehouses', { capacity: parseFloat(newWhCapacity) });
            alert('Warehouse Created!');
            setNewWhCapacity('');
        } catch (e) { alert('Error creating warehouse'); }
    };

    const createProduct = async () => {
        try {
            await api.post('/admin/products', { volume: parseFloat(newProdVolume) });
            alert('Product Created!');
            setNewProdVolume('');
        } catch (e) { alert('Error creating product'); }
    };

    return (
        <div className="min-h-screen bg-slate-50 text-slate-800">
            {/* Header */}
            <header className="bg-white shadow-sm border-b border-slate-200 sticky top-0 z-10">
                <div className="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
                    <h1 className="text-xl font-bold text-blue-600 flex items-center gap-2">
                        <Truck /> Intelligent Distribution System
                    </h1>
                    <button onClick={handleLogout} className="text-slate-500 hover:text-red-600 flex items-center gap-1 transition-colors text-sm font-medium">
                        <LogOut size={16} /> Logout
                    </button>
                </div>
            </header>

            <main className="max-w-7xl mx-auto px-4 py-8">

                {/* Кнопка перемикання режимів */}
                <div className="mb-6 flex justify-end">
                    <button
                        onClick={() => setShowForms(!showForms)}
                        className="text-sm text-blue-600 underline hover:text-blue-800"
                    >
                        {showForms ? "Hide Admin Forms" : "Create New Data (+)"}
                    </button>
                </div>

                {/* Секція створення */}
                {showForms && (
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
                        {/* Create Warehouse */}
                        <div className="bg-white p-6 rounded-xl shadow-sm border border-slate-200">
                            <h3 className="font-bold mb-4 flex items-center gap-2 text-slate-700">
                                <Warehouse size={20} className="text-purple-600"/> New Warehouse
                            </h3>
                            <div className="flex gap-2">
                                <input
                                    type="number" placeholder="Capacity (m³)"
                                    value={newWhCapacity} onChange={(e) => setNewWhCapacity(e.target.value)}
                                    className="border p-2 rounded w-full"
                                />
                                <button onClick={createWarehouse} className="bg-purple-600 text-white px-4 py-2 rounded hover:bg-purple-700">
                                    <Plus size={20}/>
                                </button>
                            </div>
                        </div>

                        {/* Create Product */}
                        <div className="bg-white p-6 rounded-xl shadow-sm border border-slate-200">
                            <h3 className="font-bold mb-4 flex items-center gap-2 text-slate-700">
                                <Package size={20} className="text-orange-600"/> New Product Type
                            </h3>
                            <div className="flex gap-2">
                                <input
                                    type="number" placeholder="Volume (m³)"
                                    value={newProdVolume} onChange={(e) => setNewProdVolume(e.target.value)}
                                    className="border p-2 rounded w-full"
                                />
                                <button onClick={createProduct} className="bg-orange-600 text-white px-4 py-2 rounded hover:bg-orange-700">
                                    <Plus size={20}/>
                                </button>
                            </div>
                        </div>

                        {/* Create Supply */}
                        <div className="bg-white p-6 rounded-xl shadow-sm border border-slate-200 md:col-span-2 bg-gradient-to-r from-blue-50 to-white">
                            <h3 className="font-bold mb-4 flex items-center gap-2 text-slate-700">
                                <Truck size={20} className="text-blue-600"/> Register Incoming Supply
                            </h3>
                            <div className="flex flex-wrap gap-4 items-end">
                                <div>
                                    <label className="block text-xs font-bold text-slate-400 mb-1">Where arrived? (Wh ID)</label>
                                    <input
                                        type="number" placeholder="1"
                                        value={newSupplyWhId} onChange={(e) => setNewSupplyWhId(e.target.value)}
                                        className="border p-2 rounded w-24"
                                    />
                                </div>
                                <div>
                                    <label className="block text-xs font-bold text-slate-400 mb-1">What product? (Prod ID)</label>
                                    <input
                                        type="number" placeholder="100"
                                        value={newSupplyProdId} onChange={(e) => setNewSupplyProdId(e.target.value)}
                                        className="border p-2 rounded w-24"
                                    />
                                </div>
                                <div>
                                    <label className="block text-xs font-bold text-slate-400 mb-1">How many?</label>
                                    <input
                                        type="number" placeholder="50"
                                        value={newSupplyQty} onChange={(e) => setNewSupplyQty(e.target.value)}
                                        className="border p-2 rounded w-24"
                                    />
                                </div>
                                <button onClick={createSupply} className="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700 font-medium h-10">
                                    Register Arrival
                                </button>
                            </div>
                        </div>
                    </div>
                )}

                <div className="bg-white p-6 rounded-xl shadow-sm border border-slate-200 mb-8">
                    <h2 className="text-lg font-semibold mb-4 text-slate-700 flex items-center gap-2">
                        <PieChart size={20} className="text-indigo-600"/> Warehouse Utilization
                    </h2>
                    <WarehouseChart />
                </div>

                {/* Control Panel */}
                <div className="bg-white p-6 rounded-xl shadow-sm border border-slate-200 mb-8">
                    <h2 className="text-lg font-semibold mb-4 text-slate-700 flex items-center gap-2">
                        <RefreshCw size={20} className="text-emerald-600"/> Operations Control
                    </h2>
                    <div className="flex flex-wrap gap-4 items-end bg-slate-50 p-4 rounded-lg border border-slate-100">
                        <div>
                            <label className="block text-xs font-bold text-slate-500 uppercase tracking-wide mb-1">Target Supply ID</label>
                            <input
                                type="number"
                                value={supplyId}
                                onChange={(e) => setSupplyId(e.target.value)}
                                className="border border-slate-300 rounded-md p-2 w-48 focus:ring-2 focus:ring-emerald-500 outline-none transition"
                                placeholder="Enter ID..."
                            />
                        </div>
                        <button
                            onClick={handleCalculate}
                            disabled={loading}
                            className={`px-6 py-2 rounded-md text-white font-medium flex items-center gap-2 transition-all shadow-sm
                ${loading ? 'bg-slate-400 cursor-not-allowed' : 'bg-emerald-600 hover:bg-emerald-700 active:transform active:scale-95'}
              `}
                        >
                            {loading ? <RefreshCw className="animate-spin" size={18}/> : <Play size={18} fill="currentColor"/>}
                            {loading ? 'Processing...' : 'Run Distribution Algorithm'}
                        </button>
                    </div>
                </div>

                {/* Results Table */}
                <div className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
                    <div className="px-6 py-4 border-b border-slate-100 flex justify-between items-center bg-slate-50/50">
                        <h2 className="text-lg font-semibold text-slate-700 flex items-center gap-2">
                            <Box size={20} className="text-blue-600"/> Distribution History
                        </h2>
                        <div className="flex items-center gap-2 text-xs text-slate-400">
                            <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
                            Live Updates
                        </div>
                    </div>

                    <div className="overflow-x-auto">
                        <table className="w-full text-left text-sm">
                            <thead className="bg-slate-50 text-slate-500 uppercase tracking-wider text-xs font-semibold border-b border-slate-200">
                            <tr>
                                <th className="px-6 py-3">ID</th>
                                <th className="px-6 py-3">Route (From → To)</th>
                                <th className="px-6 py-3">Status</th>
                                <th className="px-6 py-3">Items Content</th>
                                <th className="px-6 py-3">Timestamp</th>
                            </tr>
                            </thead>
                            <tbody className="divide-y divide-slate-100">
                            {shipments
                                .filter(s => s.sourceId !== s.destinationId)
                                .map((s) => (
                                    <tr key={s.id} className="hover:bg-slate-50 transition-colors">
                                        <td className="px-6 py-4 font-medium text-slate-700">#{s.id}</td>
                                        <td className="px-6 py-4">
                                            <div className="flex items-center gap-3">
                                                <span className="px-2 py-1 bg-white border border-slate-200 rounded text-slate-600 font-mono text-xs">WH-{s.sourceId}</span>
                                                <ArrowRight size={14} className="text-slate-400" />
                                                <span className="px-2 py-1 bg-blue-50 border border-blue-100 text-blue-700 rounded font-mono text-xs font-bold">WH-{s.destinationId}</span>
                                            </div>
                                        </td>
                                        <td className="px-6 py-4">
                      <span className={`px-2.5 py-1 rounded-full text-xs font-bold border
                        ${s.status === 'PLANNED' ? 'bg-emerald-50 text-emerald-700 border-emerald-100' : 'bg-slate-100 text-slate-600 border-slate-200'}
                      `}>
                        {s.status}
                      </span>
                                        </td>
                                        <td className="px-6 py-4 text-slate-600">
                                            {s.items.map((i, idx) => (
                                                <div key={idx} className="flex items-center gap-1">
                                                    <span className="font-bold text-slate-800">{i.quantity}x</span>
                                                    <span>{i.productName}</span>
                                                </div>
                                            ))}
                                        </td>
                                        <td className="px-6 py-4 text-slate-400 font-mono text-xs">
                                            {new Date(s.createdAt).toLocaleTimeString()}
                                        </td>
                                    </tr>
                                ))}
                            {shipments.filter(s => s.sourceId !== s.destinationId).length === 0 && (
                                <tr>
                                    <td colSpan="5" className="px-6 py-12 text-center text-slate-400 italic">
                                        No active shipments found (Items retained at source or no data).
                                    </td>
                                </tr>
                            )}
                            </tbody>
                        </table>
                    </div>
                </div>
            </main>
        </div>
    );
}