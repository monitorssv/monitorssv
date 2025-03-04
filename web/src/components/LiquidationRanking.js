import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { OperatorDisplay } from './SharedComponents';
import { AlertCircle } from 'lucide-react';

const LiquidationRanking = ({ isDarkMode, network }) => {
    const [clusters, setClusters] = useState([]);
    const [simulatedClusters, setSimulatedClusters] = useState([]);
    const [currentPageState, setCurrentPageState] = useState({
        current: 1,
        simulated: 1
    });
    const [expandedRowState, setExpandedRowState] = useState({
        current: null,
        simulated: null
    });
    const [itemsPerPage, setItemsPerPage] = useState(10);
    const [totalPages, setTotalPages] = useState({ current: 1, simulated: 1 });
    const [totalItems, setTotalItems] = useState({ current: 0, simulated: 0 });
    const [error, setError] = useState(null);
    const [activeTab, setActiveTab] = useState('current');
    const [networkFees, setNetworkFees] = useState({
        current: 0,
        upcoming: 0
    });
    const [isLoading, setIsLoading] = useState(true);

    const currentPage = currentPageState[activeTab];

    useEffect(() => {
        const initializeData = async () => {
            setIsLoading(true);
            try {
                await Promise.all([
                    fetchNetworkFees(),
                    fetchClustersInfo(),
                    fetchSimulatedClustersInfo()
                ]);
            } catch (err) {
                setError('Failed to fetch data. Please try again later.');
            } finally {
                setIsLoading(false);
            }
        };

        initializeData();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    useEffect(() => {
        const updateData = async () => {
            setIsLoading(true);
            try {
                if (activeTab === 'current') {
                    await fetchClustersInfo();
                } else {
                    await fetchSimulatedClustersInfo();
                }
            } catch (err) {
                setError('Failed to fetch data. Please try again later.');
            } finally {
                setIsLoading(false);
            }
        };

        updateData();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [currentPageState, itemsPerPage]);

    const fetchClustersInfo = async () => {
        setError(null);
        try {
            const response = await fetch(`/api/get30DayLiquidationRankingClusters?page=${currentPageState.current}&limit=${itemsPerPage}`);
            if (!response.ok) throw new Error('Failed to fetch clusters');
            const data = await response.json();
            setClusters(data.clusters);
            setTotalPages(prev => ({ ...prev, current: data.totalPages }));
            setTotalItems(prev => ({ ...prev, current: data.totalItems }));
        } catch (err) {
            throw new Error('Failed to fetch clusters');
        }
    };

    const fetchSimulatedClustersInfo = async () => {
        setError(null);
        try {
            const response = await fetch(`/api/get30DaySimulatedLiquidationRankingClusters?page=${currentPageState.simulated}&limit=${itemsPerPage}`);
            if (!response.ok) throw new Error('Failed to fetch simulated clusters');
            const data = await response.json();
            setSimulatedClusters(data.clusters);
            setTotalPages(prev => ({ ...prev, simulated: data.totalPages }));
            setTotalItems(prev => ({ ...prev, simulated: data.totalItems }));
        } catch (err) {
            throw new Error('Failed to fetch simulated clusters');
        }
    };

    const fetchNetworkFees = async () => {
        try {
            const response = await fetch('/api/getNetworkFees');
            if (!response.ok) throw new Error('Failed to fetch network fees');
            const data = await response.json();
            setNetworkFees(data);
        } catch (err) {
            setError('Failed to fetch networkFees');
        }
    };

    const toggleRowExpansion = (id) => {
        setExpandedRowState(prev => ({
            ...prev,
            [activeTab]: prev[activeTab] === id ? null : id
        }));
    };

    const truncateAddress = (address) => {
        if (!address) return '';
        return `${address.slice(0, 6)}...${address.slice(-4)}`;
    };

    const formatRunaway = (blocks) => {
        const days = Math.floor(blocks / 7200);
        const hours = Math.floor((blocks - days * 7200) * 12 / 3600);
        return `${days}d ${hours}h`;
    };

    const paginate = (pageNumber) => {
        setCurrentPageState(prev => ({
            ...prev,
            [activeTab]: pageNumber
        }));
    };

    const renderTabButtons = () => {
        const isSimulatedDisabled = networkFees.current === networkFees.upcoming;

        return (
            <div className="mb-6">
                <div className="flex justify-center">
                    <div className="w-full">
                        <div className="flex gap-4">
                            <button
                                onClick={() => setActiveTab('current')}
                                className={`flex-1 px-5 py-3 rounded-lg transition-all duration-200
                                    ${activeTab === 'current'
                                    ? isDarkMode
                                        ? 'bg-blue-600 text-white shadow-lg shadow-blue-500/20'
                                        : 'bg-blue-500 text-white shadow-lg shadow-blue-500/20'
                                    : isDarkMode
                                        ? 'bg-gray-800 text-gray-400 hover:bg-gray-700'
                                        : 'bg-white text-gray-600 hover:bg-gray-50'
                                }`}
                            >
                                <div className="flex items-center gap-8">
                                    <div className="flex-shrink-0">
                                        <h3 className="text-lg font-semibold">Current Ranking</h3>
                                    </div>
                                    <div className={`flex-1 text-right ${activeTab === 'current' ? 'text-blue-100' : 'text-gray-500'}`}>
                                        <span className="text-sm whitespace-nowrap">
                                            Based on current network fee: {networkFees.current} ssv
                                        </span>
                                    </div>
                                </div>
                            </button>
                            <button
                                onClick={() => !isSimulatedDisabled && setActiveTab('simulated')}
                                disabled={isSimulatedDisabled}
                                className={`flex-1 px-5 py-3 rounded-lg transition-all duration-200
                                    ${isSimulatedDisabled
                                    ? isDarkMode
                                        ? 'bg-gray-700 cursor-not-allowed opacity-60'
                                        : 'bg-gray-100 cursor-not-allowed opacity-60'
                                    : activeTab === 'simulated'
                                        ? isDarkMode
                                            ? 'bg-blue-600 text-white shadow-lg shadow-blue-500/20'
                                            : 'bg-blue-500 text-white shadow-lg shadow-blue-500/20'
                                        : isDarkMode
                                            ? 'bg-gray-800 text-gray-400 hover:bg-gray-700'
                                            : 'bg-white text-gray-600 hover:bg-gray-50'
                                }`}
                            >
                                <div className="flex items-center gap-8">
                                    <div className="flex-shrink-0">
                                        <h3 className="text-lg font-semibold">Simulated Ranking</h3>
                                    </div>
                                    <div className={`flex-1 text-right ${isSimulatedDisabled
                                        ? 'text-gray-500'
                                        : activeTab === 'simulated'
                                            ? 'text-blue-100'
                                            : 'text-gray-500'
                                    }`}>
                                        <span className="text-sm whitespace-nowrap">
                                            {isSimulatedDisabled
                                                ? 'No network fee changes to simulate'
                                                : `Preview of upcoming network fee: ${networkFees.upcoming} ssv`
                                            }
                                        </span>
                                    </div>
                                </div>
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        );
    };

    const renderContent = () => {
        const currentData = activeTab === 'current' ? clusters : simulatedClusters;
        const expandedRow = expandedRowState[activeTab];
        const currentTotalItems = totalItems[activeTab];

        if (isLoading) {
            return (
                <div className={`flex justify-center items-center py-8 ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>
                    <div className="animate-pulse flex flex-col items-center">
                        <div className="h-8 w-8 border-4 border-t-blue-500 border-r-transparent border-b-blue-500 border-l-transparent rounded-full animate-spin mb-4"></div>
                        <span>Loading data...</span>
                    </div>
                </div>
            );
        }

        return (
            <div>
                <div className={`mb-6 px-5 py-3 rounded-lg ${isDarkMode ? 'bg-gray-800/50' : 'bg-gray-50'} border ${isDarkMode ? 'border-gray-700' : 'border-gray-200'}`}>
                    <div className="flex items-center gap-2">
                        <AlertCircle className={`w-4 h-4 ${activeTab === 'current' ? 'text-blue-500' : 'text-yellow-500'}`} />
                        <p className={`text-sm font-medium ${isDarkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                            {activeTab === 'current'
                                ? 'These clusters are at risk of liquidation based on current cluster fee!'
                                : 'If the network fee is upgraded today with real-time parameters, these clusters risk liquidation because of the new cluster fee!'
                            }
                        </p>
                    </div>
                </div>
                <div className="mb-4 text-sm text-gray-500">
                    Showing {currentData.length} of {currentTotalItems} clusters with operational runway less than 30 days
                </div>
                <div className={`overflow-x-auto rounded-lg shadow ${isDarkMode ? 'bg-gray-800' : 'bg-white'}`}>
                    <table className="w-full">
                        <thead>
                        <tr className={isDarkMode ? 'bg-gray-700' : 'bg-gray-200'}>
                            <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Rank</th>
                            <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Cluster ID</th>
                            <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Validators</th>
                            <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Runway</th>
                        </tr>
                        </thead>
                        <tbody>
                        {currentData.map((cluster, index) => (
                            <React.Fragment key={cluster.id}>
                                <tr
                                    className={`border-b ${isDarkMode ? 'border-gray-700' : 'border-gray-200'} hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-150 cursor-pointer`}
                                    onClick={() => toggleRowExpansion(cluster.id)}
                                >
                                    <td className="p-3">{(currentPage - 1) * itemsPerPage + index + 1}</td>
                                    <td className="p-3">
                                        <Link
                                            to={`/cluster/${cluster.id}`}
                                            className={`hover:underline ${isDarkMode ? 'text-blue-400 hover:text-blue-300' : 'text-blue-600 hover:text-blue-800'
                                            }`}
                                        >
                                            {truncateAddress(cluster.id)}
                                        </Link>
                                    </td>
                                    <td className="p-3">{cluster.validatorCount}</td>
                                    <td className="p-3">
                                            <span className="px-2 py-1 rounded-full bg-red-500 text-white text-sm">
                                                {formatRunaway(cluster.operationalRunaway)}
                                            </span>
                                    </td>
                                </tr>
                                {expandedRow === cluster.id && (
                                    <tr className={`border-b ${isDarkMode ? 'border-gray-700 bg-gray-750' : 'border-gray-200 bg-gray-50'}`}>
                                        <td colSpan="4" className="p-4">
                                            <div className="text-sm flex flex-wrap items-center gap-x-6 gap-y-2">
                                                <span><strong>Owner:</strong> {cluster.owner}</span>
                                                <span><strong>On Chain Balance:</strong> {cluster.onChainBalance} ssv</span>
                                                <span><strong>Burn Fee:</strong> {activeTab === 'current' ? cluster.burnFee : cluster.upcomingBurnFee} ssv</span>
                                                <div className="w-full">
                                                    <strong>Operators:</strong>{' '}
                                                    {cluster.operators.map((op) => (
                                                        <OperatorDisplay
                                                            key={op.id}
                                                            name={op.name}
                                                            id={op.id}
                                                            network={network}
                                                        />
                                                    ))}
                                                </div>
                                            </div>
                                        </td>
                                    </tr>
                                )}
                            </React.Fragment>
                        ))}
                        </tbody>
                    </table>
                </div>
            </div>
        );
    };

    useEffect(() => {
        if (networkFees.current === networkFees.upcoming && activeTab === 'simulated') {
            setActiveTab('current');
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [networkFees.current, networkFees.upcoming, activeTab]);

    return (
        <div className={`p-8 ${isDarkMode ? 'bg-gray-900 text-gray-200' : 'bg-gray-100 text-gray-800'}`}>
            <h1 className="text-4xl font-bold mb-8">Liquidation Risk Ranking</h1>

            {renderTabButtons()}

            {error && (
                <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-lg flex items-center gap-2">
                    <AlertCircle className="w-5 h-5" />
                    {error}
                </div>
            )}

            {renderContent()}

            <div className="mt-6 flex justify-between items-center">
                <div>
                    <select
                        value={itemsPerPage}
                        onChange={(e) => {
                            setItemsPerPage(Number(e.target.value));
                            setCurrentPageState(prev => ({
                                ...prev,
                                current: 1,
                                simulated: 1
                            }));
                        }}
                        className={`p-2 rounded ${isDarkMode ? 'bg-gray-800 text-white' : 'bg-gray-100 text-gray-800'}`}
                    >
                        <option value={10}>10</option>
                        <option value={20}>20</option>
                        <option value={50}>50</option>
                    </select>
                </div>
                <div className="flex items-center space-x-2">
                    <button
                        onClick={() => paginate(1)}
                        disabled={currentPage === 1}
                        className={`px-3 py-1 rounded ${isDarkMode
                            ? 'bg-gray-800 text-white hover:bg-gray-700'
                            : 'bg-gray-200 text-gray-800 hover:bg-gray-300'
                        } ${currentPage === 1 ? 'opacity-50 cursor-not-allowed' : ''}`}
                    >
                        &lt;&lt;
                    </button>
                    <button
                        onClick={() => paginate(currentPage - 1)}
                        disabled={currentPage === 1}
                        className={`px-3 py-1 rounded ${isDarkMode
                            ? 'bg-gray-800 text-white hover:bg-gray-700'
                            : 'bg-gray-200 text-gray-800 hover:bg-gray-300'
                        } ${currentPage === 1 ? 'opacity-50 cursor-not-allowed' : ''}`}
                    >
                        &lt;
                    </button>
                    <span className={`px-3 py-1 rounded ${isDarkMode ? 'bg-blue-500 text-white' : 'bg-blue-100 text-blue-800'}`}>
                        {currentPage} / {totalPages[activeTab]}
                    </span>
                    <button
                        onClick={() => paginate(currentPage + 1)}
                        disabled={currentPage === totalPages[activeTab]}
                        className={`px-3 py-1 rounded ${isDarkMode
                            ? 'bg-gray-800 text-white hover:bg-gray-700'
                            : 'bg-gray-200 text-gray-800 hover:bg-gray-300'
                        } ${currentPage === totalPages[activeTab] ? 'opacity-50 cursor-not-allowed' : ''}`}
                    >
                        &gt;
                    </button>
                    <button
                        onClick={() => paginate(totalPages[activeTab])}
                        disabled={currentPage === totalPages[activeTab]}
                        className={`px-3 py-1 rounded ${isDarkMode
                            ? 'bg-gray-800 text-white hover:bg-gray-700'
                            : 'bg-gray-200 text-gray-800 hover:bg-gray-300'
                        } ${currentPage === totalPages[activeTab] ? 'opacity-50 cursor-not-allowed' : ''}`}
                    >
                        &gt;&gt;
                    </button>
                </div>
            </div>
        </div>
    );
};

export default LiquidationRanking;