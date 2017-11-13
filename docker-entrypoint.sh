#!/bin/bash
set -eo pipefail

# AWS_PS_VERSION - (optional) if you know a specific version of your key(s)
# AWS_PS_PATH    - Search a dir path in AWS
# AWS_PS_KEY     - Search for a particular entry in AWS

# Set version param attribute.
export _version=""
if [ ! -z "$AWS_PS_VERSION" ]
then
  _version="--version ${AWS_PS_VERSION}"
fi

# Search by path.
if [ ! -z "$AWS_PS_PATH" ]
then
	eval $(./aws-ps-client getpath "${AWS_PS_PATH}" $_version)
fi

# Search by key.
if [ ! -z "$AWS_PS_KEY" ]
then
	eval $(./aws-ps-client get "${AWS_PS_KEY}" $_version)
fi

# Dump the variables from memory.

echo "========= Environment Varuables ========"
env | sort
echo "========================================="

# Don't die for 1 hour.
sleep(3600)
