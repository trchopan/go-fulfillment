#!/bin/bash

CompileDaemon -exclude-dir=.git -exclude=".#*" -recursive=false -command="./go-rest-api"
