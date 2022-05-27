docker-img:
	docker build -t fgh151/db-server ./

docker-push:
	docker push fgh151/db-server

docker: docker-img docker-push