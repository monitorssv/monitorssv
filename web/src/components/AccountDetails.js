import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { StatusLabel, OperatorDisplay } from './SharedComponents';
import { Link } from 'react-router-dom';

const AccountDetails = ({ isDarkMode, network }) => {
    const { address } = useParams();
    const [activeTab, setActiveTab] = useState('');
    const [accountData, setAccountData] = useState(null);

    const [clusters, setClusters] = useState([]);
    const [operators, setOperators] = useState([]);
    const [events, setEvents] = useState([]);

    const [clustersPage, setClustersPage] = useState(1);
    const [operatorsPage, setOperatorsPage] = useState(1);
    const [eventsPage, setEventsPage] = useState(1);

    const [clustersPerPage, setClustersPerPage] = useState(10);
    const [operatorsPerPage, setOperatorsPerPage] = useState(10);
    const [eventsPerPage, setEventsPerPage] = useState(10);

    const [clustersTotalItems, setClustersTotalItems] = useState(0);
    const [operatorsTotalItems, setOperatorsTotalItems] = useState(0);
    const [eventsTotalItems, setEventsTotalItems] = useState(0);

    const [clustersTotalPages, setClustersTotalPages] = useState(1);
    const [operatorsTotalPages, setOperatorsTotalPages] = useState(1);
    const [eventsTotalPages, setEventsTotalPages] = useState(1);

    const [error, setError] = useState(null);

    useEffect(() => {
        fetchAccountData();
        fetchAllData();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [address]);

    useEffect(() => {
        if (activeTab) {
            fetchDataForActiveTab();
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [clustersPage, operatorsPage, eventsPage, clustersPerPage, operatorsPerPage, eventsPerPage]);

    const fetchAccountData = async () => {
        setAccountData({ address });
    };

    const fetchAllData = async () => {
        try {
            const [clustersData, operatorsData, eventsData] = await Promise.all([
                fetchClusters(),
                fetchOperators(),
                fetchEvents()
            ]);

            // Set the active tab based on data availability
            if (clustersData.clusters && clustersData.clusters.length > 0) {
                setActiveTab('clusters');
            } else if (operatorsData.operators && operatorsData.operators.length > 0) {
                setActiveTab('operators');
            } else if (eventsData.history && eventsData.history.length > 0) {
                setActiveTab('history');
            } else {
                setActiveTab('clusters');
            }
        } catch (err) {
            console.error("Error fetching data:", err);
            setError("Failed to fetch data. Please try again later.");
        }
    };

    const fetchDataForActiveTab = () => {
        // eslint-disable-next-line default-case 
        switch (activeTab) {
            case 'clusters':
                fetchClusters();
                break;
            case 'operators':
                fetchOperators();
                break;
            case 'history':
                fetchEvents();
                break;
        }
    };

    const fetchClusters = async () => {
        setError(null);
        try {
            const response = await fetch(`/api/clusters?page=${clustersPage}&limit=${clustersPerPage}&search=${address}`);
            if (!response.ok) {
                throw new Error('Failed to fetch clusters');
            }
            const data = await response.json();
            setClusters(data.clusters);
            setClustersTotalPages(data.totalPages);
            setClustersTotalItems(data.totalItems);
            return data;
        } catch (err) {
            setError('Failed to fetch clusters. Please try again later.');
            return { clusters: [] };
        }
    };

    const fetchOperators = async () => {
        setError(null);
        try {
            const response = await fetch(`/api/operators?page=${operatorsPage}&limit=${operatorsPerPage}&search=${address}`);
            if (!response.ok) {
                throw new Error('Failed to fetch operators');
            }
            const data = await response.json();
            setOperators(data.operators);
            setOperatorsTotalPages(data.totalPages);
            setOperatorsTotalItems(data.totalItems);
            return data;
        } catch (err) {
            setError('Failed to fetch operators. Please try again later.');
            return { operators: [] };
        }
    };

    const fetchEvents = async () => {
        setError(null);
        try {
            const response = await fetch(`/api/events?page=${eventsPage}&limit=${eventsPerPage}&search=${address}`);
            if (!response.ok) {
                throw new Error('Failed to fetch events');
            }
            const data = await response.json();
            setEvents(data.history);
            setEventsTotalPages(data.totalPages);
            setEventsTotalItems(data.totalItems);
            return data;
        } catch (err) {
            setError('Failed to fetch events. Please try again later.');
            return { history: [] };
        }
    };

    const getEtherscanUrl = (type, value) => {
        const baseUrl = network === 'mainnet'
            ? 'https://etherscan.io'
            : 'https://holesky.etherscan.io';
        return `${baseUrl}/${type}/${value}`;
    };

    const truncateAddr = (addr) => {
        return `${addr.slice(0, 6)}...${addr.slice(-4)}`;
    };

    const isTabDisabled = (tabName) => {
        switch (tabName) {
            case 'clusters':
                return !clusters || clusters.length === 0;
            case 'operators':
                return !operators || operators.length === 0;
            case 'history':
                return !events || events.length === 0;
            default:
                return false;
        }
    };

    const getActiveTabData = () => {
        switch (activeTab) {
            case 'clusters':
                return clusters || [];
            case 'operators':
                return operators || [];
            case 'history':
                return events || [];
            default:
                return [];
        }
    };

    const getCurrentPage = () => {
        switch (activeTab) {
            case 'clusters':
                return clustersPage;
            case 'operators':
                return operatorsPage;
            case 'history':
                return eventsPage;
            default:
                return 1;
        }
    };

    const getTotalPages = () => {
        switch (activeTab) {
            case 'clusters':
                return clustersTotalPages;
            case 'operators':
                return operatorsTotalPages;
            case 'history':
                return eventsTotalPages;
            default:
                return 1;
        }
    };

    const getItemsPerPage = () => {
        switch (activeTab) {
            case 'clusters':
                return clustersPerPage;
            case 'operators':
                return operatorsPerPage;
            case 'history':
                return eventsPerPage;
            default:
                return 10;
        }
    };

    const setCurrentPage = (page) => {
        // eslint-disable-next-line default-case 
        switch (activeTab) {
            case 'clusters':
                setClustersPage(page);
                break;
            case 'operators':
                setOperatorsPage(page);
                break;
            case 'history':
                setEventsPage(page);
                break;
        }
    };

    const setItemsPerPage = (items) => {
        // eslint-disable-next-line default-case 
        switch (activeTab) {
            case 'clusters':
                setClustersPerPage(items);
                setClustersPage(1);
                break;
            case 'operators':
                setOperatorsPerPage(items);
                setOperatorsPage(1);
                break;
            case 'history':
                setEventsPerPage(items);
                setEventsPage(1);
                break;
        }
    };

    const paginate = (pageNumber) => setCurrentPage(pageNumber);

    const renderPagination = () => {
        const totalPages = getTotalPages();
        const currentPage = getCurrentPage();

        return (
            <div className="mt-6 flex justify-between items-center">
                <div>
                    <select
                        value={getItemsPerPage()}
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
                    <span className={`px-3 py-1 rounded ${isDarkMode ? 'bg-blue-500 text-white' : 'bg-blue-100 text-blue-800'
                        }`}>
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
        );
    };

    if (!accountData) return <div className={`p-8 ${isDarkMode ? 'bg-gray-900 text-gray-200' : 'bg-white text-gray-800'}`}>Loading...</div>;

    if (error) return <div className={`p-8 ${isDarkMode ? 'bg-gray-900 text-gray-200' : 'bg-white text-gray-800'}`}>{error}</div>;

    const activeData = getActiveTabData();

    return (
        <div className={`p-8 ${isDarkMode ? 'bg-gray-900 text-gray-200' : 'bg-white text-gray-800'}`}>
            <h1 className="text-4xl font-bold mb-8">Account Details</h1>

            <div className="mb-6 flex items-center">
                <span className="font-semibold mr-6">Address:</span>
                {accountData.address}
            </div>

            <div className="flex space-x-4 mb-6">
                {['clusters', 'operators', 'history'].map((tab) => (
                    <button
                        key={tab}
                        className={`px-4 py-2 rounded-full ${activeTab === tab
                            ? 'bg-blue-500 text-white'
                            : isDarkMode
                                ? 'bg-gray-700 text-gray-300'
                                : 'bg-gray-200 text-gray-700'
                            } ${isTabDisabled(tab)
                                ? 'opacity-50 cursor-not-allowed'
                                : 'hover:bg-blue-600 hover:text-white'
                            }`}
                        onClick={() => {
                            if (!isTabDisabled(tab)) {
                                setActiveTab(tab);
                            }
                        }}
                        disabled={isTabDisabled(tab)}
                    >
                        {tab === 'clusters' && `${clustersTotalItems} Clusters`}
                        {tab === 'operators' && `${operatorsTotalItems} Operators`}
                        {tab === 'history' && `${eventsTotalItems} Account History`}
                    </button>
                ))}
            </div>

            {error && <div className="text-red-500 mb-4">{error}</div>}

            <div className={`overflow-x-auto rounded-lg shadow ${isDarkMode ? 'bg-gray-800' : 'bg-white'}`}>
                {activeTab === 'clusters' && clusters.length > 0 && (
                    <table className="w-full">
                        <thead>
                            <tr className={isDarkMode ? 'bg-gray-700' : 'bg-gray-100'}>
                                <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Cluster ID</th>
                                <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Operators</th>
                                <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Validators</th>
                                <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Status</th>
                            </tr>
                        </thead>
                        <tbody>
                            {activeData.map((cluster) => (
                                <tr key={cluster.id} className={`border-b ${isDarkMode ? 'border-gray-700' : 'border-gray-200'}`}>
                                    <td className="p-3">
                                        <Link
                                            to={`/cluster/${cluster.id}`}
                                            className="text-blue-500 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300"
                                        >
                                            {truncateAddr(cluster.id)}
                                        </Link>
                                    </td>
                                    <td className="p-3">
                                        <div className="flex flex-wrap gap-1">
                                            {cluster.operators.map((op) => (
                                                <OperatorDisplay
                                                    key={op.id}
                                                    name={op.name}
                                                    id={op.id}
                                                    network={network}
                                                />
                                            ))}
                                        </div>
                                    </td>
                                    <td className={`p-3 ${isDarkMode ? 'text-gray-300' : 'text-gray-800'}`}>{cluster.validators}</td>
                                    <td className="p-3">
                                        <StatusLabel status={cluster.status} />
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                )}
                {activeTab === 'operators' && operators.length > 0 && (
                    <table className="w-full">
                        <thead>
                            <tr className={isDarkMode ? 'bg-gray-700' : 'bg-gray-100'}>
                                <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>ID</th>
                                <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Name</th>
                                <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Validators</th>
                                <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Operator Fee</th>
                            </tr>
                        </thead>
                        <tbody>
                            {activeData.map((item) => (
                                <tr key={item.id} className={`border-b ${isDarkMode ? 'border-gray-700' : 'border-gray-200'}`}>
                                    <td className={`p-3 ${isDarkMode ? 'text-gray-300' : 'text-gray-800'}`}>{item.id}</td>
                                    <td className={`p-3 ${isDarkMode ? 'text-gray-300' : 'text-gray-800'}`}>{item.name}</td>
                                    <td className={`p-3 ${isDarkMode ? 'text-gray-300' : 'text-gray-800'}`}>{item.validators}</td>
                                    <td className={`p-3 ${isDarkMode ? 'text-gray-300' : 'text-gray-800'}`}>{item.operatorFee} ssv</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                )}
                {activeTab === 'history' && events.length > 0 && (
                    <table className="w-full">
                        <thead>
                            <tr className={isDarkMode ? 'bg-gray-700' : 'bg-gray-100'}>
                                <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Block</th>
                                <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Transaction Hash</th>
                                <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Action</th>
                            </tr>
                        </thead>
                        <tbody>
                            {activeData.map((item, index) => (
                                <tr key={`${item.transactionHash}-${index}`} className={`border-b ${isDarkMode ? 'border-gray-700' : 'border-gray-200'}`}>
                                    <td className={`p-3 ${isDarkMode ? 'text-gray-300' : 'text-gray-800'}`}>{item.block}</td>
                                    <td className={`p-3 ${isDarkMode ? 'text-gray-300' : 'text-gray-800'}`}>
                                        <a
                                            href={getEtherscanUrl('tx', item.transactionHash)}
                                            target="_blank"
                                            rel="noopener noreferrer"
                                            className={`hover:underline ${isDarkMode ? 'text-blue-400' : 'text-blue-600'}`}
                                        >
                                            {item.transactionHash}
                                        </a>
                                    </td>
                                    <td className={`p-3 ${isDarkMode ? 'text-gray-300' : 'text-gray-800'}`}>{item.action}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                )}
                {isTabDisabled(activeTab) && (
                    <p className={`p-4 ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>No data available for this tab.</p>
                )}
            </div>
            {getActiveTabData().length > 0 && renderPagination()}
        </div>
    );
};

export default AccountDetails;