#!/bin/bash

echo "====================================="
echo "Starting Cilium-Shield Dev Environment"
echo "====================================="
echo ""

# Start backend in background
echo "Starting backend on port 3001..."
cd backend
npm install
npm start &
BACKEND_PID=$!
echo "Backend PID: $BACKEND_PID"
cd ..

# Wait for backend to start
echo "Waiting for backend to start..."
sleep 3

# Start frontend in background
echo "Starting frontend on port 3000..."
cd frontend
npm install
npm start &
FRONTEND_PID=$!
echo "Frontend PID: $FRONTEND_PID"
cd ..

echo ""
echo "====================================="
echo "Dev environment is running!"
echo "====================================="
echo ""
echo "Backend:  http://localhost:3001"
echo "Frontend: http://localhost:3000"
echo ""
echo "Press Ctrl+C to stop all services"
echo ""

# Wait for Ctrl+C
trap "echo ''; echo 'Stopping services...'; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit" INT

# Keep script running
wait
