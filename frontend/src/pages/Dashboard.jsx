import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
    LogOut, Play, RefreshCw, Box, ArrowRight, Truck, Plus,
    Warehouse, Package, PieChart, Database
} from 'lucide-react';
import api from '../api';
import WarehouseChart from './WarehouseChart';

export default function Dashboard() {
    // Ініціалізуємо як порожні масиви []
    const [shipments, setShipments] = useState([]);
    const [supplyId, setSupplyId] = useState('555');
    const [loading, setLoading] = useState(false);

    const [warehousesList, setWarehousesList] = useState([]);
    const [productsList, setProductsList] = useState([]);
    const [suppliesList, setSuppliesList] = useState([]);

    const navigate = useNavigate();

    const [showForms, setShowForms] = useState(false);
    const [newSupplyWhId, setNewSupplyWhId] = useState('');
    const [newSupplyProdId, setNewSupplyProdId] = useState('');
    const [newSupplyQty, setNewSupplyQty] = useState('');
    const [newWhCapacity, setNewWhCapacity] = useState('');
    const [newProdVolume, setNewProdVolume] = useState('');

    // --- БЕЗПЕЧНЕ ЗАВАНТАЖЕННЯ ДАНИХ ---
    const fetchAllData = async () => {
        try {
            // 1. Shipments
            const resShipments = await api.get('/v1/distribution/shipments');
            // Перевіряємо, чи це масив. Якщо ні — ставимо пустий масив.
            setShipments(Array.isArray(resShipments.data) ? resShipments.data : []);

            // 2. Warehouses
            const resWh = await api.get('/admin/warehouses');
            setWarehousesList(Array.isArray(resWh.data) ? resWh.data : []);

            // 3. Products
            const resProd = await api.get('/admin/products');
            setProductsList(Array.isArray(resProd.data) ? resProd.data : []);

            // 4. Supplies
            const resSup = await api.get('/admin/supplies');
            setSuppliesList(Array.isArray(resSup.data) ? resSup.data : []);

        } catch (err) {
            console.error("Failed to fetch data", err);
            // Не кидаємо alert кожні 3 секунди, просто пишемо в консоль
        }
    };

    useEffect(() => {
        fetchAllData();
        const interval = setInterval(fetchAllData, 3000);
        return () => clearInterval(interval);
    }, []);
    // ---------------------------

    const createSupply = async () => {
        try {
            const res = await api.post('/admin/supplies', {
                warehouseId: parseInt(newSupplyWhId),
                productId: parseInt(newSupplyProdId),
                quantity: parseInt(newSupplyQty)
            });

            let newId = 'unknown';
            // Безпечна перевірка відповіді
            if (res.data && typeof res.data === 'string' && res.data.includes('ID:')) {
                newId = res.data.split('ID:')[1].trim();
            }

            alert(`Supply Created! ID: ${newId}`);
            if (newId !== 'unknown') setSupplyId(newId);

            setNewSupplyWhId(''); setNewSupplyProdId(''); setNewSupplyQty('');
            fetchAllData();
        } catch (e) {
            console.error(e);
            alert('Error creating supply. Check console/network.');
        }
    };

    const createWarehouse = async () => {
        try {
            await api.post('/admin/warehouses', { capacity: parseFloat(newWhCapacity) });
            alert('Warehouse Created!');
            setNewWhCapacity('');
            fetchAllData();
        } catch (e) { alert('Error creating warehouse'); }
    };

    const createProduct = async () => {
        try {
            await api.post('/admin/products', { volume: parseFloat(newProdVolume) });
            alert('Product Created!');
            setNewProdVolume('');
            fetchAllData();
        } catch (e) { alert('Error creating product'); }
    };

    const handleCalculate = async () => {
        setLoading(true);
        try {
            await api.post('/v1/distribution/calculate', { supplyId: parseInt(supplyId) });
        } catch (err) {
            alert('Error triggering calculation.');
        } finally {
            setTimeout(() => setLoading(false), 1000);
        }
    };

    const handleLogout = () => {
        localStorage.removeItem('token');
        navigate('/');
    };

    return (
        <div className="min-h-screen bg-slate-50 text-slate-800 font-sans">
            <header className="bg-white shadow-sm border-b border-slate-200 sticky top-0 z-20">
                <div className="max-w-7xl mx-auto px-4 py-3 flex justify-between items-center">
                    <h1 className="text-lg font-bold text-indigo-700 flex items-center gap-2">
                        <Truck className="text-indigo-600" /> Intelligent Distribution System
                    </h1>
                    <button onClick={handleLogout} className="text-slate-500 hover:text-red-600 flex items-center gap-1 text-sm font-medium transition">
                        <LogOut size={16} /> Logout
                    </button>
                </div>
            </header>

            <main className="max-w-7xl mx-auto px-4 py-6 space-y-8">

                {/* 1. Блок Графіків */}
                <section className="bg-white p-5 rounded-xl shadow-sm border border-slate-200">
                    <h2 className="text-md font-bold text-slate-700 mb-4 flex items-center gap-2">
                        <PieChart size={18} className="text-indigo-500"/> Real-time Warehouse Utilization
                    </h2>
                    <WarehouseChart />
                </section>

                {/* 2. ДОВІДКОВІ ДАНІ */}
                <section>
                    <h2 className="text-md font-bold text-slate-700 mb-4 flex items-center gap-2">
                        <Database size={18} className="text-blue-500"/> Infrastructure & Inventory Overview
                    </h2>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-6">

                        {/* Таблиця Складів */}
                        <div className="bg-white rounded-lg shadow-sm border border-slate-200 overflow-hidden">
                            <div className="bg-slate-50 px-4 py-2 border-b border-slate-200 font-semibold text-xs text-slate-500 uppercase flex justify-between">
                                <span>Warehouses</span> <Warehouse size={14}/>
                            </div>
                            <div className="max-h-48 overflow-y-auto">
                                <table className="w-full text-xs text-left">
                                    <thead className="bg-slate-50 text-slate-400 sticky top-0">
                                    <tr><th className="p-2">ID</th><th className="p-2">Capacity</th></tr>
                                    </thead>
                                    <tbody className="divide-y divide-slate-100">
                                    {/* 👇 ЗАХИСТ: Перевіряємо чи це масив перед map */}
                                    {Array.isArray(warehousesList) && warehousesList.map(w => (
                                        <tr key={w.id} className="hover:bg-slate-50">
                                            <td className="p-2 font-mono font-bold text-indigo-600">WH-{w.id}</td>
                                            <td className="p-2">{w.totalCapacity} m³</td>
                                        </tr>
                                    ))}
                                    </tbody>
                                </table>
                            </div>
                        </div>

                        {/* Таблиця Продуктів */}
                        <div className="bg-white rounded-lg shadow-sm border border-slate-200 overflow-hidden">
                            <div className="bg-slate-50 px-4 py-2 border-b border-slate-200 font-semibold text-xs text-slate-500 uppercase flex justify-between">
                                <span>Products Catalog</span> <Package size={14}/>
                            </div>
                            <div className="max-h-48 overflow-y-auto">
                                <table className="w-full text-xs text-left">
                                    <thead className="bg-slate-50 text-slate-400 sticky top-0">
                                    <tr><th className="p-2">ID</th><th className="p-2">Volume / Item</th></tr>
                                    </thead>
                                    <tbody className="divide-y divide-slate-100">
                                    {/* 👇 ЗАХИСТ ТУТ (ймовірна причина помилки p.map) */}
                                    {Array.isArray(productsList) && productsList.map(p => (
                                        <tr key={p.id} className="hover:bg-slate-50">
                                            <td className="p-2 font-mono font-bold text-orange-600">PR-{p.id}</td>
                                            <td className="p-2">{p.volumeM3} m³</td>
                                        </tr>
                                    ))}
                                    </tbody>
                                </table>
                            </div>
                        </div>

                        {/* Таблиця Поставок */}
                        <div className="bg-white rounded-lg shadow-sm border border-slate-200 overflow-hidden">
                            <div className="bg-slate-50 px-4 py-2 border-b border-slate-200 font-semibold text-xs text-slate-500 uppercase flex justify-between">
                                <span>Incoming Supplies</span> <Truck size={14}/>
                            </div>
                            <div className="max-h-48 overflow-y-auto">
                                <table className="w-full text-xs text-left">
                                    <thead className="bg-slate-50 text-slate-400 sticky top-0">
                                    <tr><th className="p-2">ID</th><th className="p-2">To</th><th className="p-2">Status</th></tr>
                                    </thead>
                                    <tbody className="divide-y divide-slate-100">
                                    {/* 👇 ЗАХИСТ */}
                                    {Array.isArray(suppliesList) && suppliesList.map(s => (
                                        <tr key={s.id} className="hover:bg-slate-50">
                                            <td className="p-2 font-mono font-bold text-blue-600">SUP-{s.id}</td>
                                            <td className="p-2">WH-{s.warehouseId}</td>
                                            <td className="p-2">
                                                    <span className={`px-1.5 py-0.5 rounded text-[10px] font-bold border ${
                                                        s.status === 'RECEIVED' ? 'bg-yellow-50 text-yellow-700 border-yellow-200' :
                                                            s.status === 'PROCESSED' ? 'bg-green-50 text-green-700 border-green-200' :
                                                                'bg-slate-100 text-slate-600 border-slate-200'
                                                    }`}>
                                                        {s.status}
                                                    </span>
                                            </td>
                                        </tr>
                                    ))}
                                    </tbody>
                                </table>
                            </div>
                        </div>

                    </div>
                </section>

                {/* 3. Кнопки створення */}
                <div className="flex justify-end">
                    <button onClick={() => setShowForms(!showForms)} className="text-xs font-bold text-indigo-600 hover:underline">
                        {showForms ? "Hide Admin Forms" : "+ Create New Resources"}
                    </button>
                </div>

                {showForms && (
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6 p-4 bg-indigo-50/50 border border-indigo-100 rounded-xl">
                        <div className="bg-white p-4 rounded shadow-sm">
                            <h3 className="text-sm font-bold mb-2 flex gap-2"><Warehouse size={16}/> New Warehouse</h3>
                            <div className="flex gap-2">
                                <input type="number" placeholder="Cap" value={newWhCapacity} onChange={e=>setNewWhCapacity(e.target.value)} className="border p-1 rounded w-full text-sm"/>
                                <button onClick={createWarehouse} className="bg-indigo-600 text-white px-3 py-1 rounded text-sm"><Plus size={16}/></button>
                            </div>
                        </div>
                        <div className="bg-white p-4 rounded shadow-sm">
                            <h3 className="text-sm font-bold mb-2 flex gap-2"><Package size={16}/> New Product</h3>
                            <div className="flex gap-2">
                                <input type="number" placeholder="Vol" value={newProdVolume} onChange={e=>setNewProdVolume(e.target.value)} className="border p-1 rounded w-full text-sm"/>
                                <button onClick={createProduct} className="bg-orange-500 text-white px-3 py-1 rounded text-sm"><Plus size={16}/></button>
                            </div>
                        </div>
                        <div className="bg-white p-4 rounded shadow-sm md:col-span-2">
                            <h3 className="text-sm font-bold mb-2 flex gap-2"><Truck size={16}/> Register Incoming Supply</h3>
                            <div className="flex gap-2 items-end">
                                <input type="number" placeholder="WhID" value={newSupplyWhId} onChange={e=>setNewSupplyWhId(e.target.value)} className="border p-1 rounded w-20 text-sm"/>
                                <input type="number" placeholder="PrID" value={newSupplyProdId} onChange={e=>setNewSupplyProdId(e.target.value)} className="border p-1 rounded w-20 text-sm"/>
                                <input type="number" placeholder="Qty" value={newSupplyQty} onChange={e=>setNewSupplyQty(e.target.value)} className="border p-1 rounded w-20 text-sm"/>
                                <button onClick={createSupply} className="bg-blue-600 text-white px-4 py-1 rounded text-sm font-bold">Register</button>
                            </div>
                        </div>
                    </div>
                )}

                {/* 4. Панель управління */}
                <section className="bg-white p-6 rounded-xl shadow-sm border border-slate-200 flex flex-col md:flex-row justify-between items-center gap-4">
                    <div>
                        <h2 className="text-md font-bold text-slate-700 flex items-center gap-2">
                            <RefreshCw size={18} className="text-emerald-600"/> Operations Control
                        </h2>
                        <p className="text-xs text-slate-400 mt-1">Select a supply ID from the table above to distribute.</p>
                    </div>
                    <div className="flex gap-2 bg-slate-50 p-2 rounded-lg border border-slate-100">
                        <input
                            type="number"
                            value={supplyId}
                            onChange={(e) => setSupplyId(e.target.value)}
                            className="border border-slate-300 rounded px-3 py-1.5 w-32 focus:outline-none focus:border-emerald-500 text-sm font-mono"
                            placeholder="Supply ID"
                        />
                        <button
                            onClick={handleCalculate}
                            disabled={loading}
                            className={`px-4 py-1.5 rounded text-white text-sm font-bold flex items-center gap-2 shadow-sm
                                ${loading ? 'bg-slate-400' : 'bg-emerald-600 hover:bg-emerald-700'}
                            `}
                        >
                            {loading ? <RefreshCw className="animate-spin" size={16}/> : <Play size={16} fill="currentColor"/>}
                            {loading ? 'Thinking...' : 'Run Algo'}
                        </button>
                    </div>
                </section>

                {/* 5. Таблиця Історії */}
                <section className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
                    <div className="px-6 py-3 border-b border-slate-100 bg-slate-50/50 flex justify-between">
                        <h2 className="text-sm font-bold text-slate-700 flex items-center gap-2">
                            <Box size={16} className="text-indigo-600"/> Distribution Routes (Results)
                        </h2>
                        <span className="text-xs text-emerald-500 font-bold animate-pulse">● Live</span>
                    </div>
                    <div className="overflow-x-auto">
                        <table className="w-full text-left text-sm">
                            <thead className="bg-slate-50 text-slate-500 uppercase text-[10px] font-bold border-b border-slate-200">
                            <tr>
                                <th className="px-6 py-2">ID</th>
                                <th className="px-6 py-2">Route</th>
                                <th className="px-6 py-2">Status</th>
                                <th className="px-6 py-2">Cargo</th>
                                <th className="px-6 py-2">Time</th>
                            </tr>
                            </thead>
                            <tbody className="divide-y divide-slate-100">
                            {/* 👇 ЗАХИСТ І ТУТ */}
                            {Array.isArray(shipments) && shipments
                                .filter(s => s.sourceId !== s.destinationId)
                                .map((s) => (
                                    <tr key={s.id} className="hover:bg-slate-50">
                                        <td className="px-6 py-3 font-mono text-xs text-slate-500">#{s.id}</td>
                                        <td className="px-6 py-3">
                                            <div className="flex items-center gap-2">
                                                <span className="px-2 py-0.5 bg-white border rounded text-xs font-mono">WH-{s.sourceId}</span>
                                                <ArrowRight size={12} className="text-slate-400" />
                                                <span className="px-2 py-0.5 bg-indigo-50 border border-indigo-100 text-indigo-700 rounded text-xs font-mono font-bold">WH-{s.destinationId}</span>
                                            </div>
                                        </td>
                                        <td className="px-6 py-3">
                                            <span className="px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-50 text-emerald-600 border border-emerald-100">
                                                {s.status}
                                            </span>
                                        </td>
                                        <td className="px-6 py-3">
                                            {s.items.map((i, idx) => (
                                                <div key={idx} className="text-xs font-bold text-slate-700">
                                                    {i.quantity} x {i.productName}
                                                </div>
                                            ))}
                                        </td>
                                        <td className="px-6 py-3 text-xs text-slate-400 font-mono">
                                            {new Date(s.createdAt).toLocaleTimeString()}
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </section>
            </main>
        </div>
    );
}