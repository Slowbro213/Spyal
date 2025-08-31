#!/bin/sh
if [ "$ENV" = "production" ]; then
    echo "Starting production server..."
    /app/server
else
    echo "Starting development server with Air..."
    air
fi
