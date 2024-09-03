#!/bin/sh

cd /app && chmod 775 server && chmod +x server && exec /app/server
