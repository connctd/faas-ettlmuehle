.PHONY: build push deploy

build:
	faas-cli build -f ettlmuehle-fn.yml --build-arg GO111MODULE=on

push:	
	faas-cli push -f ettlmuehle-fn.yml

install:
	kubectl apply -f secrets.yaml
	kubectl apply -f cronjob.yaml
	faas-cli deploy -f ettlmuehle-fn.yml

uninstall:
	kubectl delete -f secrets.yaml
	kubectl delete -f cronjob.yaml
	faas-cli remove -f ettlmuehle-fn.yml
