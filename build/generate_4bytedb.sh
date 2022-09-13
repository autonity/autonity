#!/bin/bash
SOLC_BIN="$1"
echo "{" > signer/fourbyte/autonity_4byte.json
$SOLC_BIN --hashes autonity/solidity/contracts/Autonity.sol 2>/dev/null | grep -E --color=never "^[0-9a-f]{8}: .*" | sed -E 's/^([0-9a-f]{8}): (.*)$/"\1": "\2",/' >> signer/fourbyte/autonity_4byte.json
LAST_LINE=$(cat signer/fourbyte/autonity_4byte.json | wc -l)
sed -i "$LAST_LINE s/,$//" signer/fourbyte/autonity_4byte.json
echo "}" >> signer/fourbyte/autonity_4byte.json
cat signer/fourbyte/4byte.json signer/fourbyte/autonity_4byte.json | jq -s add > signer/fourbyte/4byte_merged.json
mv signer/fourbyte/4byte_merged.json signer/fourbyte/4byte.json
rm signer/fourbyte/autonity_4byte.json
