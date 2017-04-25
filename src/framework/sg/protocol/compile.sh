#!/bin/bash
bambam -o . -p protocol message.go
capnpc -ogo schema.capnp

