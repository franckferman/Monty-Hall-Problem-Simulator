#!/usr/bin/env bash

# ============================
# 🏗️ Monty Hall Problem Simulator - Build Script
# ============================

set -e  # Stop if error occurs
set -u  # Treat unset variables as error
set -o pipefail  # Handle errors in pipelines

# ----------- Config -----------
PROJECT_NAME="Monty-Hall-Problem-Simulator"
MAIN_FILE="main.go"
OUTPUT_DIR="bin"

# Clean old binaries
echo "🧹 Cleaning previous builds..."
rm -rf "$OUTPUT_DIR"
mkdir -p "$OUTPUT_DIR"

# ----------- Build Flags -----------
BUILD_FLAGS="-trimpath -ldflags='-s -w'"
CGO_ENABLED=0

# ----------- Build Targets -----------
declare -A TARGETS
TARGETS=(
  ["Linux"]="GOOS=linux GOARCH=amd64"
  ["Windows"]="GOOS=windows GOARCH=amd64"
  ["MacOS"]="GOOS=darwin GOARCH=amd64"
  ["MacOS-ARM64"]="GOOS=darwin GOARCH=arm64"
)

# ----------- Build Process -----------
echo "🚀 Starting build process..."

for target in "${!TARGETS[@]}"; do
  echo "🔨 Building for $target..."

  # Output file name
  output_file="$OUTPUT_DIR/$PROJECT_NAME-$target"
  [ "$target" == "Windows" ] && output_file="$output_file.exe"

  # Command
  build_cmd="CGO_ENABLED=$CGO_ENABLED ${TARGETS[$target]} go build $BUILD_FLAGS -o $output_file $MAIN_FILE"

  # Execute
  echo "⚙️  $build_cmd"
  eval $build_cmd

  # Check
  if [ -f "$output_file" ]; then
    echo "✅ Built $output_file"
  else
    echo "❌ Failed to build $target"
    exit 1
  fi
done

echo "🎉 All binaries built successfully and available in '$OUTPUT_DIR/'"

