#!/usr/bin/env bash

out=${1:-cmd/ping/kodata/config.js}

echo "writing to $out"

echo "// NOTE: Please do not depend on the values in this file, as they" > $out
echo "// may change without warning. Instead, use the process in"       >> $out
echo "// DEVELOPMENT.md to create your own regional Cloud Run services" >> $out
echo "// and use those instead."                                        >> $out
echo "let _URLS = {"                                                    >> $out
echo '  "global": "https://global.gcping.com",'                         >> $out

terraform output services | tail -n+2 | sed -e 's/ =/:/g' | sed -e 's/"$/",/g' | sed -e 's/}/};/g' >> $out

cat $out
