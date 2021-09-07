#!/bin/bash

/usr/bin/python3 /opt/aws/aws-hwaccel-error-parser.py &
/usr/bin/python3 /opt/aws/accel-to-cw.py /opt/aws/binaries/nvidia-exporter >> /dev/null 2>&1 &
