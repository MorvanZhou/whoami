#!/bin/bash
source .venv/bin/activate

if [ "$1" = "dev" ]; then
  uv run uvicorn app.main:app --reload --port 8002
else
  uv run uvicorn app.main:app --port 8002
fi