import React, { useState, useEffect } from 'react';
// import { StatusLabel } from './SharedComponents';
import { Link } from 'react-router-dom';

const Operators = ({ isDarkMode }) => {
    const [operators, setOperators] = useState([]);
    const [searchTerm, setSearchTerm] = useState('');
    const [currentPage, setCurrentPage] = useState(1);
    const [itemsPerPage, setItemsPerPage] = useState(10);
    const [expandedRow, setExpandedRow] = useState(null);
    const [totalPages, setTotalPages] = useState(1);
    const [totalItems, setTotalItems] = useState(0);
    const [error, setError] = useState(null);
    const [shouldFetch, setShouldFetch] = useState(true);

    useEffect(() => {
        if (shouldFetch || searchTerm === "") {
            fetchOperators();
            setShouldFetch(false);
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [currentPage, itemsPerPage, searchTerm]);

    const fetchOperators = async () => {
        setError(null);
        try {
            const search = searchTerm ? `&search=${searchTerm}` : '';
            const response = await fetch(`/api/operators?page=${currentPage}&limit=${itemsPerPage}${search}`);
            if (!response.ok) {
                throw new Error('Failed to fetch operators');
            }
            const data = await response.json();
            setOperators(data.operators);
            setTotalPages(data.totalPages);
            setTotalItems(data.totalItems);
        } catch (err) {
            setError('Failed to fetch operators. Please try again later.');
        }
    };

    const handleSearch = (e) => {
        e.preventDefault();
        setCurrentPage(1);
        fetchOperators();
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

    const toggleRowExpansion = (id) => {
        setExpandedRow(expandedRow === id ? null : id);
    };

    const truncateAddress = (address) => {
        if (address === "") return;
        return `${address.slice(0, 6)}...${address.slice(-4)}`;
    };

    return (
        <div className={`p-8 ${isDarkMode ? 'bg-gray-900 text-gray-200' : 'bg-gray-100 text-gray-800'}`}>
            <h1 className="text-4xl font-bold mb-8">Operators</h1>

            <form onSubmit={handleSearch} className="mb-6 flex">
                <input
                    type="text"
                    placeholder="Search operator by name, owner address or operator id"
                    className={`flex-grow p-3 rounded-l-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 ${isDarkMode ? 'bg-gray-800 text-white' : 'bg-white text-black'}`}
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
                Showing {operators.length} of {totalItems} operators
            </div>

            <div className={`overflow-x-auto rounded-lg shadow ${isDarkMode ? 'bg-gray-800' : 'bg-white'}`}>
                <table className="w-full">
                    <thead>
                        <tr className={isDarkMode ? 'bg-gray-700' : 'bg-gray-200'}>
                            <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Name</th>
                            <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Owner</th>
                            <th className={`p-3 text-left font-semibold ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>Validators</th>
                        </tr>
                    </thead>
                    <tbody>
                        {operators.map((operator) => (
                            <React.Fragment key={operator.id}>
                                <tr
                                    className={`border-b ${isDarkMode ? 'border-gray-700' : 'border-gray-200'} hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-150 cursor-pointer`}
                                    onClick={() => toggleRowExpansion(operator.id)}
                                >
                                    <td className="p-3 flex items-center">
                                        <span className="mr-2">{operator.name}</span>
                                        {operator.removed ? (
                                            <span className="px-1.5 py-0.5 text-xs font-medium rounded-full bg-red-500 text-white">removed</span>
                                        ) : (
                                            !operator.privacy ? (
                                                <span className="px-1.5 py-0.5 text-xs font-medium rounded-full bg-green-500 text-white">public</span>
                                            ) : (
                                                <span className="px-1.5 py-0.5 text-xs font-medium rounded-full bg-blue-500 text-white">private</span>
                                            )
                                        )}
                                    </td>
                                    <td className="p-3">
                                        <Link
                                            to={`/account/${operator.owner}`}
                                            className={`hover:underline ${isDarkMode ? 'text-blue-400 hover:text-blue-300' : 'text-blue-600 hover:text-blue-800'}`}
                                            onClick={(e) => e.stopPropagation()}
                                        >
                                            {truncateAddress(operator.owner)}
                                        </Link>
                                    </td>
                                    <td className={`p-3 ${isDarkMode ? 'text-gray-300' : 'text-gray-800'}`}>{operator.validators}</td>
                                </tr>
                                {expandedRow === operator.id && (
                                    <tr className={`border-b ${isDarkMode ? 'border-gray-700 bg-gray-750' : 'border-gray-200 bg-gray-50'}`}>
                                        <td colSpan="4" className="p-4">
                                            <div className="text-sm flex flex-wrap items-center gap-x-6 gap-y-2">
                                                <span><strong>Operator ID:</strong> {operator.id}</span>
                                                <span><strong>Operator Fee:</strong> {operator.operatorFee} ssv</span>
                                                <span><strong>Operator Earnings:</strong> {operator.operatorEarnings} ssv</span>
                                                <div className="w-full">
                                                    <strong>Whitelisted Addresses:</strong>{' '}
                                                    {operator.whitelistedAddress.map((address, index) => (
                                                        <React.Fragment key={address}>
                                                            <Link
                                                                to={`/account/${address}`}
                                                                className={`hover:underline ${isDarkMode ? 'text-blue-400 hover:text-blue-300' : 'text-blue-600 hover:text-blue-800'}`}
                                                                onClick={(e) => e.stopPropagation()}
                                                            >
                                                                {truncateAddress(address)}
                                                            </Link>
                                                            {index < operator.whitelistedAddress.length - 1 ? ', ' : ''}
                                                        </React.Fragment>
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

export default Operators;