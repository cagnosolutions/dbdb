#!/bin/bash

### NOTES: ###
# Added 1<<20 entries (a little over 1 million) at the end of the tests.
# Memory stats indicated 150mb allocated total (ram), however the system...
#  ... monitor had me hanging around 50-60mb ram (GC'd)
# 10 total GC's for 1<<20 entries, and only 150mb RAMnot too bad
#
# Took about 3 seconds on a Intel I-5 (4 Cores)

go test -memprofile mem.out
