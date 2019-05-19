#!/bin/sh

SCRIPTPATH=$(cd "$(dirname "$0")"; pwd)
"$SCRIPTPATH/notification-ms" -importPath notification-ms -srcPath "$SCRIPTPATH/src" -runMode dev
