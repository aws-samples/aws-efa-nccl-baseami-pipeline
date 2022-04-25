import subprocess
import os

if __name__ =='__main__':
    subprocess.call(['env'])
    subprocess.call(['/workspace/nccl-tests/build/all_reduce_perf', \
            '-b', \
            '8', \
            '-e', \
            '1G', \
            '-f', \
            '2', \
            '-t', \
            '1', \
            '-g', \
            '1', \
            '-c',\
            '1', \
            '-n',\
            '100'])
