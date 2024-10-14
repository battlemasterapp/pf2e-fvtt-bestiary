#!/bin/bash

# Clone the git repository
git clone --depth 1 git@github.com:foundryvtt/pf2e.git ./pf2e

# Copy the bestiary packs
mkdir -p ./tmp
cp -r ./pf2e/packs/*-bestiary* ./tmp/

mkdir -p ./bestiaries

# Get all directories in the tmp folder
for dir in ./tmp/*; do
  # Get the name of the directory
  dir_name=$(basename $dir)
  echo "Processing $dir_name"

  # Concat all internal json files into a single file
  jq -s 'flatten' ./tmp/$dir_name/*.json > ./bestiaries/$dir_name.json
done

# Create an index file
echo "Creating index file"
touch ./bestiaries/index.json
echo "{" > ./bestiaries/index.json
for file in ./bestiaries/*.json; do
  # Add the name and path to the index file
  file_name=$(basename $file .json)
  # Remove the `-bestiary` suffix
  key=${file_name%-bestiary}

  # if the key is `index`, skip it
  if [ "$key" == "index" ]; then
    continue
  fi
  echo "\"$key\": \"${file_name}.json\"," >> ./bestiaries/index.json
done
# Clear the last comma
sed -i '$ s/.$//' ./bestiaries/index.json
echo "}" >> ./bestiaries/index.json

# Get the spells
echo "Processing spells"
mkdir -p ./spells
jq -s 'flatten' ./pf2e/packs/spells/*.json > ./spells/spells.json