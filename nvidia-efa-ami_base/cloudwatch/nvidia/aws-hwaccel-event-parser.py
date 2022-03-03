import subprocess
import platform

plat=platform.platform()
if 'Ubuntu' in plat:
    f = subprocess.Popen(['tail','-F','-n','1000','/var/log/syslog'],\
        stdout=subprocess.PIPE,stderr=subprocess.PIPE)
else:
    f = subprocess.Popen(['tail','-F','-n','1000','/var/log/messages'],\
        stdout=subprocess.PIPE,stderr=subprocess.PIPE)

while True:
    line = f.stdout.readline()
    if "NVRM:" in str(line):
       file = open('/var/log/gpuevent.log','a')
       print(line.decode('utf8', errors='strict').strip())
       file.write(line.decode('utf8', errors='strict').strip())
       file.write('\n')
       file.close()
