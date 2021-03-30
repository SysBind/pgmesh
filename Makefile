pgmesh:
	CGO_ENABLED=0 go build

e2e: e2e_clean docker_image
	@cd e2e/moodle && ./run.sh

docker_image: pgmesh Dockerfile
	docker build . -t pgmesh

e2e_clean:
	@cd e2e/moodle && ./cleanup.sh
