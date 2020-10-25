rm ../bin/gymClassBooker
env GOOS=linux GOARCH=arm GOARM=5 go build -o ../bin ./...
cd ../bin
mv ./main ./gymClassBooker
#./gymClassBooker