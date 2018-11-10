#!/bin/bash

ADDR=0.0.0.0 \
JWT_SECRET="./jwtRS256.key" \
buffalo dev
