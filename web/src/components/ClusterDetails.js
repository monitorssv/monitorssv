import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { StatusLabel, ValidatorStatusLabel, OperatorDisplay, PublicKeyDisplay } from './SharedComponents';
import { AlertTriangle } from 'lucide-react';

const ClusterDetails = ({ isDarkMode, network }) => {
    const { id } = useParams();
    const [activeTab, setActiveTab] = useState('');
    const [clusterData, setClusterData] = useState(null);
    const [clusterPosData, setClusterPosData] = useState(null);

    const [validators, setValidators] = useState([]);
    const [blocks, setBlocks] = useState([]);
    const [events, setEvents] = useState([]);

    const [validatorsPage, setValidatorsPage] = useState(1);
    const [blocksPage, setBlocksPage] = useState(1);
    const [eventsPage, setEventsPage] = useState(1);

    const [validatorsPerPage, setValidatorsPerPage] = useState(10);
    const [blocksPerPage, setBlocksPerPage] = useState(10);
    const [eventsPerPage, setEventsPerPage] = useState(10);

    const [validatorsTotalItems, setValidatorsTotalItems] = useState(0);
    const [blocksTotalItems, setBlocksTotalItems] = useState(0);
    const [eventsTotalItems, setEventsTotalItems] = useState(0);

    const [validatorsTotalPages, setValidatorsTotalPages] = useState(1);
    const [blocksTotalPages, setBlocksTotalPages] = useState(1);
    const [eventsTotalPages, setEventsTotalPages] = useState(1);

    const [error, setError] = useState(null);

    useEffect(() => {
        fetchClusterDetails();
        fetchClusterPosData();
        fetchAllData();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [id]);

    useEffect(() => {
        if (activeTab) {
            fetchDataForActiveTab();
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [validatorsPage, blocksPage, eventsPage, validatorsPerPage, blocksPerPage, eventsPerPage]);

    const fetchClusterDetails = async () => {
        setError(null);
        try {
            const response = await fetch(`/api/clusterDetails?clusterId=${id}`);
            if (!response.ok) {
                throw new Error('Failed to fetch cluster details');
            }

            const data = await response.json();
            setClusterData(data.clusterDetails);
        } catch (error) {
            setError('Failed to fetch cluster details. Please try again later.');
        }
    };

    const fetchClusterPosData = async () => {
        if (network !== 'mainnet') {
            return;
        }
        setError(null);
        try {
            const response = await fetch(`/api/posData?clusterId=${id}`);
            if (!response.ok) {
                throw new Error('Failed to fetch cluster pos data');
            }

            const data = await response.json();
            console.log('Received cluster pos data:', data);
            setClusterPosData(data.posData);
        } catch (error) {
            setError('Failed to fetch cluster pos data. Please try again later.');
        }
    };

    const fetchAllData = async () => {
        try {
            const [validatorsData, blocksData, eventsData] = await Promise.all([
                fetchValidators(),
                fetchEvents(),
                fetchBlocks()
            ]);

            // Set the active tab based on data availability
            if (validatorsData.validators && validatorsData.validators.length > 0) {
                setActiveTab('validators');
            } else if (blocksData.blocks && blocksData.blocks.length > 0) {
                setActiveTab('blocks');
            } else if (eventsData.history && eventsData.history.length > 0) {
                setActiveTab('history');
            } else {
                setActiveTab('validators');
            }
        } catch (err) {
            console.error("Error fetching data:", err);
            setError("Failed to fetch data. Please try again later.");
        }
    };

    const fetchValidators = async () => {
        setError(null);
        try {
            const response = await fetch(`/api/validators?page=${validatorsPage}&limit=${validatorsPerPage}&search=${id}`);
            if (!response.ok) {
                throw new Error('Failed to fetch validators');
            }
            const data = await response.json();
            console.log("Validators data:", data);
            setValidators(data.validators);
            setValidatorsTotalPages(data.totalPages);
            setValidatorsTotalItems(data.totalItems);
            return data;
        } catch (err) {
            setError('Failed to fetch validators. Please try again later.');
            return { validators: [] };
        }
    };

    const fetchBlocks = async () => {
        setError(null);
        try {
            const response = await fetch(`/api/blocks?page=${blocksPage}&limit=${blocksPerPage}&clusterId=${id}`);
            if (!response.ok) {
                throw new Error('Failed to fetch blocks');
            }
            const data = await response.json();
            setBlocks(data.blocks);
            setBlocksTotalPages(data.totalPages);
            setBlocksTotalItems(data.totalItems);
            return data;
        } catch (err) {
            setError('Failed to fetch blocks. Please try again later.');
            return { blocks: [] };
        }
    };

    const fetchEvents = async () => {
        setError(null);
        try {
            const response = await fetch(`/api/events?page=${eventsPage}&limit=${eventsPerPage}&search=${id}`);
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

    const fetchDataForActiveTab = () => {
        // eslint-disable-next-line default-case 
        switch (activeTab) {
            case 'validators':
                fetchValidators();
                break;
            case 'blocks':
                fetchBlocks();
                break;
            case 'history':
                fetchEvents();
                break;
        }
    };

    const isTabDisabled = (tabName) => {
        switch (tabName) {
            case 'validators':
                return !validators || validators.length === 0;
            case 'history':
                return !events || events.length === 0;
            case 'blocks':
                return !blocks || blocks.length === 0;
            default:
                return false;
        }
    };

    const getActiveTabData = () => {
        switch (activeTab) {
            case 'validators':
                return validators || [];
            case 'history':
                return events || [];
            case 'blocks':
                return blocks || [];
            default:
                return [];
        }
    };

    const getCurrentPage = () => {
        switch (activeTab) {
            case 'validators':
                return validatorsPage;
            case 'blocks':
                return blocksPage;
            case 'history':
                return eventsPage;
            default:
                return 1;
        }
    };

    const getTotalPages = () => {
        switch (activeTab) {
            case 'validators':
                return validatorsTotalPages;
            case 'blocks':
                return blocksTotalPages;
            case 'history':
                return eventsTotalPages;
            default:
                return 1;
        }
    };

    const getItemsPerPage = () => {
        switch (activeTab) {
            case 'validators':
                return validatorsPerPage;
            case 'blocks':
                return blocksPerPage;
            case 'history':
                return eventsPerPage;
            default:
                return 10;
        }
    };

    const setCurrentPage = (page) => {
        // eslint-disable-next-line default-case 
        switch (activeTab) {
            case 'validators':
                setValidatorsPage(page);
                break;
            case 'blocks':
                setBlocksPage(page);
                break;
            case 'history':
                setEventsPage(page);
                break;
        }
    };

    const setItemsPerPage = (items) => {
        // eslint-disable-next-line default-case 
        switch (activeTab) {
            case 'validators':
                setValidatorsPerPage(items);
                setValidatorsPage(1);
                break;
            case 'blocks':
                setBlocksPerPage(items);
                setBlocksPage(1);
                break;
            case 'history':
                setEventsPerPage(items);
                setEventsPage(1);
                break;
        }
    };

    const paginate = (pageNumber) => setCurrentPage(pageNumber);

    const formatRunaway = (blocks) => {
        const days = Math.floor(blocks / 7200);
        const hours = Math.floor((blocks - days * 7200) * 12 / 3600);
        if (days === 0 && hours === 0) {
            return 'liquidatable';
        }
        return `${days}d ${hours}h`;
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

    const getSSVExploereUrl = (type, value) => {
        const baseUrl = network === 'mainnet'
            ? 'https://explorer.ssv.network'
            : 'https://holesky.explorer.ssv.network';
        return `${baseUrl}/${type}/${value}`;
    };

    const renderDetailsSection = () => (
        <div className={`p-4 rounded-lg mb-6 space-y-1 ${isDarkMode ? 'bg-gray-800' : 'bg-gray-200'}`}>
            <div className="mb-1">
                <span className="font-semibold mr-2">Active:</span>
                <StatusLabel status={clusterData.active} />
            </div>
            <div className="mb-1">
                <span className="font-semibold mr-2">Validator Count:</span>
                {clusterData.validatorCount}
            </div>
            <div className="mb-1">
                <span className="font-semibold mr-2">On Chain Balance:</span>
                {clusterData.onChainBalance}
                <span className="ml-1 text-sm text-gray-400">ssv</span>
            </div>
            <div className="mb-1">
                <span className="font-semibold mr-2">Cluster Burn Fee:</span>
                {(clusterData.burnFee * clusterData.validatorCount).toFixed(2)}
                <span className="ml-1 text-sm text-gray-400">ssv/year</span>
            </div>
            <div className="mb-1 flex items-center">
                <span className="font-semibold mr-2">Operational Runway:</span>
                <span >
                    {clusterData.validatorCount === 0 || !clusterData.active ? '--' : formatRunaway(clusterData.operationalRunaway)}
                </span>
                {clusterData.burnFee !== clusterData.upcomingBurnFee && (
                    <div className="relative group ml-2">
                        <AlertTriangle className="w-4 h-4 text-yellow-500" />
                        <div className={`absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-2 
                            ${isDarkMode
                                ? 'bg-gray-900 text-white'
                                : 'bg-white text-gray-800 shadow-lg border border-gray-200'
                            } 
                            text-sm rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap z-10`}>
                            <div>Forecasted Cluster Burn Fee: {(clusterData.upcomingBurnFee * clusterData.validatorCount).toFixed(2)} ssv/year</div>
                            <div>Forecasted Operational Runway: {formatRunaway(clusterData.upcomingOperationalRunaway)}</div>
                            <div className={`absolute bottom-0 left-1/2 transform -translate-x-1/2 translate-y-1 border-4 border-transparent 
                                ${isDarkMode
                                    ? 'border-t-gray-900'
                                    : 'border-t-white'
                                }`}>
                            </div>
                        </div>
                    </div>
                )}
            </div>
            <div className="mb-1">
                <span className="font-semibold mr-2">Owner:</span>
                {clusterData.owner}
            </div>
            <div className="mb-1">
                <span className="font-semibold mr-2">Fee Recipient Address:</span>
                {clusterData.feeRecipientAddress}
            </div>
            <div className="mb-1">
                <span className="font-semibold mr-2">Cluster ID:</span>
                {clusterData.id}
            </div>
            <div className="flex items-start">
                <span className="font-semibold w-36 mt-1">Operators:</span>
                <div className="flex flex-wrap -m-0.5">
                    {clusterData.operators.map((op) => (
                        <div key={op.id} className="m-0.5">
                            <OperatorDisplay
                                name={op.name}
                                id={op.id}
                                network={network}
                            />
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );

    const renderRewardsSection = () => {
        if (!clusterPosData) return null;
        return (
            <div className={`p-4 rounded-lg mb-6 grid grid-cols-2 md:grid-cols-4 gap-4 ${isDarkMode ? 'bg-gray-800 text-gray-200' : 'bg-gray-200 text-gray-800'}`}>
                <div>
                    <div className="text-sm text-gray-400">Total Proposed Blocks</div>
                    <div className="text-xl font-bold">{clusterPosData.totalProposedBlocks ?? 'N/A'}</div>
                </div>
                <div>
                    <div className="text-sm text-gray-400">Total Missed Blocks</div>
                    <div className="text-xl font-bold">{clusterPosData.totalMissedBlocks ?? 'N/A'}</div>
                </div>
                <div>
                    <div className="text-sm text-gray-400">Offline Validators</div>
                    <div className="text-xl font-bold">{clusterPosData.totalOfflineCount ?? 'N/A'}</div>
                </div>
                <div>
                    <div className="text-sm text-gray-400">Pending Removal Validators</div>
                    <div className="text-xl font-bold">{clusterPosData.pendingRemovalCount}</div>
                </div>
            </div>
        );
    };

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

    if (!clusterData) return <div className={`p-8 ${isDarkMode ? 'bg-gray-900 text-gray-200' : 'bg-white text-gray-800'}`}>Loading...</div>;

    if (error) return <div className={`p-8 ${isDarkMode ? 'bg-gray-900 text-gray-200' : 'bg-white text-gray-800'}`}>{error}</div>;

    const activeData = getActiveTabData();

    return (
        <div className={`p-8 ${isDarkMode ? 'bg-gray-900 text-gray-200' : 'bg-gray-100 text-gray-800'}`}>
            <h1 className="text-4xl font-bold mb-8">Cluster</h1>

            {renderDetailsSection()}
            {network === 'mainnet' && renderRewardsSection()}

            <div className="flex space-x-4 mb-6">
                {['validators', 'blocks', 'history'].map((tab) => (
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
                        {tab === 'validators' && `${validatorsTotalItems} Validators`}
                        {tab === 'blocks' && `${blocksTotalItems} Blocks`}
                        {tab === 'history' && `${eventsTotalItems} Account History`}
                    </button>
                ))}
            </div>

            {error && <div className="text-red-500 mb-4">{error}</div>}

            <div className={`overflow-x-auto rounded-lg shadow ${isDarkMode ? 'bg-gray-800' : 'bg-white'}`}>
                {activeTab === 'validators' && validators.length > 0 && (
                    <table className="w-full">
                        <thead>
                            <tr className={isDarkMode ? 'bg-gray-700' : 'bg-gray-100'}>
                                <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Public Key</th>
                                <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Status</th>
                                <th className={`p-3 text-center font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Online</th>
                            </tr>
                        </thead>
                        <tbody>
                            {activeData.map((validator, index) => (
                                <tr key={validator.publicKey} className={`border-b ${isDarkMode ? 'border-gray-700' : 'border-gray-200'} ${isDarkMode
                                    ? index % 2 === 0 ? 'bg-gray-800' : 'bg-gray-750'
                                    : index % 2 === 0 ? 'bg-gray-50' : 'bg-white'
                                    }`}>
                                    <td className={`p-3 ${isDarkMode ? 'text-gray-300' : 'text-gray-800'}`}>
                                        <PublicKeyDisplay
                                            publicKey={validator.publicKey}
                                            beaconchainLink={getBeaconscanUrl('validator', `0x${validator.publicKey}`)}
                                            explorerssvLink={getSSVExploereUrl('validators', `0x${validator.publicKey}`)}
                                            isDarkMode={isDarkMode}
                                            isTruncate={false}
                                        />
                                    </td>
                                    <td className="p-3"><ValidatorStatusLabel status={validator.status} /></td>
                                    <td className="p-3">
                                        <div className="flex justify-center items-center">
                                            {network === 'mainnet' ? (
                                                <span className={`inline-block w-3 h-3 rounded-full ${validator.online ? 'bg-green-500' : 'bg-gray-500'
                                                    }`}></span>
                                            ) : (
                                                <span>-</span>
                                            )}
                                        </div>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                )}
                {activeTab === 'blocks' && blocks.length > 0 && (
                    <table className="w-full">
                        <thead>
                            <tr className={isDarkMode ? 'bg-gray-700' : 'bg-gray-100'}>
                                <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Epoch</th>
                                <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Slot</th>
                                <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Block</th>
                                <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Proposer</th>
                                <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Status</th>
                            </tr>
                        </thead>
                        <tbody>
                            {activeData.map((item) => (
                                <tr key={item.id} className={`border-b ${isDarkMode ? 'border-gray-700' : 'border-gray-200'}`}>
                                    <td className={`p-3 ${isDarkMode ? 'text-gray-300' : 'text-gray-800'}`}>{item.epoch}</td>
                                    <td className="p-2">
                                        <a
                                            href={getBeaconscanUrl('slot', item.slot)}
                                            target="_blank"
                                            rel="noopener noreferrer"
                                            className="text-blue-500 hover:underline"
                                        >
                                            {item.slot}
                                        </a>
                                    </td>
                                    <td className={`p-3 ${isDarkMode ? 'text-gray-300' : 'text-gray-800'}`}>
                                        {item.blockNumber !== 0 ? (
                                            <a
                                                href={getEtherscanUrl('block', item.blockNumber)}
                                                target="_blank"
                                                rel="noopener noreferrer"
                                                className={`hover:underline ${isDarkMode ? 'text-blue-400' : 'text-blue-600'}`}
                                            >
                                                {item.blockNumber}
                                            </a>
                                        ) : (
                                            '--'
                                        )}
                                    </td>
                                    <td className={`p-3 ${isDarkMode ? 'text-gray-300' : 'text-gray-800'}`}>{item.proposer}</td>
                                    <td className={`p-3 ${isDarkMode ? 'text-gray-300' : 'text-gray-800'}`}>
                                        {item.blockNumber !== 0 ? (
                                            <span className="text-green-500 font-bold">Proposed</span>
                                        ) : (
                                            <span className="text-red-500 font-bold">Missed</span>
                                        )}
                                    </td>
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
                    </table >
                )}
                {
                    isTabDisabled(activeTab) && (
                        <p className={`p-4 ${isDarkMode ? 'text-gray-400' : 'text-gray-600'}`}>No data available for this tab.</p>
                    )
                }
            </div >
            {getActiveTabData().length > 0 && renderPagination()}
        </div >
    );
};

export default ClusterDetails;