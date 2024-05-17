#!/bin/bash

# Check if the correct number of arguments are passed
if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <directory_to_monitor> <file_extension> <command_to_execute>"
    exit 1
fi

# Directory to monitor
DIR_TO_MONITOR="$1"

# File extension to monitor
FILE_EXTENSION="$2"

# Command to execute on change
COMMAND_TO_EXECUTE="$3"

# Use inotifywait to monitor the directory for create, modify, and delete events
inotifywait -m -e create -e modify -e delete --format '%w%f' "$DIR_TO_MONITOR" | while read FILE
do

        $COMMAND_TO_EXECUTE
        echo "Change detected in file: $FILE"

done
