# Copyright (c) 2026 José María Micoli
# Licensed under {'license_type': 'BSL', 'change_date': '2033-02-17', 'convert_to': 'Apache-2.0'}
# You may:
# ✔ Study
# ✔ Modify
# ✔ Use for internal security testing
# You may NOT:
# ✘ Offer as a commercial service
# ✘ Sell derived competing products
/*
Copyright (c) 2026 José María Micoli
Licensed under the Business Source License 1.1
Change Date: 2033-02-17
Change License: Apache-2.0

You may:
✔ Study
✔ Modify
✔ Use for internal security testing

You may NOT:
✘ Offer as a commercial service
✘ Sell derived competing products
*/

#!/bin/bash

# Diagnostic test script
# Captures all logs while running VaporTrace scan

echo "=== Starting VaporTrace Diagnostic Test ==="
echo "Target: https://httpbin.org"
echo ""
echo "Starting VaporTrace with logging..."
echo ""

# Run VaporTrace and capture output
timeout 60 ./VaporTrace 2>&1 &
VAPID=$!

sleep 3

# Send commands
echo "target https://httpbin.org" | nc localhost 3000 2>/dev/null &
sleep 2

echo "swagger" | nc localhost 3000 2>/dev/null &
sleep 5

echo "miner" | nc localhost 3000 2>/dev/null &
sleep 5

echo "scrape" | nc localhost 3000 2>/dev/null &
sleep 5

echo "map" | nc localhost 3000 2>/dev/null &
sleep 5

echo "quit" | nc localhost 3000 2>/dev/null &

# Wait for VaporTrace to finish
wait $VAPID 2>/dev/null || true

echo ""
echo "=== Test Complete ==="
