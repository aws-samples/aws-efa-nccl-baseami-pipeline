#!/bin/bash -xe

nvidia-smi mig -cgi $1 -C
