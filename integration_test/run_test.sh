#!/bin/bash
# Integration test for gobpftool
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
CONTAINER_NAME="gobpftool-test-$$"
IMAGE_NAME="gobpftool-test"

cleanup() {
    echo "Cleaning up..."
    docker rm -f "$CONTAINER_NAME" 2>/dev/null || true
}
trap cleanup EXIT

echo "=== Building gobpftool ==="
cd "$PROJECT_DIR"
make build

echo "=== Building test container ==="
cp gobpftool integration_test/
docker build -t "$IMAGE_NAME" integration_test/

echo "=== Starting container with BPF capabilities ==="
docker run -d --name "$CONTAINER_NAME" \
    --privileged \
    "$IMAGE_NAME" sleep 300

echo "=== Loading eBPF program with bpftool ==="
docker exec "$CONTAINER_NAME" bpftool prog load test_prog.o /sys/fs/bpf/test_prog type xdp || {
    echo "Note: Could not load with bpftool, trying ip link..."
    docker exec "$CONTAINER_NAME" ip link set dev lo xdp obj test_prog.o sec xdp 2>/dev/null || true
}

echo "=== Testing: prog show ==="
OUTPUT=$(docker exec "$CONTAINER_NAME" ./gobpftool prog show 2>&1) || true
echo "$OUTPUT"
if echo "$OUTPUT" | grep -q "test_prog"; then
    echo "✓ prog show: Found test_prog"
else
    echo "✗ prog show: test_prog not found"
    exit 1
fi

echo "=== Testing: prog show (JSON) ==="
OUTPUT=$(docker exec "$CONTAINER_NAME" ./gobpftool -j prog show 2>&1) || true
echo "$OUTPUT"
if echo "$OUTPUT" | grep -q '"name":"test_prog"'; then
    echo "✓ prog show -j: Found test_prog in JSON"
else
    echo "✗ prog show -j: test_prog not found in JSON"
    exit 1
fi

echo "=== Testing: map show ==="
OUTPUT=$(docker exec "$CONTAINER_NAME" ./gobpftool map show 2>&1) || true
echo "$OUTPUT"
if echo "$OUTPUT" | grep -q "test_map"; then
    echo "✓ map show: Found test_map"
else
    echo "✗ map show: test_map not found"
    exit 1
fi

echo "=== Testing: map show (JSON) ==="
OUTPUT=$(docker exec "$CONTAINER_NAME" ./gobpftool -j map show 2>&1) || true
echo "$OUTPUT"
if echo "$OUTPUT" | grep -q '"name":"test_map"'; then
    echo "✓ map show -j: Found test_map in JSON"
else
    echo "✗ map show -j: test_map not found in JSON"
    exit 1
fi

echo "=== Testing: version ==="
OUTPUT=$(docker exec "$CONTAINER_NAME" ./gobpftool --version)
echo "$OUTPUT"
if echo "$OUTPUT" | grep -q "gobpftool version"; then
    echo "✓ version: Output valid"
else
    echo "✗ version: Invalid output"
    exit 1
fi

echo ""
echo "=== All integration tests passed! ==="
