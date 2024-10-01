import React, { useState, useEffect } from 'react';
import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts';
import { Link } from 'react-router-dom';
import { PublicKeyDisplay } from './SharedComponents';

const Dashboard = ({ isDarkMode, network }) => {
    const [data, setData] = useState({
        activeOperators: 0,
        activeValidators: 0,
        activeClusters: 0,
        stakedETH: 0,
        proposedBlocks: 0,
        networkFee: "",
        operatorValidatorLimit: 0,
        liquidationThreshold: 0,
        minimumCollateral: "",
        events: [],
        blocks: [],
        charts: [],
        validators: []
    });
    const [error, setError] = useState(null);

    useEffect(() => {
        fetchData();
    }, []);

    const fetchData = async () => {
        setError(null);
        try {
            const response = await fetch('/api/dashboard');
            if (!response.ok) {
                throw new Error('Failed to fetch dashboard data');
            }

            const jsonData = await response.json();
            setData({
                activeOperators: jsonData.activeOperators,
                activeValidators: jsonData.activeValidators,
                activeClusters: jsonData.activeClusters,
                stakedETH: jsonData.stakedETH,
                proposedBlocks: jsonData.proposedBlocks,
                networkFee: jsonData.networkFee,
                operatorValidatorLimit: jsonData.operatorValidatorLimit,
                liquidationThreshold: jsonData.liquidationThreshold,
                minimumCollateral: jsonData.minimumCollateral,
                events: jsonData.events || [],
                blocks: jsonData.blocks || [],
                validators: jsonData.validators || [],
                charts: jsonData.charts || []
            });
        } catch (error) {
            setError('Failed to fetch dashboard data. Please try again later.');
        }
    };

    const getEtherscanUrl = (type, value) => {
        const baseUrl = network === 'mainnet'
            ? 'https://etherscan.io'
            : 'https://holesky.etherscan.io';
        return `${baseUrl}/${type}/${value}`;
    };

    const getBeaconscanUrl = (type, value) => {
        const baseUrl = network === 'mainnet'
            ? 'https://beaconcha.in'
            : 'https://holesky.beaconcha.in';
        return `${baseUrl}/${type}/${value}`;
    };

    const truncateHash = (hash) => {
        return `${hash.slice(0, 6)}...${hash.slice(-4)}`;
    };

    const calculateYAxisMax = (data) => {
        const maxValidators = Math.max(...data.map(item => item.validators));

        const baseMax = Math.ceil(maxValidators / 10000) * 10000;

        const buffer = maxValidators * 0.1;
        const adjustedMax = Math.ceil((maxValidators + buffer) / 1000) * 1000;

        return Math.min(baseMax, adjustedMax);
    };

    if (error) return <div className={`p-8 ${isDarkMode ? 'bg-gray-900 text-gray-200' : 'bg-white text-gray-800'}`}>{error}</div>;

    return (
        <div className={`${isDarkMode ? 'text-white' : 'text-black'}`}>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-4">
                <div className={`${isDarkMode ? 'bg-gray-800' : 'bg-gray-200'} p-4 rounded`}>
                    <h2 className={`text-sm ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>Active Operators</h2>
                    <p className="text-2xl font-bold">{(data.activeOperators || 0).toLocaleString()}</p>
                </div>
                <div className={`${isDarkMode ? 'bg-gray-800' : 'bg-gray-200'} p-4 rounded`}>
                    <h2 className={`text-sm ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>Active Validators</h2>
                    <p className="text-2xl font-bold">{(data.activeValidators || 0).toLocaleString()}</p>
                </div>
                <div className={`${isDarkMode ? 'bg-gray-800' : 'bg-gray-200'} p-4 rounded`}>
                    <h2 className={`text-sm ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>Staked ETH</h2>
                    <p className="text-2xl font-bold">{(data.stakedETH || 0).toLocaleString()}</p>
                </div>
                <div className={`${isDarkMode ? 'bg-gray-800' : 'bg-gray-200'} p-4 rounded`}>
                    {network === 'mainnet' ? (
                        <>
                            <h2 className={`text-sm ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>Proposed Blocks</h2>
                            <p className="text-2xl font-bold">{(data.proposedBlocks || 0).toLocaleString()}</p>
                        </>
                    ) : (
                        <>
                            <h2 className={`text-sm ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>Active Clusters</h2>
                            <p className="text-2xl font-bold">{(data.activeClusters || 0).toLocaleString()}</p>
                        </>
                    )}
                </div>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-4">
                <div className={`${isDarkMode ? 'bg-gray-800' : 'bg-gray-200'} p-4 rounded`}>
                    <h2 className={`text-sm ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>Operator Validator Limit</h2>
                    <p className="text-2xl font-bold">{(data.operatorValidatorLimit || 0).toLocaleString()}</p>
                </div>
                <div className={`${isDarkMode ? 'bg-gray-800' : 'bg-gray-200'} p-4 rounded`}>
                    <h2 className={`text-sm ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>Liquidation Threshold</h2>
                    <p className="text-2xl font-bold">
                        {(data.liquidationThreshold || 0).toLocaleString()}
                        <span className="text-sm font-normal ml-1">blocks</span>
                    </p>
                </div>
                <div className={`${isDarkMode ? 'bg-gray-800' : 'bg-gray-200'} p-4 rounded`}>
                    <h2 className={`text-sm ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>Minimum Collateral</h2>
                    <p className="text-2xl font-bold">
                        {data.minimumCollateral}
                        <span className="text-sm font-normal ml-1">ssv</span>
                    </p>
                </div>
                <div className={`${isDarkMode ? 'bg-gray-800' : 'bg-gray-200'} p-4 rounded`}>
                    <h2 className={`text-sm ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>Network Fee</h2>
                    <p className="text-2xl font-bold">
                        {data.networkFee}
                        <span className="text-sm font-normal ml-1">ssv</span>
                    </p>
                </div>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-4">
                <div className={`${isDarkMode ? 'bg-gray-800' : 'bg-gray-200'} p-4 rounded`}>
                    <h2 className="text-xl mb-4">Validators</h2>
                    <ResponsiveContainer width="100%" height={200}>
                        <LineChart data={data.charts}>
                            <XAxis dataKey="name" stroke={isDarkMode ? "#888" : "#333"} />
                            <YAxis
                                stroke={isDarkMode ? "#888" : "#333"}
                                domain={[0, calculateYAxisMax(data.charts)]}
                            />
                            <Tooltip
                                contentStyle={{
                                    backgroundColor: isDarkMode ? '#1f2937' : '#f3f4f6',
                                    border: 'none',
                                    borderRadius: '8px',
                                    color: isDarkMode ? '#fff' : '#000'
                                }}
                            />
                            <Line type="monotone" dataKey="validators" stroke="#8884d8" />
                        </LineChart>
                    </ResponsiveContainer>
                </div>
                <div className={`${isDarkMode ? 'bg-gray-800' : 'bg-gray-200'} p-4 rounded`}>
                    <h2 className="text-xl mb-4">Operators</h2>
                    <ResponsiveContainer width="100%" height={200}>
                        <LineChart data={data.charts}>
                            <XAxis dataKey="name" stroke={isDarkMode ? "#888" : "#333"} />
                            <YAxis stroke={isDarkMode ? "#888" : "#333"} />
                            <Tooltip
                                contentStyle={{
                                    backgroundColor: isDarkMode ? '#1f2937' : '#f3f4f6',
                                    border: 'none',
                                    borderRadius: '8px',
                                    color: isDarkMode ? '#fff' : '#000'
                                }}
                            />
                            <Line type="monotone" dataKey="operators" stroke="#82ca9d" />
                        </LineChart>
                    </ResponsiveContainer>
                </div>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
                <div className={`${isDarkMode ? 'bg-gray-800' : 'bg-gray-200'} p-4 rounded`}>
                    <h2 className="text-xl mb-4">Latest events</h2>
                    {data.events.length > 0 ? (
                        <div className="overflow-x-auto">
                            <table className="w-full">
                                <thead>
                                    <tr className="text-left">
                                        <th className="p-2">Block</th>
                                        <th className="p-2">Tx Hash</th>
                                        <th className="p-2">Action</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {data.events.map((event, index) => (
                                        <tr key={index} className={`${index % 2 === 0 ? 'bg-gray-100 dark:bg-gray-700' : ''}`}>
                                            <td className="p-2">
                                                <a
                                                    href={getEtherscanUrl('block', event.block)}
                                                    target="_blank"
                                                    rel="noopener noreferrer"
                                                    className="text-blue-500 hover:underline"
                                                >
                                                    {event.block}
                                                </a>
                                            </td>
                                            <td className="p-2">
                                                <a
                                                    href={getEtherscanUrl('tx', event.transactionHash)}
                                                    target="_blank"
                                                    rel="noopener noreferrer"
                                                    className="text-blue-500 hover:underline"
                                                >
                                                    {truncateHash(event.transactionHash)}
                                                </a>
                                            </td>
                                            <td className="p-2">{event.action}</td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    ) : (
                        <p className={`${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>No latest events</p>
                    )}
                </div>
                <div className={`${isDarkMode ? 'bg-gray-800' : 'bg-gray-200'} p-4 rounded`}>
                    {network === 'mainnet' ? (
                        <>
                            <h2 className="text-xl mb-4">Latest blocks</h2>
                            {data.blocks.length && data.blocks.length > 0 ? (
                                <div className="overflow-x-auto">
                                    <table className="w-full">
                                        <thead>
                                            <tr className="text-left">
                                                <th className="p-2">Epoch</th>
                                                <th className="p-2">Slot</th>
                                                <th className="p-2">Block</th>
                                                <th className="p-2">Proposer</th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {data.blocks.map((block, index) => (
                                                <tr key={index} className={`${index % 2 === 0 ? 'bg-gray-100 dark:bg-gray-700' : ''}`}>
                                                    <td className="p-2">{block.epoch}</td>
                                                    <td className="p-2">
                                                        <a
                                                            href={getBeaconscanUrl('slot', block.slot)}
                                                            target="_blank"
                                                            rel="noopener noreferrer"
                                                            className="text-blue-500 hover:underline"
                                                        >
                                                            {block.slot}
                                                        </a>
                                                    </td>
                                                    <td className="p-2">
                                                        <a
                                                            href={getEtherscanUrl('block', block.blockNumber)}
                                                            target="_blank"
                                                            rel="noopener noreferrer"
                                                            className="text-blue-500 hover:underline"
                                                        >
                                                            {block.blockNumber}
                                                        </a>
                                                    </td>
                                                    <td className="p-2">{block.proposer}</td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                </div>
                            ) : (
                                <p className={`${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>No latest blocks</p>
                            )}
                        </>
                    ) : (
                        <>
                            <h2 className="text-xl mb-4">Latest validators</h2>
                            {data.validators && data.validators.length > 0 ? (
                                <div className="overflow-x-auto">
                                    <table className="w-full">
                                        <thead>
                                            <tr className="text-left">
                                                <th className="p-2">Public Key</th>
                                                <th className="p-2">Owner</th>
                                                <th className="p-2">Cluster ID</th>
                                            </tr>
                                        </thead>
                                        <tbody>
                                            {data.validators.map((validator, index) => (<tr key={index} className={`${index % 2 === 0 ? 'bg-gray-100 dark:bg-gray-700' : ''}`}>
                                                <td className={`p-2 ${isDarkMode ? 'text-gray-100' : 'text-gray-700'}`}>
                                                    <PublicKeyDisplay
                                                        publicKey={validator.publicKey}
                                                        beaconchainLink={getBeaconscanUrl('validator', `0x${validator.publicKey}`)}
                                                        isDarkMode={isDarkMode}
                                                        isTruncate={true}
                                                    />
                                                </td>


                                                <td className="p-2">
                                                    <Link
                                                        to={`/account/${validator.owner}`}
                                                        className={`hover:underline ${isDarkMode ? 'text-blue-400 hover:text-blue-300' : 'text-blue-600 hover:text-blue-800'}`}
                                                    >
                                                        {truncateHash(validator.owner)}
                                                    </Link>
                                                </td>
                                                <td className="p-2">
                                                    <Link
                                                        to={`/cluster/${validator.clusterId}`}
                                                        className={`hover:underline ${isDarkMode ? 'text-blue-400 hover:text-blue-300' : 'text-blue-600 hover:text-blue-800'
                                                            }`}
                                                    >
                                                        {truncateHash(validator.clusterId)}
                                                    </Link>
                                                </td>
                                            </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                </div>
                            ) : (
                                <p className={`${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>No validators data</p>
                            )}
                        </>
                    )}
                </div>
            </div>
        </div >
    );
};

export default Dashboard;