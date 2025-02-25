#!/bin/bash

# Check if both arguments are provided
if [ -z "$1" ]; then
    echo "Usage: $0 <path_to_executable> <time_file_name>"
    exit 1
fi

# Define the path to your executable and the time file name from arguments
EXECUTABLE="$1"
TIME_FILE="$2"

# Check if the executable exists
if [ ! -x "$EXECUTABLE" ]; then
    echo "Error: $EXECUTABLE not found or not executable."
    exit 1
fi

# Clear the time file if it exists, to start fresh
> "$TIME_FILE"

# Loop to run the executable 100 times
for i in $(seq 1 100)
do
    # Print the current iteration number to track progress
    echo "Running iteration $i of 100..."
    # Capture the time taken for each run and append it to the file
    { time "$EXECUTABLE"; } 2>> "$TIME_FILE"
done

echo "Execution times have been saved to $TIME_FILE."
