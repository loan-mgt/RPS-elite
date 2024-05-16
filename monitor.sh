#!/bin/bash
 
 
 
# Check if a file argument is provided
if [ $# -ne 1 ]; then
    echo "Usage: $0 <file_to_monitor>"
    exit 1
fi
 
 
 
# File to monitor passed as argument
FILE="$1"
 
 
 
# Check if the file exists
if [ ! -f "$FILE" ]; then
    echo "Error: File '$FILE' not found."
    exit 1
fi
 
 
 
# Get the initial modification time of the file
LTIME=$(stat -c %Z "$FILE")
 
 
 
# Function to execute the command when file is modified
execute_command() {
    go run main.go
    # Add your command to run here
}
 
 
 
while true; do
    # Get the current modification time of the file
    ATIME=$(stat -c %Z "$FILE")
 
 
 
    # Check if the modification time has changed
    if [[ "$ATIME" != "$LTIME" ]]; then
        execute_command
        # Update the last modification time
        LTIME=$ATIME
    fi
 
 
 
    # Sleep for a short duration before checking again
    sleep 5
done
 