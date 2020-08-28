import os
import stat

def prepare_exp(SSHHost, SSHPort, REMOTEROOT, optpt):
    f = open("config", 'w')
    f.write("Host benchmark\n")
    f.write("   Hostname %s\n" % SSHHost)
    f.write("   Port %d\n" % SSHPort)
    f.write("   User ubuntu\n")
    f.close()
    

    f = open("run-experiment.sh", 'w')
    f.write("#!/bin/bash\n")
    f.write("set -x\n\n")
    
    f.write("ssh -i /home/ubuntu/id_rsa -F config benchmark -o StrictHostKeyChecking=no -p 22 \"nohup memcached -u ubuntu -P memcached.pid > memcached.out 2> memcached.err &\"\n")
    
    f.write("RESULT=`ssh -i /home/ubuntu/id_rsa -F config benchmark  -o StrictHostKeyChecking=no -p 22 \"pidof memcached\"`\n")

    f.write("sleep 5\n")

    f.write("if [ -z \"$RESULT\"]; then echo \"memcached process not running\"; CODE=1; else CODE=0; fi\n")
        
    f.write("mcperf -N %d -R %d -n %d --server=server --port=11211 2> stats.log\n\n" % (optpt["noRequests"]*10,optpt["noRequests"],optpt["concurrency"]))
    f.write("REQPERSEC=$(cat stats.log | grep \"Response rate\" | awk '{print $3}')\n")
    f.write("LATENCY=$(cat stats.log | grep \"Response time \[ms\]: avg\" | awk '{print $5}')\n")    

    f.write("ssh -i /home/ubuntu/id_rsa -F config benchmark -o StrictHostKeyChecking=no \"kill -9 $RESULT\"\n")

    f.write("echo \"requests latency\" > stats.csv\n")
    f.write("echo \"$REQPERSEC $LATENCY\" >> stats.csv\n")
    

    f.write("if [ $(wc -l <stats.csv) -le 1 ]; then CODE=1; fi\n\n")
    
    f.write("exit $CODE\n")

    f.close()
    
    os.chmod("run-experiment.sh", stat.S_IRWXU)
