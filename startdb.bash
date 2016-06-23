#!/bin/bash

sudo service mongodb stop
numactl --interleave=all mongod --dbpath=/home/dorival/data/db
