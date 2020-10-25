rm ../bin/gymClassBooker
env GOOS=linux GOARCH=arm64 go build -o ../bin ./...
cd ../bin
mv ./main ./gymClassBooker
#./gymClassBooker