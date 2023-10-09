ms_build:
	docker build -t mysql1:1.0 .

ms_run:
	docker run -d --name=mysqlserver1 -p 3306:3306 mysql1:1.0

ms_stop:
	docker container stop mysqlserver1

ms_clean:
	docker container rm mysqlserver1
	docker rmi $(shell docker images 'mysql1' -a -q)

run:
	go run main.go