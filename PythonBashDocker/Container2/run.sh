#!/bin/bash
export PYTHONPATH="$PYTHONPATH::$HOME/local/lib/python"
dude run
dude sum
Rscript graphs.R
