import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Route, NavLink, Routes, useLocation, useNavigate } from 'react-router-dom';
import { Sun, Moon, Award, Grid, Users, Hexagon, Shield, Menu, X, Bell } from 'lucide-react';
import { FaGithub } from 'react-icons/fa';
import { createWeb3Modal, defaultConfig, useWeb3ModalTheme } from '@web3modal/ethers/react'

import Dashboard from './components/Dashboard';
import Operators from './components/Operators';
import Clusters from './components/Clusters';
import Validators from './components/Validators';
import AccountDetails from './components/AccountDetails';
import ClusterDetails from './components/ClusterDetails';
import Claim from './components/Claim';
import ClusterMonitor from './components/ClusterMonitor';

const projectId = process.env.REACT_APP_WALLET_CONNECT_PROJECT_ID
const mainnet = {
  chainId: 1,
  name: 'Ethereum',
  currency: 'ETH',
  explorerUrl: 'https://etherscan.io',
  rpcUrl: 'https://cloudflare-eth.com'
}
const holesky = {
  chainId: 17000,
  name: 'Ethereum',
  currency: 'ETH',
  explorerUrl: 'https://holesky.etherscan.io',
  rpcUrl: 'https://ethereum-holesky.publicnode.com'
}

const metadata = {
  name: 'Monitor SSV APP',
  description: 'MonitorSSV - SSV Network Monitoring Tool',
  url: 'https://monitorssv.xyz',
  icons: ['']
}

createWeb3Modal({
  ethersConfig: defaultConfig({ metadata }),
  chains: [mainnet, holesky],
  projectId
})

function App() {
  const [isDarkMode, setIsDarkMode] = useState(true);
  const [network, setNetwork] = useState('mainnet');
  const currentYear = new Date().getFullYear();
  const { setThemeMode } = useWeb3ModalTheme()

  const toggleDarkMode = () => {
    setIsDarkMode(!isDarkMode);
  };

  useEffect(() => {
    setThemeMode(isDarkMode ? 'dark' : 'light')
  }, [isDarkMode, setThemeMode]);

  return (
    <Router>
      <AppContent
        isDarkMode={isDarkMode}
        toggleDarkMode={toggleDarkMode}
        network={network}
        setNetwork={setNetwork}
        currentYear={currentYear}
      />
    </Router>
  );
}

