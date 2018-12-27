for i in `cat ./iplist.txt | awk '{print $1 }'` ; do ssh root@$i 'cd /root/fabric-test/first-network/; sh createdir.sh' ; done
