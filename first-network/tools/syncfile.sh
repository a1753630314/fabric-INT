for i in `cat ./iplist.txt | awk '{print $1}'`;do  ssh $i 'rm -rf /root/fabric-test';rsync -HavP /root/fabric-test root@$i:/root/;done
