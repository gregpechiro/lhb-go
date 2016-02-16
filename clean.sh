#!/bin/bash

# a script to remove old files and folders
# this should only be run on the server where the project lives right before
# redeployment. It will remove all files that are contained in the new .tar

rm -rf static/ templates/ auth.go main.go webmasterRoutes.go lhb-go lhb-go.tar
