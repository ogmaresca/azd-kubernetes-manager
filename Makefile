go-lint:
	golint -min_confidence=0.01 -set_exit_status=1

go-build:
	go build -o ../bin/azd-kubernetes-manager .

go-run:
	../bin/azd-kubernetes-manager --token=${AZURE_DEVOPS_TOKEN} --url=${AZURE_DEVOPS_URL} --config-file example-config.yaml --log=debug --username=a --password=b

go-test:
	go clean -testcache && go test -cover ./...

go-test-azuredevops:
	go clean -testcache && go test -cover -test.v ./pkg/azuredevops -args --log=debug

go-test-config:
	go clean -testcache && go test -cover -test.v ./pkg/config -args --log=debug

go-test-kubernetes:
	go clean -testcache && go test -cover -test.v ./pkg/kubernetes -args --log=debug

go-test-processors:
	go clean -testcache && go test -cover -test.v ./pkg/processors -args --log=debug

docker-build:
	docker build -t azd-kubernetes-manager:dev .

docker-run:
	docker run -it --rm --name=azd-kubernetes-manager -v ${HOME}/.kube:/home/azd-kubernetes-manager/.kube:ro -v `pwd`/example-config.yaml://home/azd-kubernetes-manager/configuration.yaml:ro --network=host azd-kubernetes-manager:dev --token=${AZURE_DEVOPS_TOKEN} --url=${AZURE_DEVOPS_URL} --log=debug --config-file=/home/azd-kubernetes-manager/configuration.yaml

docker-push:
	sh docker-push.sh

docker-clean:
	docker rmi azd-kubernetes-manager:dev

helm-lint:
	helm lint charts/azd-kubernetes-manager

helm-template:
	helm template charts/azd-kubernetes-manager --values=example-helm-values.yaml

helm-install:
	helm upgrade --debug --install azd-kubernetes-manager charts/azd-kubernetes-manager --values=example-helm-values.yaml --set image.repository=azd-kubernetes-manager,image.tag=dev

helm-package:
	helm package charts/azd-kubernetes-manager -d charts && \
	helm repo index --merge charts/index.yaml charts
