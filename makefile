include .env

migration:
	goose -dir internal/db/migrations create $(name) sql

migrate:
	goose -dir 'internal/db/migrations' postgres ${DATABASE_URL} up

migrate-up-one:
	goose -dir 'internal/db/migrations' postgres ${DATABASE_URL} up-by-one

migratedbs:
	make migrate

rollback:
	goose -dir 'internal/db/migrations' postgres ${DATABASE_URL} down

api:
	go run cmd/api/*.go

docker-build:
	docker build -t savannah-pos:latest .

docker-run:
	docker run --rm -p 8080:8080 --env-file .env savannah-pos:latest

k8s-apply:
	kubectl apply -f k8s/ -n savannah-pos

k8s-delete:
	kubectl delete -f k8s/ -n savannah-pos

k8s-get-all:
	kubectl get all -n savannah-pos

k8s-describe-pod:
	kubectl describe pod $(pod) -n savannah-pos

k8s-logs:
	kubectl logs -f $(pod) -n savannah-pos

k8s-port-forward:
	kubectl port-forward svc/savannah-pos-service 8080:80 -n savannah-pos

k8s-rollout-restart:
	kubectl rollout restart deployment/savannah-pos -n savannah-pos

k8s-rollout-status:
	kubectl rollout status deployment/savannah-pos -n savannah-pos

k8s-delete-pods:
	kubectl delete pods -l app=savannah-pos -n savannah-pos

k8s-migrate:
	kubectl apply -f k8s/migration-job.yaml -n savannah-pos

k8s-migrate-logs:
	kubectl logs -f job/db-migration -n savannah-pos
