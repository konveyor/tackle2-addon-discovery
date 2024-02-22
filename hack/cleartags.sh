#!/bin/bash

host="${HUB_API:-localhost:8080}"
application="${1}"
curl -XPUT -k ${host}/applications/${application}/tags\?source\=discovery -d "[]"
