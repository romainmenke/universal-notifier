#!/bin/bash

echo "$@"

"$@" $WERCKER_STEP_ROOT/universal-notifier
