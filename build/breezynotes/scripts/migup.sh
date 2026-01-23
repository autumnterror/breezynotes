#!/bin/bash
CONFIG_PATH=./configs/migrator.yaml
export CONFIG_PATH
./bin/migrator --type up
unset CONFIG_PATH
