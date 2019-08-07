go-lint:
	golint -min_confidence=0.01 -set_exit_status=1

go-build:
	go build -o ../bin/azd-kubernetes-manager .

go-run:
	../bin/azd-kubernetes-manager --token=${AZURE_DEVOPS_TOKEN} --url=${AZURE_DEVOPS_URL} --log-level=Trace

go-test:
	go test github.com/ggmaresca/azd-kubernetes-manager/pkg/tests

docker-build:
	docker build -t azd-kubernetes-manager:dev .

docker-run:
	docker run -it --rm --name=azd-kubernetes-manager -v ${HOME}/.kube:/home/azd-kubernetes-manager/.kube:ro --network=host azd-kubernetes-manager:dev --name=azd-kubernetes-manager --namespace=default --token=${AZURE_DEVOPS_TOKEN} --url=${AZURE_DEVOPS_URL} --log-level=Trace

docker-push:
	sh docker-push.sh

docker-clean:
	docker rmi azd-kubernetes-manager:dev

helm-lint:
	helm lint charts/azd-kubernetes-manager

helm-template:
	helm template charts/azd-kubernetes-manager --set azd.url=https://dev.azure.com/test,azd.token=abc123def456ghi789jkl,pdb.enabled=true,serviceMonitor.enabled=true

helm-install:
	helm upgrade --debug --install azd-kubernetes-manager charts/azd-kubernetes-manager --set azd.url=${AZURE_DEVOPS_URL},azd.existingSecret=azd-agent,azd.existingSecretKey=azd-token,logLevel=trace

helm-package:
	helm package charts/azd-kubernetes-manager -d charts && \
	helm repo index --merge charts/index.yaml charts
