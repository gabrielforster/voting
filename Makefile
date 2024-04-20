build:
	@echo "Building..."

docker-services: docker-auth docker-poll
	@echo "Built Docker image of auth and poll services"

docker-auth: tidy-auth
	@echo "Building Docker of auth service"
	@docker build -t auth -f ./auth/Dockerfile .

docker-poll: tidy-poll
	@echo "Building Docker of poll service"
	@docker build -t poll -f ./poll/Dockerfile .

tidy-auth:
	@echo "Tidying auth service"
	@cd auth && go mod tidy

tidy-poll:
	@echo "Tidying pol service"
	@cd poll && go mod tidy
