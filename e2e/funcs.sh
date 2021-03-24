#!/bin/bash


# stdout outputs meesage to standard output
stdout() {
    echo ">>> $1"
}

# fails echos message to standard error and bails out
fail() {
    echo "!!! $1"
    exit 5
}
