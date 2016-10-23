#!/bin/bash -ex
tugboat ssh -n monitoring -o "-L 3306:localhost:3306 -L 9102:localhost:9102"
