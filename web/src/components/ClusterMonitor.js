import React, { useState, useEffect } from 'react';
import { useWeb3Modal } from '@web3modal/ethers/react';
import { useWeb3ModalAccount } from '@web3modal/ethers/react';
import { useWeb3ModalProvider } from '@web3modal/ethers/react';
import { BrowserProvider } from 'ethers';
import { AlertCircle, CheckCircle, Wallet, LogOut, Send, Eye, Trash2 } from 'lucide-react';

const ClusterMonitor = ({ isDarkMode, network }) => {
    const [isOwner, setIsOwner] = useState(false);
    const [alertMethod, setAlertMethod] = useState('discord');
    const [discordWebhook, setDiscordWebhook] = useState('');
    const [telegramAccessToken, setTelegramAccessToken] = useState('');
    const [telegramChatId, setTelegramChatId] = useState('');
    const [liquidationThresholdDays, setLiquidationThresholdDays] = useState(30);
    const [tempLiquidationThreshold, setTempLiquidationThreshold] = useState('30');
    const [error, setError] = useState(null);
    const [msg, setMsg] = useState(null);
    const [block, setBlock] = useState(null);
    const [clusterInfo, setClusterInfo] = useState(null);
    const [showConfig, setShowConfig] = useState(false);
    const [configOptions, setConfigOptions] = useState({
        reportOperatorFeeChanges: false,
        reportNetworkFeeChanges: false,
        reportBlockProposals: false,
        reportMissedBlocks: false,
        reportBalanceDecrease: false,
        reportExitedButNotRemovedValidators: false,
        weeklyReport: false,
    });
    const [isMonitoring, setIsMonitoring] = useState(false);

    const [isSaving, setIsSaving] = useState(false);
    const [isDeleting, setIsDeleting] = useState(false);
    const [isTesting, setIsTesting] = useState(false);
    const [isViewing, setIsViewing] = useState(false);

    const { open } = useWeb3Modal();
    const { address, isConnected } = useWeb3ModalAccount();
    const { walletProvider } = useWeb3ModalProvider();

    useEffect(() => {
        if (isConnected && address) {
            checkOwnerStatus(address);
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [isConnected, address]);

    const checkOwnerStatus = async (address) => {
        await fetchClusterInfo(address);
    };

    const fetchClusterInfo = async (owner) => {
        setError('');
        try {
            const response = await fetch(`/api/clusterMonitorInfo?owner=${owner}`);
            if (!response.ok) {
                throw new Error('Failed to fetch clusterMonitorInfo');
            }
            const data = await response.json();
            setClusterInfo(data);
            setBlock(data.block);
            setIsMonitoring(data.isMonitoring);
            if (data.totalActiveCluster > 0) {
                setIsOwner(true);
            }
        } catch (error) {
            console.error("Error fetching cluster info:", error);
            showMessage("Failed to fetch cluster information. Please try again later.", true);
        }
    };

    const fetchMonitorConfig = async (ownerAddress, blockNumber) => {
        if (!isConnected || !isOwner || isViewing) return;
        setError('');
        setIsViewing(true);

        try {
            const ethersProvider = new BrowserProvider(walletProvider);
            const signer = await ethersProvider.getSigner();
            const message = `Signature required for cluster ownership. Block: ${blockNumber}`;
            const signature = await signer.signMessage(message);

            const response = await fetch(`/api/clusterMonitorConfig?owner=${ownerAddress}&block=${blockNumber}&signature=${signature}`)
            if (!response.ok) {
                throw new Error('Failed to fetch clusterMonitorConfig');
            }
            const data = await response.json();
            updateConfigFromResponse(data.monitorConfig);
            setBlock(data.block);
            setShowConfig(true);
        } catch (error) {
            console.error('Error fetching monitor config:', error);
            showMessage("Failed to view monitoring configuration. Please try again later.", true);
        }
        setIsViewing(false);
    };

    const deleteMonitorConfig = async () => {
        if (!isConnected || !isOwner || isDeleting) return;

        setError('');
        setIsDeleting(true);
        try {
            const ethersProvider = new BrowserProvider(walletProvider);
            const signer = await ethersProvider.getSigner();
            const message = `Signature required for cluster ownership. Block: ${block}`;
            const signature = await signer.signMessage(message);

            const response = await fetch('/api/deleteClusterMonitorConfig', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    owner: address,
                    signature,
                    block,
                })
            });
            if (!response.ok) {
                throw new Error('Failed to save configuration');
            }
            const data = await response.json();
            console.log("delete config response:", data);
            showMessage("Configuration deleted successfully!")
            setIsMonitoring(false);
            setAlertMethod('discord');
            setDiscordWebhook('');
            setTelegramAccessToken('');
            setTelegramChatId('');
            setLiquidationThresholdDays(30);
            setConfigOptions({
                reportOperatorFeeChanges: false,
                reportNetworkFeeChanges: false,
                reportBlockProposals: false,
                reportMissedBlocks: false,
                reportBalanceDecrease: false,
                reportExitedButNotRemovedValidators: false,
                weeklyReport: false,
            });
            setShowConfig(true);
        } catch (error) {
            console.error('Error fetching monitor config:', error);
            showMessage("Failed to delete monitoring configuration. Please try again later.", true);
        }

        setIsDeleting(false);
    };

    const updateConfigFromResponse = (config) => {
        setAlertMethod(config.alarm_type === 0 ? 'discord' : 'telegram');
        if (config.alarm_type === 0) {
            setDiscordWebhook(config.alarm_channel);
        }
        if (config.alarm_type === 1) {
            const [telegramAccessToken, telegramChatId] = config.alarm_channel.split(',');
            setTelegramAccessToken(telegramAccessToken);
            setTelegramChatId(telegramChatId);
        }
        setLiquidationThresholdDays(config.report_liquidation_threshold);
        setConfigOptions({
            reportOperatorFeeChanges: config.report_operator_fee_change,
            reportNetworkFeeChanges: config.report_network_fee_change,
            reportBlockProposals: network === 'holesky' ? false : config.report_propose_block,
            reportMissedBlocks: network === 'holesky' ? false : config.report_missed_block,
            reportBalanceDecrease: network === 'holesky' ? false : config.report_balance_decrease,
            reportExitedButNotRemovedValidators: config.report_exited_but_not_removed,
            weeklyReport: config.report_weekly,
        });
    };

    const handleAlertMethodChange = (method) => {
        if (!isConnected || !isOwner) return;
        setAlertMethod(method);
    };

    const handleConfigOptionChange = (option) => {
        if (!isConnected || !isOwner) return;
        setConfigOptions(prev => ({ ...prev, [option]: !prev[option] }));
    };

    const handleTestAlert = async () => {
        if (!isConnected || !isOwner || isTesting) return;
        setError('');
        setMsg('');
        setIsTesting(true);
        try {
            const response = await fetch('/api/testAlarm', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    alarm_type: alertMethod === 'discord' ? 0 : 1,
                    alarm_channel: alertMethod === 'discord' ? discordWebhook : `${telegramAccessToken},${telegramChatId}`,
                })
            });
            if (!response.ok) {
                throw new Error('Failed to testAlarm');
            }
            const data = await response.json();
            console.log("Test alarm response:", data);
            showMessage("Alarm test successful!")
        } catch (error) {
            console.error('Error testing alarm:', error);
            showMessage("Failed to test the alert. Please try again later.", true)
        }

        setIsTesting(false);
    };

    const handleSaveConfig = async () => {
        if (!isConnected || !isOwner || !walletProvider || isSaving) return;
        setError('');
        setMsg('');
        setIsSaving(true);
        try {
            const ethersProvider = new BrowserProvider(walletProvider);
            const signer = await ethersProvider.getSigner();

            const configData = {
                alarm_type: alertMethod === 'discord' ? 0 : 1,
                alarm_channel: alertMethod === 'discord' ? discordWebhook : `${telegramAccessToken},${telegramChatId}`,
                report_liquidation_threshold: liquidationThresholdDays,
                report_operator_fee_change: configOptions.reportOperatorFeeChanges,
                report_network_fee_change: configOptions.reportNetworkFeeChanges,
                report_propose_block: network === 'holesky' ? false : configOptions.reportBlockProposals,
                report_missed_block: network === 'holesky' ? false : configOptions.reportMissedBlocks,
                report_balance_decrease: network === 'holesky' ? false : configOptions.reportBalanceDecrease,
                report_exited_but_not_removed: configOptions.reportExitedButNotRemovedValidators,
                report_weekly: configOptions.weeklyReport,
            };

            const message = `Signature required for cluster ownership. Block: ${block}\n${JSON.stringify(configData)}`;
            const signature = await signer.signMessage(message);

            const response = await fetch('/api/saveClusterMonitorConfig', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    monitorConfig: JSON.stringify(configData),
                    owner: address,
                    signature,
                    block,
                })
            });
            if (!response.ok) {
                throw new Error('Failed to save configuration');
            }
            const data = await response.json();
            setIsMonitoring(true);
            console.log("Save config response:", data);
            showMessage("Configuration saved successfully!")
        } catch (error) {
            console.error('Error saving configuration:', error);
            showMessage("Failed to save configuration. Please try again later.", true);
        }

        setIsSaving(false);
    };

    const showMessage = (message, isError = false) => {
        if (isError) {
            setError(message);
            setMsg(null);
        } else {
            setMsg(message);
            setError(null);
        }

        setTimeout(() => {
            setError('');
            setMsg('');
        }, 5000);
    };

    const handleDisconnect = async () => {
        await open({ view: 'Account' });
    };

    const handleViewConfig = () => {
        if (isConnected && isOwner && block) {
            fetchMonitorConfig(address, block);
        }
    };

    const handleStartConfiguration = () => {
        if (isConnected && isOwner && clusterInfo && clusterInfo.totalActiveCluster > 0) {
            setShowConfig(true);
        }
    }

    const shouldRenderOption = (optionKey) => {
        if (network === 'holesky') {
            return !['reportBlockProposals', 'reportMissedBlocks', 'reportBalanceDecrease'].includes(optionKey);
        }
        return true;
    };

    const bgColor = isDarkMode ? 'bg-gray-900' : 'bg-gray-100';
    const textColor = isDarkMode ? 'text-white' : 'text-gray-800';
    const cardBgColor = isDarkMode ? 'bg-gray-800' : 'bg-white';
    const inputBgColor = isDarkMode ? 'bg-gray-700' : 'bg-gray-200';
    const buttonBgColor = isDarkMode ? 'bg-blue-600 hover:bg-blue-700' : 'bg-blue-500 hover:bg-blue-600';
    const disabledClass = (!isConnected || !isOwner) ? 'opacity-50 pointer-events-none' : '';

    return (
        <div className={`p-8 ${bgColor} ${textColor} min-h-screen`}>
            <div className={`${cardBgColor} rounded-3xl shadow-xl overflow-hidden w-full max-w-4xl mx-auto relative`}>
                <div className={`p-8 ${disabledClass}`}>
                    <h2 className="text-3xl font-bold mb-6 text-center">Cluster Monitoring Configuration</h2>
                    {error && (
                        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-4" role="alert">
                            <span className="block sm:inline">{error}</span>
                        </div>
                    )}
                    {msg && (
                        <div className="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded relative mb-4" role="alert">
                            <span className="block sm:inline">{msg}</span>
                        </div>
                    )}
                    <div className="space-y-6">
                        <div>
                            <h3 className="text-xl font-semibold mb-3">Alert Method</h3>
                            <div className="flex space-x-4">
                                <button
                                    onClick={() => handleAlertMethodChange('discord')}
                                    className={`flex-1 py-2 px-4 rounded-lg ${alertMethod === 'discord' ? buttonBgColor : `${inputBgColor} hover:bg-opacity-80`} transition duration-300`}
                                >
                                    Discord
                                </button>
                                <button
                                    onClick={() => handleAlertMethodChange('telegram')}
                                    className={`flex-1 py-2 px-4 rounded-lg ${alertMethod === 'telegram' ? buttonBgColor : `${inputBgColor} hover:bg-opacity-80`} transition duration-300`}
                                >
                                    Telegram
                                </button>
                            </div>
                        </div>

                        {alertMethod === 'discord' && (
                            <div>
                                <label className="block mb-2">Discord Webhook</label>
                                <input
                                    type="text"
                                    value={discordWebhook}
                                    onChange={(e) => isConnected && isOwner && setDiscordWebhook(e.target.value)}
                                    className={`w-full ${inputBgColor} rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-400`}
                                    placeholder="Enter Discord webhook URL"
                                    disabled={!isConnected || !isOwner}
                                />
                            </div>
                        )}

                        {alertMethod === 'telegram' && (
                            <div className="space-y-4">
                                <div>
                                    <label className="block mb-2">Telegram Access Token</label>
                                    <input
                                        type="text"
                                        value={telegramAccessToken}
                                        onChange={(e) => isConnected && isOwner && setTelegramAccessToken(e.target.value)}
                                        className={`w-full ${inputBgColor} rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-400`}
                                        placeholder="Enter Telegram access token"
                                        disabled={!isConnected || !isOwner}
                                    />
                                </div>
                                <div>
                                    <label className="block mb-2">Telegram Chat ID</label>
                                    <input
                                        type="text"
                                        value={telegramChatId}
                                        onChange={(e) => isConnected && isOwner && setTelegramChatId(e.target.value)}
                                        className={`w-full ${inputBgColor} rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-400`}
                                        placeholder="Enter Telegram chat ID"
                                        disabled={!isConnected || !isOwner}
                                    />
                                </div>
                            </div>
                        )}

                        <button
                            onClick={handleTestAlert}
                            className={`${buttonBgColor} text-white font-bold py-2 px-4 rounded-lg transition duration-300 flex items-center justify-center ${isTesting ? 'opacity-50 cursor-not-allowed' : ''}`}
                            disabled={!isConnected || !isOwner || isTesting}
                        >
                            <Send size={20} className="mr-2" />
                            {isTesting ? 'Testing...' : 'Test Alert Channel'}
                        </button>

                        <div className="space-y-3">
                            <h3 className="text-xl font-semibold mb-3">Monitoring Options</h3>
                            <div className="flex items-center space-x-2 mb-4">
                                <span>Alert when cluster liquidation runway is less than</span>
                                <input
                                    type="number"
                                    id="liquidationThresholdDays"
                                    value={tempLiquidationThreshold}
                                    onChange={(e) => {
                                        if (isConnected && isOwner) {
                                            setTempLiquidationThreshold(e.target.value);
                                        }
                                    }}
                                    onBlur={() => {
                                        const value = parseInt(tempLiquidationThreshold, 10);
                                        const finalValue = isNaN(value) ? 10 : Math.max(10, value);
                                        setLiquidationThresholdDays(finalValue);
                                        setTempLiquidationThreshold(finalValue.toString());
                                    }}
                                    className={`w-16 ${inputBgColor} rounded-lg px-2 py-1 text-center focus:outline-none focus:ring-2 focus:ring-blue-400`}
                                    min="10"
                                    disabled={!isConnected || !isOwner}
                                />
                                <span>days</span>
                            </div>
                            {network === 'mainnet' && <div>Alert when the validator is slashed</div>}
                            {Object.entries(configOptions).map(([key, value]) => (
                                shouldRenderOption(key) && <div key={key} className="flex items-center">
                                    <input
                                        type="checkbox"
                                        id={key}
                                        checked={value}
                                        onChange={() => handleConfigOptionChange(key)}
                                        className="mr-3"
                                        disabled={!isConnected || !isOwner}
                                    />
                                    <label htmlFor={key} className="flex-1">{key.split(/(?=[A-Z])/).join(" ")}</label>
                                </div>
                            ))}
                        </div>

                        <div className={`flex ${showConfig ? 'space-x-4' : ''}`}>
                            <button
                                onClick={handleSaveConfig}
                                className={`${isMonitoring ? 'flex-1' : 'w-full'} ${buttonBgColor} text-white font-bold py-3 px-4 rounded-lg transition duration-300 flex items-center justify-center ${isSaving ? 'opacity-50 cursor-not-allowed' : ''}`}
                                disabled={!isConnected || !isOwner || isSaving}
                            >
                                <CheckCircle size={20} className="mr-2" />
                                {isSaving ? 'Saving...' : 'Save Configuration'}
                            </button>

                            {isMonitoring && (
                                <button
                                    onClick={deleteMonitorConfig}
                                    className={`flex-1 bg-red-500 hover:bg-red-600 text-white font-bold py-3 px-4 rounded-lg transition duration-300 flex items-center justify-center ${isDeleting ? 'opacity-50 cursor-not-allowed' : ''}`}
                                    disabled={!isConnected || !isOwner || isDeleting}
                                >
                                    <Trash2 size={20} className="mr-2" />
                                    {isDeleting ? 'Deleting...' : 'Delete Configuration'}
                                </button>
                            )}
                        </div>

                        {isConnected && (
                            <button
                                onClick={handleDisconnect}
                                className="w-full bg-red-500 hover:bg-red-600 text-white font-bold py-3 px-4 rounded-lg transition duration-300 flex items-center justify-center mt-4"
                            >
                                <LogOut size={20} className="mr-2" />
                                Disconnect Wallet
                            </button>
                        )}
                    </div>
                </div>

                {!isConnected && (
                    <div className="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50">
                        <div className={`${cardBgColor} rounded-2xl shadow-2xl overflow-hidden w-full max-w-sm p-6 m-4 border border-gray-200 dark:border-gray-700`}>
                            <h2 className="text-2xl font-bold mb-4 text-center">Connect Wallet</h2>
                            <button
                                onClick={() => open()}
                                className={`w-full ${buttonBgColor} text-white font-bold py-3 px-4 rounded-lg transition duration-300 flex items-center justify-center text-base mb-4`}
                            >
                                <Wallet size={20} className="mr-2" />
                                Connect to Configure
                            </button>
                            <p className="text-center text-sm">Signature required for cluster ownership</p>
                        </div>
                    </div>
                )}

                {isConnected && clusterInfo && !showConfig && (
                    <div className="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50">
                        <div className={`${cardBgColor} rounded-2xl shadow-2xl overflow-hidden w-full max-w-sm p-6 m-4 border border-gray-200 dark:border-gray-700`}>
                            {clusterInfo.totalClusters === 0 ? (
                                <>
                                    <AlertCircle size={64} className="mx-auto mb-6 text-yellow-500" />
                                    <h2 className="text-2xl font-bold mb-4 text-center">Not Authorized</h2>
                                    <p className="text-center mb-6">You are not a cluster owner and cannot configure monitoring settings.</p>
                                </>
                            ) : !isMonitoring && clusterInfo.totalActiveCluster > 0 ? (
                                <>
                                    <h2 className="text-2xl font-bold mb-4 text-center">Owner Information</h2>
                                    <p className="text-center mb-4">Total clusters: {clusterInfo.totalClusters}</p>
                                    <p className="text-center mb-6">Active clusters: {clusterInfo.totalActiveCluster}</p>
                                    <button
                                        onClick={handleStartConfiguration}
                                        className={`w-full ${buttonBgColor} text-white font-bold py-3 px-4 rounded-lg transition duration-300 flex items-center justify-center text-base`}
                                    >
                                        Start Monitoring Configuration
                                    </button>
                                </>
                            ) : isMonitoring ? (
                                <>
                                    <h2 className="text-2xl font-bold mb-4 text-center">Monitoring Configured</h2>
                                    <p className="text-center mb-6">You have already set up monitoring for your clusters.</p>
                                    <button
                                        onClick={handleViewConfig}
                                        className={`w-full ${buttonBgColor} text-white font-bold py-3 px-4 rounded-lg transition duration-300 flex items-center justify-center ${isViewing ? 'opacity-50 cursor-not-allowed' : ''}`}
                                        disabled={!isConnected || isViewing}
                                    >
                                        <Eye size={20} className="mr-2" />
                                        {isViewing ? 'View...' : 'View Configuration'}
                                    </button>
                                </>
                            ) : null}
                            <button
                                onClick={handleDisconnect}
                                className="w-full bg-red-500 hover:bg-red-600 text-white font-bold py-3 px-4 rounded-lg transition duration-300 flex items-center justify-center mt-4"
                            >
                                <LogOut size={20} className="mr-2" />
                                Disconnect Wallet
                            </button>
                        </div>
                    </div>
                )}

                <div className={`${isDarkMode ? 'bg-gray-700' : 'bg-gray-200'} p-4 text-center text-sm ${textColor}`}>
                    {isConnected ? `Connected: ${address}` : 'Not connected'}
                </div>
            </div>
        </div>
    );
};


export default ClusterMonitor;