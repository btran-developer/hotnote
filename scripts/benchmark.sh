#!/bin/bash

# HotNote Benchmark Runner
# Runs performance benchmarks and generates reports

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
BENCHMARK_DIR="$PROJECT_DIR/.benchmarks"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
REPORT_FILE="$BENCHMARK_DIR/report_$TIMESTAMP.txt"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "HotNote Benchmark Runner"
echo "========================"
echo ""

# Create benchmark directory
mkdir -p "$BENCHMARK_DIR"

# Parse arguments for profiling
PROFILE=""
if [[ "${1:-}" == "--cpuprofile" ]]; then
    PROFILE="-cpuprofile=cpu.prof"
elif [[ "${1:-}" == "--memprofile" ]]; then
    PROFILE="-memprofile=mem.prof"
fi

# Run benchmarks
echo "Running benchmarks..."
echo "Results will be saved to: $REPORT_FILE"
echo ""

# Run benchmarks with memory stats and output to file
cd "$PROJECT_DIR"
go test -bench=. -benchmem $PROFILE -run=^$ ./cmd/ 2>&1 | tee "$REPORT_FILE"

echo ""
echo "========================"
echo "Benchmark Report"
echo "========================"
echo ""

# Analyze results
if command -v awk >/dev/null 2>&1; then
    # Parse benchmark results
    echo "Performance Summary:"
    echo "-------------------"

    # Check for benchmarks that exceed 100ms (100000000 ns)
    SLOW_BENCHMARKS=$(grep -E "^Benchmark" "$REPORT_FILE" | awk '{
        # Extract ns/op value (4th field)
        nsop = $3
        # Remove commas from numbers (Go formats large numbers with commas)
        gsub(/,/, "", nsop)
        # Convert to ms
        ms = nsop / 1000000
        if (ms > 100) {
            printf "  ⚠️  %s: %.2f ms (exceeds 100ms threshold)\n", $1, ms
        } else {
            printf "  ✅ %s: %.2f ms\n", $1, ms
        }
    }')

    echo -e "$SLOW_BENCHMARKS"
    echo ""

    # Calculate average memory usage
    echo "Memory Usage:"
    echo "-------------"
    grep -E "^Benchmark" "$REPORT_FILE" | awk '{
        printf "  %s: %s B/op, %s allocs/op\n", $1, $5, $7
    }'
else
    echo "awk not available, showing raw results:"
    cat "$REPORT_FILE"
fi

echo ""
echo "========================"
echo "Full report saved to: $REPORT_FILE"
echo "========================"
