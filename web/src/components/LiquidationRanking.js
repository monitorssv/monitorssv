import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { OperatorDisplay } from './SharedComponents';

const LiquidationRanking = ({ isDarkMode, network }) => {
    const [clusters, setClusters] = useState([]);
    const [currentPage, setCurrentPage] = useState(1);
    const [itemsPerPage, setItemsPerPage] = useState(10);
    const [expandedRow, setExpandedRow] = useState(null);
    const [totalPages, setTotalPages] = useState(1);
    const [totalItems, setTotalItems] = useState(0);
    const [error, setError] = useState(null);

    useEffect(() => {
        fetchClusters();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [currentPage, itemsPerPage]);

    const fetchClusters = async () => {
        setError(null);
        try {
            const response = await fetch(`/api/get30DayLiquidationRankingClusters?page=${currentPage}&limit=${itemsPerPage}`);
            if (!response.ok) {
                throw new Error('Failed to fetch clusters');
            }
            const data = await response.json();
            setClusters(data.clusters);
            setTotalPages(data.totalPages);
            setTotalItems(data.totalItems);
        } catch (err) {
            setError('Failed to fetch clusters. Please try again later.');
        }
    };

    const toggleRowExpansion = (id) => {
        setExpandedRow(expandedRow === id ? null : id);
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
        setCurrentPage(pageNumber);
    };

    return (
        <div className={`p-8 ${isDarkMode ? 'bg-gray-900 text-gray-200' : 'bg-gray-100 text-gray-800'}`}>
            <h1 className="text-4xl font-bold mb-8">Liquidation Risk Ranking</h1>

            <div className="mb-4 text-sm text-gray-500">
                Showing {clusters.length} of {totalItems} clusters with operational runway less than 30 days
            </div>

            {error && (
                <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-lg">
                    {error}
                </div>
            )}

            <div className={`overflow-x-auto rounded-lg shadow ${isDarkMode ? 'bg-gray-800' : 'bg-white'}`}>
                <table className="w-full">
                    <thead>
                        <tr className={isDarkMode ? 'bg-gray-700' : 'bg-gray-200'}>
                            <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Rank</th>
                            <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Cluster ID</th>
                            <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Validators</th>
                            <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Runaway</th>
                        </tr>
                    </thead>
                    <tbody>
                        {clusters.map((cluster, index) => (
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
                                                <span><strong>Burn Fee:</strong> {cluster.burnFee} ssv</span>
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

            <div className="mt-6 flex justify-between items-center">
                <div>
                    <select
                        value={itemsPerPage}
                        onChange={(e) => {
                            setItemsPerPage(Number(e.target.value));
                            setCurrentPage(1);
                        }}
                        className={`p-2 rounded ${isDarkMode ? 'bg-gray-800 text-white' : 'bg-gray-100 text-gray-800'}`}
                    >
                        <option value={10}>10</option>
                        <option value={20}>20</option>
                        <option value={50}>50</option>
                    </select>
                </div>
                <div className="flex items-center space-x-2">
                    {/* Pagination buttons - same as original component */}
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
                        {currentPage} / {totalPages}
                    </span>
                    <button
                        onClick={() => paginate(currentPage + 1)}
                        disabled={currentPage === totalPages}
                        className={`px-3 py-1 rounded ${isDarkMode
                            ? 'bg-gray-800 text-white hover:bg-gray-700'
                            : 'bg-gray-200 text-gray-800 hover:bg-gray-300'
                            } ${currentPage === totalPages ? 'opacity-50 cursor-not-allowed' : ''}`}
                    >
                        &gt;
                    </button>
                    <button
                        onClick={() => paginate(totalPages)}
                        disabled={currentPage === totalPages}
                        className={`px-3 py-1 rounded ${isDarkMode
                            ? 'bg-gray-800 text-white hover:bg-gray-700'
                            : 'bg-gray-200 text-gray-800 hover:bg-gray-300'
                            } ${currentPage === totalPages ? 'opacity-50 cursor-not-allowed' : ''}`}
                    >
                        &gt;&gt;
                    </button>
                </div>
            </div>
        </div>
    );
};

export default LiquidationRanking;