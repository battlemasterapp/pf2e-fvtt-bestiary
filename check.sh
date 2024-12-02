#/bin/bash

# Checks if all json files are valid since the parser can sometimes delete control characters

for file in ./bestiaries/*.json; do
  echo "Checking $file"
  jq . $file > /dev/null
  if [ $? -ne 0 ]; then
    echo "Error in $file"
    exit 1
  fi
done