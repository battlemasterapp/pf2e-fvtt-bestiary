#!/bin/bash

# Clone the git repository
git clone --depth 1 git@github.com:foundryvtt/pf2e.git ./pf2e

# Copy the bestiary packs
mkdir -p ./tmp
cp -r ./pf2e/packs/*-bestiary* ./tmp/