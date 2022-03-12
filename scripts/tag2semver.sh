#!/bin/bash
# -*- Mode: Shell-script -*-
#
# tag2semver.sh --- Convert a git tag to semver info.
#
# Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
#
# Author:     Paul Ward <asmodai@gmail.com>
# Maintainer: Paul Ward <asmodai@gmail.com>
# Created:    09 Jan 2022 10:31:16
#
# {{{ License:
#
# This program is free software: you can redistribute it
# and/or modify it under the terms of the GNU General Public
# License as published by the Free Software Foundation,
# either version 3 of the License, or (at your option) any
# later version.
#
# This program is distributed in the hope that it will be
# useful, but WITHOUT ANY  WARRANTY; without even the implied
# warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR
# PURPOSE.  See the GNU General Public License for more
# details.
#
# You should have received a copy of the GNU General Public
# License along with this program.  If not, see
# <http://www.gnu.org/licenses/>.
#
# }}}
# {{{ Commentary:
#
# }}}

revs=$(git rev-list --tags --max-count=1 2>/dev/null)
tags=$(git describe --tags ${revs} 2>/dev/null)
commit=$(git rev-parse --short HEAD 2>/dev/null)

# Check if VERSION is defined
if [[ ${VERSION:-""} == v*.*.* ]]
then
    tags=${VERSION}
fi

# If not, check if we've been given an argument.
if [ $# -eq 1 ]
then
    tags=$1
fi

# If not, test the current git tag.
if [[ ${tags:-""} == v*.*.* ]]
then
    tags=$(echo ${tags} | tr -d [:alpha:])
    semver=(${tags//\./ })
else
    # Well, we give up right here.
    echo "-X main.GitVers=\"0:${commit:-\"<local>\"}\""
    exit 0
fi

# Compute number.
maj=$((semver[0]*10000000))
min=$((semver[1]*10000))
vers=$((maj+min+semver[2]))

# Print it.
echo "-X main.GitVers=\"${vers}:${commit:-\"<local>\"}\""

# tag2semver.sh ends here.
