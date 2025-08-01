#!/bin/bash

# Fix absolute paths in JavaScript files
echo "Fixing absolute paths in JavaScript files..."

# Replace /static/ with static/ in all JS files
find web/static/js -name "*.js" -type f -exec sed -i '' 's|/static/|static/|g' {} +

echo "Done! All absolute paths have been converted to relative paths."