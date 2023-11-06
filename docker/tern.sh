#!/bin/ash

echo "executing db migration via tern..."

if [ "z$NOMAD_ALLOC_INDEX" == "z0" ]; then
	tern status
	tern migrate
fi