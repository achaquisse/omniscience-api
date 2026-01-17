#!/bin/bash

echo "Running integration tests..."
go test -v -run Test .

exit_code=$?

if [ $exit_code -eq 0 ]; then
    echo "✅ All tests passed!"
else
    echo "❌ Some tests failed"
fi

exit $exit_code
