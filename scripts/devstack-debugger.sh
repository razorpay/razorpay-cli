#!/bin/sh
# =============================================================================
# Devstack Remote Debugger Script
# =============================================================================
# This script attaches Delve debugger to the running application process.
# It is executed by DevSpace after file sync to enable remote debugging.
#
# Usage: ./devstack-debugger.sh <port>
#   - port: The port on which Delve will listen (e.g., 2345)
# =============================================================================

echo "Remote Debugger"

echo "Getting the PID of Delve process if it is already running"
DLV_PIDS=$(ps -ef | grep dlv | grep -v grep | awk '{print $2}')
if [ -n "$DLV_PIDS" ]; then
  echo "Killing existing Delve processes: $DLV_PIDS"
  for pid in $DLV_PIDS; do
    kill -9 "$pid" 2>/dev/null || true
  done
fi

echo "[Debugger] Waiting for process to start..."
while ! curl -s -w '%{http_code}\n' 'http://localhost:8081/ready' 2>/dev/null | grep -q 200; do
    echo "[Debugger] Application not ready yet, retrying in 5s..."
    sleep 5
done

echo "Application is ready, getting PID of the running application"
PID=$(ps -ef | grep "/tmp/app" | grep -v "CompileDaemon" | grep -v "grep" | awk '{print $2}' | head -1)

if [ -n "$PID" ]; then
  echo "Attaching debugger to process: $PID on port :$1"
  dlv attach "$PID" --listen=":$1" --headless --api-version=2 --accept-multiclient --check-go-version=false
else
  echo "ERROR: Could not find application PID"
  exit 1
fi

