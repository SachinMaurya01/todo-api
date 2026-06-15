#!/bin/bash

# 1. Load environment variables from .env file, ignoring comments and empty lines
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

COMMAND=$1
NAME=$2

case "$COMMAND" in
    "up")
        migrate -path migrations -database "$DATABASE_URL" up
        ;;
    "down")
        # Default to 1 rollback if no count parameter is provided
        COUNT=${NAME:-1}
        
        # Prompt user for confirmation
        read -p "Rolling back $COUNT migration(s). Continue? [y/N]: " CONFIRM
        
        if [ "$CONFIRM" = "y" ] || [ "$CONFIRM" = "Y" ]; then
            migrate -path migrations -database "$DATABASE_URL" down "$COUNT"
        else
            echo "Operation cancelled."
        fi
        ;;
    "create")
        if [ -z "$NAME" ]; then
            echo "Error: Please provide a name for the migration. Example: ./migrate.sh create init_schema"
            exit 1
        fi
        migrate create -ext sql -dir migrations -seq "$NAME"
        ;;
    "force")
        if [ -z "$NAME" ]; then
            echo "Error: Please provide a version number to force. Example: ./migrate.sh force 1"
            exit 1
        fi
        migrate -path migrations -database "$DATABASE_URL" force "$NAME"
        ;;
    *)
        echo "Usage: $0 {up|down|create|force} [argument]"
        exit 1
        ;;
esac
