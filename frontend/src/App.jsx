import React, { useState, useEffect } from 'react';
import Login from './components/Auth/Login';
import Signup from './components/Auth/Signup';
import UserDashboard from './components/Dashboard/UserDashboard';
import EventList from './components/Store/EventList';
import ShoppingCart from './components/Store/ShoppingCart';
import MessageAlert from './components/Shared/MessageAlert';
import { Ticket, ShoppingCart as CartIcon, User, LogOut } from 'lucide-react';

const App = () => {
  const [currentUser, setCurrentUser] = useState(null);
  const [view, setView] = useState('login');
  const [events, setEvents] = useState([]);
  const [cart, setCart] = useState([]);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState(null);

  useEffect(() => {
    const savedUser = localStorage.getItem('currentUser');
    if (savedUser) {
      const user = JSON.parse(savedUser);
      setCurrentUser(user);
      setView('events');
      loadEvents();
    }
  }, []);

  const showMessage = (type, text) => {
    setMessage({ type, text });
    setTimeout(() => setMessage(null), 20000);
  };

  const loadEvents = async () => {
    try {
      const response = await fetch('/api/events');
      const data = await response.json();
      setEvents(data);
    } catch (error) {
      showMessage('error', 'Failed to load events');
    }
  };

  const handleLogout = () => {
    setCurrentUser(null);
    localStorage.removeItem('currentUser');
    setView('login');
    setCart([]);
    showMessage('success', 'Logged out successfully');
  };

  // Auth Views
  if (view === 'login') return (
    <>
      <MessageAlert message={message} />
      <Login setView={setView} setCurrentUser={setCurrentUser} loadEvents={loadEvents} showMessage={showMessage} />
    </>
  );

  if (view === 'signup') return (
    <>
      <MessageAlert message={message} />
      <Signup setView={setView} setCurrentUser={setCurrentUser} loadEvents={loadEvents} showMessage={showMessage} />
    </>
  );
  // Dashboard View
  if (view === 'dashboard') return <UserDashboard currentUser={currentUser} setView={setView} handleLogout={handleLogout} showMessage={showMessage} />;

  // Main Storefront View
  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-600 via-pink-500 to-red-500">
      <div className="container mx-auto px-4 py-8">
        <header className="bg-white rounded-lg shadow-2xl p-6 mb-6 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <Ticket className="w-10 h-10 text-purple-600" />
            <h1 className="text-4xl font-bold text-gray-800">Ticket Queen</h1>
          </div>
          <div className="flex items-center gap-3">
            <button onClick={() => setView('dashboard')} className="px-4 py-2 bg-purple-100 text-purple-600 rounded-lg hover:bg-purple-200 transition flex items-center gap-2">
              <User className="w-5 h-5" /> {currentUser?.username}
            </button>
            <div className="flex items-center gap-2 bg-purple-100 px-4 py-2 rounded-full">
              <CartIcon className="w-5 h-5 text-purple-600" />
              <span className="font-semibold text-purple-600">{cart.length} items</span>
            </div>
            <button onClick={handleLogout} className="px-4 py-2 bg-red-100 text-red-600 rounded-lg hover:bg-red-200 transition">
              <LogOut className="w-5 h-5" />
            </button>
          </div>
        </header>

        <MessageAlert message={message} />

        <div className="grid md:grid-cols-3 gap-6">
          <div className="md:col-span-2">
            <EventList events={events} setCart={setCart} cart={cart} showMessage={showMessage} />
          </div>
          <div className="md:col-span-1">
            <ShoppingCart cart={cart} setCart={setCart} currentUser={currentUser} loadEvents={loadEvents} showMessage={showMessage} />
          </div>
        </div>
      </div>
    </div>
  );
};

export default App;