import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { ValidatorStatusLabel, OperatorDisplay, PublicKeyDisplay } from './SharedComponents';

const Validators = ({ isDarkMode, network }) => {
    const [validators, setValidators] = useState([]);
    const [searchTerm, setSearchTerm] = useState('');
    const [currentPage, setCurrentPage] = useState(1);
    const [itemsPerPage, setItemsPerPage] = useState(10);
    const [totalPages, setTotalPages] = useState(1);
    const [totalItems, setTotalItems] = useState(0);
    const [error, setError] = useState(null);
    const [shouldFetch, setShouldFetch] = useState(true);

    useEffect(() => {
        if (shouldFetch || searchTerm === "") {
            fetchValidators();
            setShouldFetch(false);
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [currentPage, itemsPerPage, searchTerm]);

    const fetchValidators = async () => {
        setError(null);
        try {
            const search = searchTerm ? `&search=${searchTerm}` : '';
            const response = await fetch(`/api/validators?page=${currentPage}&limit=${itemsPerPage}${search}`);
            if (!response.ok) {
                throw new Error('Failed to fetch operators');
            }
            const data = await response.json();
            setValidators(data.validators);
            setTotalPages(data.totalPages);
            setTotalItems(data.totalItems);
        } catch (err) {
            setError('Failed to fetch validators. Please try again later.');
        }
    };

    const handleSearch = (e) => {
        e.preventDefault();
        setCurrentPage(1);
        fetchValidators();
    };

    const handleSearchInputChange = (e) => {
        setSearchTerm(e.target.value);
    };

    const handleKeyPress = (e) => {
        if (e.key === 'Enter') {
            handleSearch(e);
        }
    };

    const paginate = (pageNumber) => {
        setCurrentPage(pageNumber);
        setShouldFetch(true);
    }

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

    const truncateAddr = (addr) => {
        return `${addr.slice(0, 6)}...${addr.slice(-4)}`;
    };

    return (
        <div className={`p-8 ${isDarkMode ? 'bg-gray-900 text-gray-200' : 'bg-gray-100 text-gray-800'}`}>
            <h1 className="text-4xl font-bold mb-8">Validators</h1>

            <form onSubmit={handleSearch} className="mb-6">
                <input
                    type="text"
                    placeholder="Search validator by owner address, public key or cluster id"
                    className={`w-full p-3 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 ${isDarkMode ? 'bg-gray-800 text-white' : 'bg-white text-black'
                        }`}
                    value={searchTerm}
                    onChange={handleSearchInputChange}
                    onKeyPress={handleKeyPress}
                />
            </form>

            {error && (
                <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-lg">
                    {error}
                </div>
            )}

            <div className="mb-4 text-sm text-gray-500">
                Showing {validators.length} of {totalItems} validators
            </div>

            <div className={`overflow-x-auto rounded-lg shadow ${isDarkMode ? 'bg-gray-800' : 'bg-white'}`}>
                <table className="w-full">
                    <thead>
                        <tr className={isDarkMode ? 'bg-gray-700' : 'bg-gray-200'}>
                            <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Public Key</th>
                            <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Owner</th>
                            <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Operators</th>
                            <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Cluster ID</th>
                            <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Status</th>
                            <th className={`p-3 text-center font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Online</th>
                        </tr>
                    </thead>
                    <tbody>
                        {validators.map((validator, index) => (
                            <tr key={validator.publicKey} className={`border-b ${isDarkMode ? 'border-gray-700' : 'border-gray-200'}`}>
                                <td className={`p-3 ${isDarkMode ? 'text-gray-300' : 'text-gray-800'}`}>
                                    <PublicKeyDisplay
                                        publicKey={validator.publicKey}
                                        beaconchainLink={getBeaconscanUrl('validator', `0x${validator.publicKey}`)}
                                        explorerssvLink={getSSVExploereUrl('validators', `0x${validator.publicKey}`)}
                                        isDarkMode={isDarkMode}
                                        isTruncate={true}
                                    />
                                </td>

                                <td className="p-3">
                                    <Link
                                        to={`/account/${validator.owner}`}
                                        className={`hover:underline ${isDarkMode ? 'text-blue-400 hover:text-blue-300' : 'text-blue-600 hover:text-blue-800'}`}
                                    >
                                        {truncateAddr(validator.owner)}
                                    </Link>
                                </td>
                                <td className="p-3">
                                    <div className="flex flex-wrap gap-1">
                                        {validator.operators.map((op) => (
                                            <OperatorDisplay
                                                key={op.id}
                                                name={op.name}
                                                id={op.id}
                                                network={network}
                                            />
                                        ))}
                                    </div>
                                </td>
                                <td className="p-3">
                                    <Link
                                        to={`/cluster/${validator.clusterId}`}
                                        className={`hover:underline ${isDarkMode ? 'text-blue-400 hover:text-blue-300' : 'text-blue-600 hover:text-blue-800'
                                            }`}
                                    >
                                        {truncateAddr(validator.clusterId)}
                                    </Link>
                                </td>
                                <td className="p-3">
                                    <ValidatorStatusLabel status={validator.status} />
                                </td>
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
            </div>

            <div className="mt-6 flex justify-between items-center">
                <div>
                    <select
                        value={itemsPerPage}
                        onChange={(e) => {
                            setItemsPerPage(Number(e.target.value));
                            setCurrentPage(1);
                            setShouldFetch(true)
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
        </div>
    );
};

export default Validators;