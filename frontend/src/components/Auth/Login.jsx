import React, { useState } from 'react';
import { Ticket, LogIn, CheckCircle, AlertCircle } from 'lucide-react';
import MessageAlert from '../Shared/MessageAlert';

const Login = ({ setView, setCurrentUser, loadEvents, showMessage }) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);

  const handleLogin = async () => {
    setLoading(true);
    try {
      const response = await fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
      });

      if (response.ok) {
        const user = await response.json();
        setCurrentUser(user);
        localStorage.setItem('currentUser', JSON.stringify(user));
        setView('events');
        loadEvents();
        showMessage('success', `Welcome back!`);
      } else {
        showMessage('error', 'Invalid email or password');
      }
    } catch (error) {
      showMessage('error', 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-600 via-pink-500 to-red-500 flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl shadow-2xl p-8 w-full max-w-md">
        <div className="text-center mb-8">
          <Ticket className="w-16 h-16 text-purple-600 mx-auto mb-4" />
          <h1 className="text-3xl font-bold text-gray-800">Ticket Queen</h1>
          <p className="text-gray-600 mt-2">Login to purchase tickets</p>

        </div>
        <div className="space-y-4">
          <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} className="w-full px-4 py-2 border rounded-lg" placeholder="Email" />
          <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} className="w-full px-4 py-2 border rounded-lg" placeholder="Password" />
          <button onClick={handleLogin} disabled={loading} className="w-full bg-purple-600 text-white py-3 rounded-lg flex items-center justify-center gap-2">
            <LogIn className="w-5 h-5" /> {loading ? 'Logging in...' : 'Login'}
          </button>
        </div>
        <p className="mt-6 text-center text-gray-600">
          New here? <button onClick={() => setView('signup')} className="text-purple-600 font-bold">Sign up</button>
        </p>
      </div>
    </div>
  );
};

export default Login;