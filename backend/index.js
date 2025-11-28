const express = require('express');
const cors = require('cors');
const app = express();
const port = process.env.PORT || 3001;

// Middleware
app.use(cors());
app.use(express.json());

// In-memory event store
let events = [];

// Health check endpoint
app.get('/', (req, res) => {
  res.json({
    status: 'ok',
    message: 'Cilium-Shield Observer Backend',
    totalEvents: events.length
  });
});

// Endpoint to receive redaction events from Wasm filter or log agents
app.post('/api/v1/events', (req, res) => {
  const { source_pod_ip, destination_url, redacted_type } = req.body;

  if (!source_pod_ip || !destination_url || !redacted_type) {
    return res.status(400).json({
      error: 'Missing required fields: source_pod_ip, destination_url, redacted_type'
    });
  }

  const event = {
    timestamp: new Date().toISOString(),
    source_pod_ip,
    destination_url,
    redacted_type
  };

  events.push(event);
  console.log(`[EVENT] Received redaction event:`, event);

  res.status(201).json({
    message: 'Event received successfully',
    event
  });
});

// Endpoint to list all events (for dashboard)
app.get('/api/v1/events/list', (req, res) => {
  // Return events in reverse chronological order
  const sortedEvents = [...events].sort((a, b) =>
    new Date(b.timestamp) - new Date(a.timestamp)
  );

  res.json(sortedEvents);
});

// Endpoint to get event statistics
app.get('/api/v1/events/stats', (req, res) => {
  const stats = {
    total: events.length,
    byType: {},
    byDestination: {}
  };

  events.forEach(event => {
    // Count by redaction type
    stats.byType[event.redacted_type] = (stats.byType[event.redacted_type] || 0) + 1;

    // Count by destination
    stats.byDestination[event.destination_url] = (stats.byDestination[event.destination_url] || 0) + 1;
  });

  res.json(stats);
});

// Clear all events (useful for testing)
app.delete('/api/v1/events', (req, res) => {
  const count = events.length;
  events = [];
  res.json({
    message: `Cleared ${count} events`,
    count
  });
});

app.listen(port, () => {
  console.log(`[SERVER] Cilium-Shield Observer Backend listening at http://localhost:${port}`);
  console.log(`[SERVER] API endpoints:`);
  console.log(`  - GET  /                     - Health check`);
  console.log(`  - POST /api/v1/events        - Submit redaction event`);
  console.log(`  - GET  /api/v1/events/list   - List all events`);
  console.log(`  - GET  /api/v1/events/stats  - Get event statistics`);
  console.log(`  - DELETE /api/v1/events      - Clear all events`);
});
