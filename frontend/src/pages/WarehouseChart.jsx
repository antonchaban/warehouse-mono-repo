import React, { useEffect, useState } from 'react';
import { PieChart, Pie, Cell, Tooltip, Legend, ResponsiveContainer } from 'recharts';

const COLORS = ['#00C49F', '#FF8042']; // Зелений (вільно), Помаранчевий (зайнято)

const WarehouseChart = () => {
    const [data, setData] = useState([]);

    useEffect(() => {
        fetch('http://localhost:8080/api/admin/warehouses/stats', {
            headers: {
                'Authorization': 'Bearer ' + localStorage.getItem('token')
            }
        })
            .then(res => res.json())
            .then(data => setData(data))
            .catch(err => console.error(err));
    }, []);

    if (!data.length) return <div>Loading charts...</div>;

    return (
        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '20px', marginTop: '20px' }}>
            {data.map((wh, idx) => (
                <div key={idx} style={cardStyle}>
                    <h3>{wh.name}</h3>
                    <div style={{ width: '100%', height: 200 }}>
                        <ResponsiveContainer>
                            <PieChart>
                                <Pie
                                    data={[
                                        { name: 'Free', value: wh.freeCapacity },
                                        { name: 'Used', value: wh.usedCapacity }
                                    ]}
                                    innerRadius={60}
                                    outerRadius={80}
                                    paddingAngle={5}
                                    dataKey="value"
                                >
                                    <Cell fill={COLORS[0]} />
                                    <Cell fill={COLORS[1]} />
                                </Pie>
                                <Tooltip />
                                <Legend />
                            </PieChart>
                        </ResponsiveContainer>
                    </div>
                    <p><strong>{wh.utilizationPercentage.toFixed(1)}%</strong> Full</p>
                    <p style={{fontSize: '0.8em', color: '#666'}}>
                        {wh.usedCapacity.toFixed(1)} / {wh.totalCapacity} m³
                    </p>
                </div>
            ))}
        </div>
    );
};

const cardStyle = {
    border: '1px solid #e0e0e0',
    borderRadius: '10px',
    padding: '15px',
    width: '250px',
    backgroundColor: 'white',
    boxShadow: '0 4px 6px rgba(0,0,0,0.1)',
    textAlign: 'center'
};

export default WarehouseChart;