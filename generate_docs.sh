#!/bin/bash

# Generate Go HTML documentation
# This script creates HTML documentation files in the godoc/ directory

echo "Generating Go HTML documentation..."

# Create godoc directory if it doesn't exist
mkdir -p godoc

# Generate HTML documentation for each package
echo "Generating main package documentation..."
godoc -url="/pkg/myutilsport/" > godoc/main.html

echo "Generating utils package documentation..."
godoc -url="/pkg/myutilsport/utils/" > godoc/utils.html

echo "Generating cmd package documentation..."
godoc -url="/pkg/myutilsport/cmd/listcodes/" > godoc/cmd.html

echo "Documentation generation complete!"
echo ""
echo "Generated files:"
echo "  HTML documentation: godoc/"
echo "    - godoc/main.html"
echo "    - godoc/utils.html"
echo "    - godoc/cmd.html"
echo ""
echo "To view HTML documentation, open the files in your browser:"
echo "  open godoc/main.html    # macOS"
echo "  xdg-open godoc/main.html # Linux"
