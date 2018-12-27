for i in `grep "/data" *.yaml|awk '{print $3}'|awk -F':' '{print $1}'`
do
	rm -rf $i
	mkdir -p $i
        chmod 777 $i
done
