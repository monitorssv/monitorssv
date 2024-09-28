import React, { useState } from 'react';

export const StatusLabel = ({ status }) => (
    <span className={`
      inline-flex items-center justify-center
      w-20 px-2.5 py-0.5 rounded-full text-xs font-medium
      ${status
            ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'
            : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300'
        }
    `}>
        {status ? 'Active' : 'Inactive'}
    </span>
);

export const ValidatorStatusLabel = ({ status }) => {
    const getStatusColor = (status) => {
        switch (status) {
            case 'Active':
                return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300';
            case 'Slashed':
                return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300';
            default:
                return 'bg-gray-100 text-gray-800 dark:bg-gray-600 dark:text-gray-300';
        }
    };

    return (
        <span className={`
        inline-flex items-center justify-center
        px-2.5 py-0.5 rounded-full text-xs font-medium
        ${getStatusColor(status)}
      `}>
            {status || 'Unknown'}
        </span>
    );
};

export const OperatorDisplay = ({ name, id, network = 'mainnet' }) => {
    const baseUrl = network === 'mainnet'
        ? 'https://explorer.ssv.network/operators/'
        : 'https://holesky.explorer.ssv.network/operators/';

    const handleClick = () => {
        window.open(`${baseUrl}${id}`, '_blank');
    };

    return (
        <Tooltip content={
            <span className="text-xs font-medium text-white">
                ID: {id}
            </span>
        }>
            <span
                className="inline-block px-2 py-1 m-1 text-xs font-medium text-blue-800 bg-blue-100 rounded-full cursor-pointer transition-all duration-200 hover:bg-blue-200 dark:bg-blue-900 dark:text-blue-200 dark:hover:bg-blue-800"
                onClick={handleClick}
            >
                {name}
            </span>
        </Tooltip>
    );
};

export const Tooltip = ({ content, children }) => {
    const [isVisible, setIsVisible] = useState(false);

    return (
        <div className="relative inline-block">
            <div
                onMouseEnter={() => setIsVisible(true)}
                onMouseLeave={() => setIsVisible(false)}
            >
                {children}
            </div>
            <div
                className={`
                    absolute z-10 px-3 py-2 text-sm font-medium text-white bg-gray-900 
                    rounded-lg shadow-lg dark:bg-gray-700 tooltip transition-opacity duration-300
                    ${isVisible ? 'opacity-100' : 'opacity-0 pointer-events-none'}
                `}
                style={{
                    bottom: 'calc(100% + 5px)',
                    left: '50%',
                    transform: 'translateX(-50%)',
                    whiteSpace: 'nowrap'
                }}
            >
                {content}
                <div className="tooltip-arrow" />
            </div>
        </div>
    );
};