#!/bin/bash

# Usage check
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <pod_name_1> <pod_name_2>"
    exit 1
fi

POD1=$1
POD2=$2

# Get Pod IPs
echo "Fetching IP for $POD1..."
IP1=$(kubectl get pod "$POD1" -o jsonpath='{.status.podIP}')
echo "Fetching IP for $POD2..."
IP2=$(kubectl get pod "$POD2" -o jsonpath='{.status.podIP}')

if [ -z "$IP1" ]; then
    echo "Error: Could not find IP for $POD1"
    exit 1
fi

if [ -z "$IP2" ]; then
    echo "Error: Could not find IP for $POD2"
    exit 1
fi

echo "Pod 1: $POD1 IP: $IP1"
echo "Pod 2: $POD2 IP: $IP2"

# Function to check connectivity
check_conn() {
    local src_pod=$1
    local dst_ip=$2
    
    echo "Checking $src_pod -> $dst_ip ..."
    # We use curl with a timeout of 5 seconds.
    # Exit code 7 means connection refused (host up, no service).
    # Exit code 28 means timeout (likely blocked or host down).
    kubectl exec "$src_pod" -- curl -s --connect-timeout 5 http://"$dst_ip" > /dev/null 2>&1
    local status=$?
    
    if [ $status -eq 0 ] || [ $status -eq 7 ]; then
        return 0
    else
        return 1
    fi
}

echo "----------------------------------------"
if check_conn "$POD1" "$IP2"; then
    RESULT_1="WORKS"
else
    RESULT_1="FAILS"
fi

echo "----------------------------------------"
if check_conn "$POD2" "$IP1"; then
    RESULT_2="WORKS"
else
    RESULT_2="FAILS"
fi

echo "----------------------------------------"
echo "Summary:"
echo "$POD1 -> $POD2: $RESULT_1"
echo "$POD2 -> $POD1: $RESULT_2"
