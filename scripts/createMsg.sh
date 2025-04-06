#!/bin/bash

# Define usage function
usage() {
  echo "Usage: $0 [-r | -c | -b]"
  echo "  -r  Code review mode only"
  echo "  -c  Code comment mode only"
  echo "  -b  Both modes (default)"
  exit 1
}

# Parse command-line options
while getopts "rcb" opt; do
  case "$opt" in
    r) MODE="review";;
    c) MODE="comment";;
    b) MODE="both";;
    *) usage;;
  esac
done

# Default to both modes if no flags are provided
MODE=${MODE:-"both"}

# Run the main program with the selected mode
export OPEN_ROUTER_KEY="sk-or-v1-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" # Replace with your actual API key
go run main.go --mode "$MODE"
