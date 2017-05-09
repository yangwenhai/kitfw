#!/bin/bash
# bambam and capnpc can only run under the mac os
# capnp will need capnpc-go to generate golang code
bambam -o . -p protocol message.go
capnp compile -ogo ./schema.capnp

