#!/bin/bash
bambam -o . -p protocol sum.go
capnpc -ogo schema.capnp

