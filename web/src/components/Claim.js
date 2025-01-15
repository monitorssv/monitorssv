import React, { useState, useEffect } from 'react';
import { Search, Wallet, LogOut, CheckCircle, Loader } from 'lucide-react';
import { useWeb3Modal } from '@web3modal/ethers/react';
import { useWeb3ModalAccount } from '@web3modal/ethers/react';
import { useWeb3ModalProvider } from '@web3modal/ethers/react';
import { JsonRpcProvider, BrowserProvider, Contract, toBigInt } from 'ethers';

const Claim = ({ isDarkMode }) => {
    const [recipient, setRecipient] = useState('');
    const [eligibleRewards, setEligibleRewards] = useState(0);
    const [claimedRewards, setClaimedRewards] = useState(0);
    const [rewardsToClaim, setRewardsToClaim] = useState(0);
    const [isLoading, setIsLoading] = useState(false);
    const [isClaiming, setIsClaiming] = useState(false);
    const [txHash, setTxHash] = useState('');
    const [showTxModal, setShowTxModal] = useState(false);
    const [expectedMerkleRoot, setExpectedMerkleRoot] = useState('');
    const [merkleProof, setMerkleProof] = useState([]);
    const [errorMessage, setErrorMessage] = useState('');
    const [safeTransactionStatus, setSafeTransactionStatus] = useState(null);
    const { open } = useWeb3Modal();
    const { address, isConnected } = useWeb3ModalAccount();
    const { walletProvider } = useWeb3ModalProvider();

    const contractAddress = '0xe16d6138B1D2aD4fD6603ACdb329ad1A6cD26D9f';
    const abi = [{ "inputs": [{ "internalType": "address", "name": "token_", "type": "address" }], "stateMutability": "nonpayable", "type": "constructor" }, { "anonymous": false, "inputs": [{ "indexed": false, "internalType": "address", "name": "account", "type": "address" }, { "indexed": false, "internalType": "uint256", "name": "amount", "type": "uint256" }], "name": "Claimed", "type": "event" }, { "anonymous": false, "inputs": [{ "indexed": false, "internalType": "bytes32", "name": "oldMerkleRoot", "type": "bytes32" }, { "indexed": false, "internalType": "bytes32", "name": "newMerkleRoot", "type": "bytes32" }], "name": "MerkelRootUpdated", "type": "event" }, { "anonymous": false, "inputs": [{ "indexed": true, "internalType": "address", "name": "previousOwner", "type": "address" }, { "indexed": true, "internalType": "address", "name": "newOwner", "type": "address" }], "name": "OwnershipTransferred", "type": "event" }, { "inputs": [{ "internalType": "address", "name": "account", "type": "address" }, { "internalType": "uint256", "name": "cumulativeAmount", "type": "uint256" }, { "internalType": "bytes32", "name": "expectedMerkleRoot", "type": "bytes32" }, { "internalType": "bytes32[]", "name": "merkleProof", "type": "bytes32[]" }], "name": "claim", "outputs": [], "stateMutability": "nonpayable", "type": "function" }, { "inputs": [{ "internalType": "address", "name": "", "type": "address" }], "name": "cumulativeClaimed", "outputs": [{ "internalType": "uint256", "name": "", "type": "uint256" }], "stateMutability": "view", "type": "function" }, { "inputs": [], "name": "merkleRoot", "outputs": [{ "internalType": "bytes32", "name": "", "type": "bytes32" }], "stateMutability": "view", "type": "function" }, { "inputs": [], "name": "owner", "outputs": [{ "internalType": "address", "name": "", "type": "address" }], "stateMutability": "view", "type": "function" }, { "inputs": [], "name": "renounceOwnership", "outputs": [], "stateMutability": "nonpayable", "type": "function" }, { "inputs": [{ "internalType": "bytes32", "name": "merkleRoot_", "type": "bytes32" }], "name": "setMerkleRoot", "outputs": [], "stateMutability": "nonpayable", "type": "function" }, { "inputs": [], "name": "token", "outputs": [{ "internalType": "address", "name": "", "type": "address" }], "stateMutability": "view", "type": "function" }, { "inputs": [{ "internalType": "address", "name": "newOwner", "type": "address" }], "name": "transferOwnership", "outputs": [], "stateMutability": "nonpayable", "type": "function" }];

    useEffect(() => {
        if (isConnected && address) {
            setRecipient(address);
        }
    }, [isConnected, address]);

    useEffect(() => {
        if (recipient) {
            handleSearch();
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [recipient]);

    const getClaimedFromContract = async (address) => {
        setIsLoading(true);
        setErrorMessage('');

        try {
            const provider = new JsonRpcProvider("https://ethereum.publicnode.com");
            const contract = new Contract(contractAddress, abi, provider);
            const claimed = await contract.cumulativeClaimed(address);
            return claimed;
        } catch (error) {
            console.error("Failed to fetch claimed amount from contract:", error);
            throw new Error('Failed to fetch claimed amount from contract');
        }
    };

    const handleSearch = async (e) => {
        if (e) e.preventDefault();
        if (!recipient) return;
        setIsLoading(true);
        setErrorMessage('');

        try {
            console.log("Searching for:", recipient);
            const response = await fetch(`/api/claim?account=${recipient}`);
            if (!response.ok) {
                throw new Error('Failed to fetch rewards');
            }
            const data = await response.json();

            console.log("============", data);
            const cumulativeAmount = toBigInt(data.ssvRewardInfo.cumulativeAmount);
            let claimed = 0n;
            if (cumulativeAmount !== 0n) {
                claimed = await getClaimedFromContract(recipient);
                console.log("======getClaimedFromContract======", claimed);
            }

            const rewardsToClaim = cumulativeAmount - claimed;

            setEligibleRewards(cumulativeAmount);
            setClaimedRewards(claimed);
            setRewardsToClaim(rewardsToClaim > 0n ? rewardsToClaim : '0');
            setExpectedMerkleRoot(data.ssvRewardInfo.expectedMerkleRoot);
            setMerkleProof(data.ssvRewardInfo.merkleProof);
        } catch (error) {
            console.error("Failed to fetch data:", error);
            setErrorMessage('Failed to fetch rewards data. Please try again.');
        } finally {
            setIsLoading(false);
        }
    };

    const handleClaim = async () => {
        if (!walletProvider || !isConnected || rewardsToClaim <= 0) return;

        setIsClaiming(true);
        setErrorMessage('');

        try {
            const ethersProvider = new BrowserProvider(walletProvider);
            const signer = await ethersProvider.getSigner();

            const contract = new Contract(contractAddress, abi, signer);
            console.log("============", recipient);
            console.log("============", eligibleRewards);
            console.log("============", expectedMerkleRoot);
            console.log("============", merkleProof);
            const isSafe = await checkIsSafeWallet(signer.address);
            console.log("is Safe wallet", isSafe);

            if (isSafe) {
                contract.claim(
                    recipient,
                    eligibleRewards,
                    expectedMerkleRoot,
                    merkleProof
                ).catch(error => {
                    console.error("Failed to claim rewards:", error);
                    setErrorMessage(error.message || 'Failed to claim rewards. Please try again.');
                    setShowTxModal(false);
                });

                console.log("propose safe transaction");

                await new Promise(resolve => setTimeout(resolve, 5000));

                setSafeTransactionStatus({
                    status: 'pending',
                    message: 'Transaction proposal created. Please sign it in Safe Wallet!',
                    link: `https://app.safe.global/transactions/queue?safe=eth:${recipient}`
                });

                setShowTxModal(true);
                setIsClaiming(false);
                return;
            }

            const tx = await contract.claim(
                recipient,
                eligibleRewards,
                expectedMerkleRoot,
                merkleProof
            );

            setTxHash(tx.hash);
            await tx.wait();
            console.log("Transaction confirmed");
            setClaimedRewards(eligibleRewards);
            setRewardsToClaim(0);
            setShowTxModal(true);
        } catch (error) {
            console.error("Failed to claim rewards:", error);
            setErrorMessage(error.message || 'Failed to claim rewards. Please try again.');
        }
        finally {
            setIsClaiming(false);
        }
    };

    const checkIsSafeWallet = async (address) => {
        try {
            const provider = new JsonRpcProvider("https://ethereum.publicnode.com");
            const safeCode = await provider.getCode(address);
            return safeCode !== '0x';
        } catch (error) {
            console.error("Error checking Safe wallet:", error);
            return false;
        }
    };

    const ErrorMessage = ({ errorMessage }) => {
        const truncatedMessage = errorMessage.length > 80
            ? errorMessage.slice(0, 100) + '...'
            : errorMessage;

        return (
            errorMessage && (
                <div className="mb-4 text-red-500 text-sm">
                    {truncatedMessage}
                </div>
            )
        );
    };

    const truncateSafeLink = (hash) => {
        if (typeof hash !== 'string') {
            return hash;
        }
        return `${hash.slice(0, 60)}...${hash.slice(-7)}`;
    };

    const truncateHash = (hash) => {
        if (typeof hash !== 'string' || hash.length < 58) {
            return hash;
        }
        return `${hash.slice(0, 28)}...${hash.slice(-29)}`;
    };

    const bgColor = isDarkMode ? 'bg-gray-900' : 'bg-gray-100';
    const textColor = isDarkMode ? 'text-white' : 'text-gray-800';
    const cardBgColor = isDarkMode ? 'bg-gray-800' : 'bg-white';
    const inputBgColor = isDarkMode ? 'bg-gray-700' : 'bg-gray-200';
    const buttonBgColor = isDarkMode ? 'bg-blue-600 hover:bg-blue-700' : 'bg-blue-500 hover:bg-blue-600';

    const TransactionModal = ({ txHash, safeTransactionStatus, isDarkMode, closeModal }) => {
        console.log('Modal rendered with:', { txHash, safeTransactionStatus });

        if (safeTransactionStatus) {
            return (
                <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
                    <div className={`${isDarkMode ? 'bg-gray-800' : 'bg-white'} p-6 rounded-lg shadow-xl`}>
                        <h3 className={`text-xl font-bold mb-4 ${isDarkMode ? 'text-white' : 'text-black'}`}>
                            Transaction Proposed
                        </h3>
                        <p className={`mb-4 ${isDarkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                            {safeTransactionStatus.message}
                        </p>
                        <p className={`mb-4 break-all ${isDarkMode ? 'text-blue-300' : 'text-blue-600'}`}>
                            {truncateSafeLink(safeTransactionStatus.link)}
                        </p>
                        <div className="flex justify-between">
                            <button
                                onClick={closeModal}
                                className={`px-4 py-2 rounded ${isDarkMode ? 'bg-gray-600 text-white' : 'bg-gray-200 text-black'}`}
                            >
                                Close
                            </button>
                            <a
                                href={safeTransactionStatus.link}
                                target="_blank"
                                rel="noopener noreferrer"
                                className={`px-4 py-2 rounded ${isDarkMode ? 'bg-blue-600 text-white' : 'bg-blue-500 text-white'}`}
                            >
                                View in Safe Wallet
                            </a>
                        </div>
                    </div>
                </div>
            );
        }

        return (
            <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
                <div className={`${isDarkMode ? 'bg-gray-800' : 'bg-white'} p-6 rounded-lg shadow-xl`}>
                    <h3 className={`text-xl font-bold mb-4 ${isDarkMode ? 'text-white' : 'text-black'}`}>
                        Transaction Sent
                    </h3>
                    <p className={`mb-4 ${isDarkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                        Transaction has been sent. Transaction hash:
                    </p>
                    <p className={`mb-4 break-all ${isDarkMode ? 'text-blue-300' : 'text-blue-600'}`}>
                        {truncateHash(txHash)}
                    </p>
                    <div className="flex justify-between">
                        <button
                            onClick={closeModal}
                            className={`px-4 py-2 rounded ${isDarkMode ? 'bg-gray-600 text-white' : 'bg-gray-200 text-black'}`}
                        >
                            Close
                        </button>
                        <a
                            href={`https://etherscan.io/tx/${txHash}`}
                            target="_blank"
                            rel="noopener noreferrer"
                            className={`px-4 py-2 rounded ${isDarkMode ? 'bg-blue-600 text-white' : 'bg-blue-500 text-white'}`}
                        >
                            View on Etherscan
                        </a>
                    </div>
                </div>
            </div>
        );
    };

    return (
        <div className={`p-8 ${bgColor} ${textColor} bg-gradient-to-br flex items-center justify-center`}>
            <div className={`${cardBgColor} rounded-3xl shadow-xl overflow-hidden w-full max-w-xl relative`}>
                <div className="p-8">
                    <h2 className="text-3xl font-bold mb-6 text-center">Claim Your Rewards</h2>
                    <ErrorMessage errorMessage={errorMessage} />

                    <form onSubmit={handleSearch} className="mb-6">
                        <div className="relative">
                            <input
                                type="text"
                                value={recipient}
                                onChange={(e) => !isConnected && setRecipient(e.target.value)}
                                className={`w-full ${inputBgColor} rounded-lg px-4 py-3 pr-12 focus:outline-none focus:ring-2 focus:ring-blue-400 ${textColor} placeholder-gray-400 text-sm`}
                                placeholder="Enter account address"
                                disabled={isConnected || isLoading}
                                style={{ minHeight: '48px' }}
                            />
                            <button
                                type="submit"
                                className={`absolute right-3 top-1/2 transform -translate-y-1/2 ${textColor} hover:text-blue-300 transition-colors`}
                                disabled={isLoading}
                            >
                                <Search size={20} />
                            </button>
                        </div>
                    </form>

                    <div className="space-y-4 mb-8">
                        <RewardItem label="Eligible Rewards" value={eligibleRewards} isDarkMode={isDarkMode} />
                        <RewardItem label="Claimed Rewards" value={claimedRewards} isDarkMode={isDarkMode} />
                        <RewardItem label="Rewards to Claim" value={rewardsToClaim} highlight={true} isDarkMode={isDarkMode} />
                    </div>

                    {isConnected ? (
                        <div className="space-y-4">
                            <button
                                onClick={handleClaim}
                                className={`w-full ${buttonBgColor} text-white font-bold py-3 px-4 rounded-lg transition duration-300 flex items-center justify-center`}
                                disabled={rewardsToClaim <= 0 || isLoading}
                            >
                                {rewardsToClaim > 0 ? (
                                    <>
                                        <CheckCircle size={20} className="mr-2" />
                                        Claim Rewards
                                    </>
                                ) : (
                                    'No Rewards to Claim'
                                )}
                            </button>
                            <button
                                onClick={() => open({ view: 'Account' })}
                                className="w-full bg-red-500 hover:bg-red-600 text-white font-bold py-3 px-4 rounded-lg transition duration-300 flex items-center justify-center"
                                disabled={isLoading}
                            >
                                <LogOut size={20} className="mr-2" />
                                Disconnect Wallet
                            </button>
                        </div>
                    ) : (
                        <button
                            onClick={() => open()}
                            className={`w-full ${buttonBgColor}  text-white font-bold py-3 px-4 rounded-lg transition duration-300 flex items-center justify-center`}
                            disabled={isLoading}
                        >
                            <Wallet size={20} className="mr-2" />
                            Connect to Claim
                        </button>
                    )}
                </div>

                {isConnected && address && (
                    <div className={`${isDarkMode ? 'bg-gray-700' : 'bg-gray-200'} p-4 text-center text-sm ${textColor}`}>
                        Connected: {address}
                    </div>
                )}

                {(isLoading || isClaiming) && (
                    <div className={`absolute inset-0 ${isDarkMode ? 'bg-black bg-opacity-50' : 'bg-white bg-opacity-70'} flex items-center justify-center z-10`}>
                        <div className={`${isDarkMode ? 'bg-white bg-opacity-20' : 'bg-gray-200 bg-opacity-90'} rounded-lg p-4 flex flex-col items-center`}>
                            <Loader className={`animate-spin ${isDarkMode ? 'text-blue-400' : 'text-blue-600'}`} size={40} />
                            <p className={`mt-2 ${isDarkMode ? 'text-white' : 'text-gray-800'}`}>
                                {isClaiming ? 'Claiming reward...' : 'Fetching reward...'}
                            </p>
                        </div>
                    </div>
                )}
            </div>
            {showTxModal && (
                <TransactionModal
                    txHash={txHash}
                    safeTransactionStatus={safeTransactionStatus}
                    isDarkMode={isDarkMode}
                    closeModal={() => setShowTxModal(false)}
                />
            )}
        </div>
    );
};

const RewardItem = ({ label, value, highlight = false, isDarkMode }) => {
    const bgColor = isDarkMode
        ? (highlight ? 'bg-blue-900 bg-opacity-50' : 'bg-gray-700')
        : (highlight ? 'bg-blue-100' : 'bg-gray-200');
    const textColor = isDarkMode ? 'text-white' : 'text-gray-800';
    const borderColor = highlight ? 'border-blue-400' : 'border-transparent';

    const numericValue = Number(value) || 0;

    const formattedValue = numericValue >= 1e18
        ? (numericValue / 1e18).toFixed(3)
        : numericValue.toFixed(3);

    return (
        <div className={`${bgColor} rounded-lg p-4 ${highlight ? `border-2 ${borderColor}` : ''}`}>
            <label className={`block ${isDarkMode ? 'text-gray-300' : 'text-gray-600'} mb-1 text-sm`}>{label}</label>
            <div className={`text-xl font-semibold ${textColor} flex justify-between items-center`}>
                <span>{formattedValue}</span>
                <span className={`${isDarkMode ? 'text-blue-300' : 'text-blue-500'} text-sm`}>SSV</span>
            </div>
        </div>
    );
};

export default Claim;