function AppContent({ isDarkMode, toggleDarkMode, network, setNetwork, currentYear }) {
  const location = useLocation();
  const isDashboardPage = location.pathname === '/';
  const navigate = useNavigate();
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const handleClaimClick = () => {
    navigate('/claim');
  };

  // Add a class to the body to always show scrollbar
  useEffect(() => {
    document.body.classList.add('overflow-y-scroll');
    return () => {
      document.body.classList.remove('overflow-y-scroll');
    };
  }, []);

  return (
    <div className={`min-h-screen flex flex-col font-sans ${isDarkMode ? 'dark bg-gray-900 text-gray-100' : 'bg-gray-50 text-gray-900'}`}>
      <header className={`sticky top-0 z-50 backdrop-blur-md bg-opacity-80 ${isDarkMode ? 'bg-gray-900' : 'bg-white'} shadow-lg`}>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-4">
            <div className="flex items-center">
              <h1 className="text-2xl font-bold">
                Monitor<span className={`${isDarkMode ? 'text-indigo-400' : 'text-indigo-600'}`}>SSV</span>
              </h1>
            </div>
            <div className="hidden md:flex items-center space-x-4">
              <NavItems isDarkMode={isDarkMode} />
            </div>
            <div className="flex items-center space-x-4">
              <select
                value={network}
                onChange={(e) => setNetwork(e.target.value)}
                className={`rounded-md px-3 py-1 text-sm focus:outline-none focus:ring-2 transition-all duration-300 ${isDarkMode
                  ? 'bg-gray-800 text-gray-200 border-gray-700 focus:ring-indigo-500'
                  : 'bg-gray-100 text-gray-800 border-gray-300 focus:ring-indigo-500'
                  }`}
              >
                <option value="Mainnet">Mainnet</option>
                {/*<option value="Holesky">Holesky</option>*/}
              </select>
              <button
                onClick={toggleDarkMode}
                className={`p-2 rounded-full transition-all duration-300 focus:outline-none focus:ring-2 ${isDarkMode
                  ? 'bg-gray-800 text-yellow-400 hover:bg-gray-700 focus:ring-yellow-500'
                  : 'bg-gray-200 text-indigo-600 hover:bg-gray-300 focus:ring-indigo-500'
                  }`}
                aria-label={isDarkMode ? "Switch to light mode" : "Switch to dark mode"}
              >
                {isDarkMode ? <Sun className="w-4 h-4" /> : <Moon className="w-4 h-4" />}
              </button>
              <button
                onClick={() => setIsMenuOpen(!isMenuOpen)}
                className={`md:hidden p-2 rounded-full transition-all duration-300 focus:outline-none focus:ring-2 ${isDarkMode
                  ? 'bg-gray-800 text-gray-200 hover:bg-gray-700 focus:ring-indigo-500'
                  : 'bg-gray-200 text-gray-800 hover:bg-gray-300 focus:ring-indigo-500'
                  }`}
              >
                {isMenuOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
              </button>
            </div>
          </div>
        </div>
      </header>

      <nav className={`md:hidden ${isMenuOpen ? 'block' : 'hidden'} ${isDarkMode ? 'bg-gray-800' : 'bg-white'}`}>
        <div className="px-2 pt-2 pb-3 space-y-1">
          <NavItems isDarkMode={isDarkMode} mobile />
        </div>
      </nav>

      <main className="flex-grow max-w-7xl mx-auto px-4 py-8 w-full">
        {isDashboardPage && (
          <div className={`mb-8 p-6 rounded-xl shadow-lg flex flex-col md:flex-row items-center justify-between transform hover:scale-102 transition-all duration-300 ${isDarkMode
            ? 'bg-gradient-to-r from-blue-900 via-indigo-900 to-purple-900 text-white'
            : 'bg-gradient-to-r from-blue-400 to-indigo-500 text-white'
            }`}>
            <div className="flex items-center mb-4 md:mb-0">
              <Award className="w-10 h-10 mr-4 text-yellow-300" />
              <span className="text-xl font-semibold">Claim your rewards from Incentivized Mainnet Program now!</span>
            </div>
            <button
              onClick={handleClaimClick}
              className={`px-6 py-2 rounded-full font-bold transition-colors shadow-md hover:shadow-lg focus:outline-none focus:ring-2 focus:ring-offset-2 ${isDarkMode
                ? 'bg-indigo-500 text-white hover:bg-indigo-600 focus:ring-indigo-400'
                : 'bg-white text-indigo-600 hover:bg-indigo-50 focus:ring-indigo-500'
                }`}
            >
              CLAIM
            </button>
          </div>
        )}
        <div className={`rounded-xl shadow-lg p-6 transition-all duration-300 ${isDarkMode ? 'bg-gray-800 border border-gray-700' : 'bg-white border border-gray-200'
          }`}>
          <Routes>
            <Route path="/" element={<Dashboard isDarkMode={isDarkMode} network={network} />} />
            <Route path="/operators" element={<Operators isDarkMode={isDarkMode} network={network} />} />
            <Route path="/clusters" element={<Clusters isDarkMode={isDarkMode} network={network} />} />
            <Route path="/cluster/:id" element={<ClusterDetails isDarkMode={isDarkMode} network={network} />} />
            <Route path="/account/:address" element={<AccountDetails isDarkMode={isDarkMode} network={network} />} />
            <Route path="/validators" element={<Validators isDarkMode={isDarkMode} network={network} />} />
            <Route path="/claim" element={<Claim isDarkMode={isDarkMode} />} />
            <Route path="/monitor" element={<ClusterMonitor isDarkMode={isDarkMode} network={network} />} />
          </Routes>
        </div>
      </main>

      <footer className={`mt-auto py-6 transition-all duration-300 ${isDarkMode ? 'bg-gray-800 border-t border-gray-700' : 'bg-gray-100 border-t border-gray-200'
        }`}>
        <div className="max-w-7xl mx-auto px-4 flex flex-col md:flex-row justify-between items-center">
          <div className="text-sm font-medium mb-4 md:mb-0">
            Project granted by SSV Network
          </div>
          <div className="flex items-center space-x-4">
            <a
              href="https://github.com/monitorssv"
              target="_blank"
              rel="noopener noreferrer"
              className={`transition-colors ${isDarkMode
                ? 'text-gray-400 hover:text-indigo-400'
                : 'text-gray-600 hover:text-indigo-500'
                }`}
            >
              <FaGithub size={24} />
            </a>
            <span className="text-sm">©{currentYear} Monitorssv</span>
          </div>
        </div>
      </footer>
    </div>
  );
}

function NavItems({ isDarkMode, mobile = false }) {
  const navItems = [
    { to: "/", icon: <Grid className="w-5 h-5" />, label: "Dashboard" },
    { to: "/operators", icon: <Users className="w-5 h-5" />, label: "Operators" },
    { to: "/clusters", icon: <Hexagon className="w-5 h-5" />, label: "Clusters" },
    { to: "/validators", icon: <Shield className="w-5 h-5" />, label: "Validators" },
    { to: "/claim", icon: <Award className="w-5 h-5" />, label: "Claim" },
    { to: "/monitor", icon: <Bell className="w-5 h-5" />, label: "Monitor" },
  ];

  return navItems.map((item) => (
    <NavLink
      key={item.to}
      to={item.to}
      className={({ isActive }) =>
        `flex items-center px-3 py-2 rounded-md text-sm font-medium transition-all duration-300 ${mobile ? 'flex' : ''
        } ${isActive
          ? isDarkMode
            ? 'text-white bg-gray-900'
            : 'text-indigo-600 bg-indigo-50'
          : isDarkMode
            ? 'text-gray-300 hover:bg-gray-700 hover:text-white'
            : 'text-gray-600 hover:bg-gray-100 hover:text-indigo-600'
        }`
      }
    >
      {item.icon}
      <span className={mobile ? 'ml-3' : 'hidden md:block ml-2'}>{item.label}</span>
    </NavLink>
  ));
}

export default App;