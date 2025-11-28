"use client";

import { useEffect, useState, useMemo } from 'react';

// Define the Event type to match the backend
const RedactionDashboard = () => {
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchEvents = async () => {
    try {
      // Connect to the backend API
      const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:3001';
      const response = await fetch(`${API_URL}/api/v1/events/list`);

      if (!response.ok) {
        throw new Error(`Failed to fetch events: ${response.statusText}`);
      }

      const data = await response.json();

      // Sort events by timestamp in descending order
      setEvents(data.sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()));
      setError(null);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    // Initial fetch
    fetchEvents();

    // Fetch events every 5 seconds
    const intervalId = setInterval(fetchEvents, 5000);

    // Cleanup interval on component unmount
    return () => clearInterval(intervalId);
  }, []);

  const totalRedactions = useMemo(() => events.length, [events]);
  const redactionsByType = useMemo(() => {
    return events.reduce((acc, event) => {
      acc[event.redacted_type] = (acc[event.redacted_type] || 0) + 1;
      return acc;
    }, {});
  }, [events]);

  return (
    <div className="bg-gray-900 text-gray-200 min-h-screen font-sans">
      <div className="container mx-auto p-4 md:p-8">
        <header className="flex justify-between items-center mb-8 border-b border-cyan-900 pb-4">
          <div>
            <h1 className="text-3xl font-bold text-cyan-400 tracking-wider">CISO COMMAND CENTER</h1>
            <p className="text-gray-500">Cilium-Shield: AI-DLP Firewall</p>
          </div>
          <div className="text-right">
             <div className="text-lg text-gray-400">Total Redactions</div>
             <div className="text-4xl font-bold text-red-500">{totalRedactions}</div>
          </div>
        </header>

        {error && (
          <div className="bg-red-900 border border-red-700 text-red-200 p-4 rounded-md mb-6 animate-pulse">
            <p><span className="font-bold">ERROR:</span> Could not connect to the Shield Observer. {error}</p>
          </div>
        )}

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          {Object.entries(redactionsByType).map(([type, count]) => (
            <div key={type} className="bg-gray-800 p-4 rounded-lg shadow-lg border border-gray-700">
              <h3 className="text-md font-semibold text-gray-400">{type}</h3>
              <p className="text-3xl font-bold text-yellow-400">{count}</p>
            </div>
          ))}
        </div>

        <div className="bg-gray-800 shadow-2xl rounded-lg overflow-hidden border border-gray-700">
          <div className="px-6 py-4">
            <h2 className="text-xl font-semibold text-gray-300">Live Redaction Feed</h2>
          </div>
          <div className="overflow-x-auto">
            <table className="min-w-full">
              <thead className="bg-gray-900">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-cyan-400 uppercase tracking-wider">Timestamp</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-cyan-400 uppercase tracking-wider">Source Pod IP</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-cyan-400 uppercase tracking-wider">Data Type</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-cyan-400 uppercase tracking-wider">Destination</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-700">
                {loading ? (
                   <tr>
                    <td colSpan={4} className="text-center py-10 text-gray-500">Loading initial events...</td>
                  </tr>
                ) : events.length > 0 ? (
                  events.map((event, index) => (
                    <tr key={index} className="hover:bg-gray-700/50 transition-colors duration-150">
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-400">{new Date(event.timestamp).toLocaleString()}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-cyan-300 font-mono">{event.source_pod_ip}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm">
                        <span className="px-3 py-1 inline-flex text-xs leading-5 font-bold rounded-full bg-red-600 text-white">
                          {event.redacted_type}
                        </span>
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-400 font-mono">{event.destination_url}</td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan={4} className="text-center py-10 text-gray-500">
                      No redaction events detected. Monitoring live...
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
        <footer className="text-center text-gray-600 mt-8 text-xs">
          <p>Cilium-Shield | eBPF-Powered L7 Security</p>
        </footer>
      </div>
    </div>
  );
};

export default RedactionDashboard;
