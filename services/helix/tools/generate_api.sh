#!/bin/bash

MSG_DIR="$(dirname "$0")/../message"

protoc --go_out="$MSG_DIR" --go_opt=paths=source_relative --proto_path="$MSG_DIR" "$MSG_DIR"/*.proto
