#!/bin/bash

set -e

echo -e "Installing plugin..."

#go install
cf uninstall-plugin multiapps

PLUGIN_NAME_WIN_32=multiapps-plugin.win32
PLUGIN_NAME_WIN_64=multiapps-plugin.win64
PLUGIN_NAME_LINUX_32=multiapps-plugin.linux32
PLUGIN_NAME_LINUX_64=multiapps-plugin.linux64
PLUGIN_NAME_OSX=multiapps-plugin.osx

# get OS version with 'uname -s'
plugin_name=$PLUGIN_NAME_WIN_64

#cf install-plugin $GOPATH/build/$plugin_name -f
cf install-plugin /c/Users/I356120/go/src/github.com/cloudfoundry-incubator/multiapps-cli-plugin/build/$plugin_name -f