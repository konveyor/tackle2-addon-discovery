#!/bin/bash

host="${HUB_API:-localhost:8080}"
application="${1}"
curl -k -X POST ${host}/tasks \
  -H 'Content-Type:application/x-yaml' \
  -H 'Accept:application/x-yaml' \
 -d \
"
---
state: Ready
addon: language-discovery
application:
  id: "${application}"
data: {}
"
