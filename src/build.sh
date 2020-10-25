rm ../bin/gymClassBooker
go build -o ../bin ./...
cd ../bin
mv ./main ./gymClassBooker
#./gymClassBooker