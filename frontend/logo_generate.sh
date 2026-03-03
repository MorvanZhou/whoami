#!/bin/bash
# Generate favicon and icon variants from logo.png
# Uses macOS built-in sips (no extra dependencies needed)

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SRC="${SCRIPT_DIR}/public/logo.png"
OUT="${SCRIPT_DIR}/public"

if [ ! -f "$SRC" ]; then
  echo "Error: ${SRC} not found"
  exit 1
fi

echo "Generating icons from ${SRC}..."

generate() {
  local size=$1
  local name=$2
  cp "$SRC" "${OUT}/${name}"
  sips -z "$size" "$size" "${OUT}/${name}" --out "${OUT}/${name}" > /dev/null 2>&1
}

generate 180 "apple-touch-icon.png"
generate 16  "favicon-16x16.png"
generate 32  "favicon-32x32.png"
generate 192 "icon-192.png"
generate 512 "icon-512.png"

echo "Done:"
ls -lh "${OUT}/apple-touch-icon.png" \
       "${OUT}/favicon-16x16.png" \
       "${OUT}/favicon-32x32.png" \
       "${OUT}/icon-192.png" \
       "${OUT}/icon-512.png"
