out=`curl -s -X POST -F file=@scripts.sh localhost:8080/isula/task`
echo "task id created is: " $out 

curl localhost:8080/isula/task/$out
echo "\n"

out2=`curl -s -X POST -F file=@scripts.sh localhost:8080/isula/task?output=isula-0.0.1.iso`
echo "task id created is: " $out2

curl localhost:8080/isula/task/$out2
echo "\n"

curl -s -X DELETE localhost:8080/isula/task/$out2
echo "\n"
curl localhost:8080/isula/task/$out2
echo "\n"
curl -s -X DELETE localhost:8080/isula/task/$out2
echo "\n"
