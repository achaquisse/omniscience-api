#!/bin/bash

echo "Running integration tests..."
cd src && go test -v ./...

exit_code=$?

if [ $exit_code -eq 0 ]; then
    echo "✅ All tests passed!"
else
    echo "❌ Some tests failed"
fi

exit $exit_code
