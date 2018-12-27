for i in `cat ./ipnew.txt | awk '{print $1}'`;do  ssh $i 'rm -rf /root/fabric-test';rsync -HavP /root/fabric-test root@$i:/root/;done
