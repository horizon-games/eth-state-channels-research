DEPLOY_TARGET ?= dev

deps:
	yarn

run:
	yarn ganache:verbose

docker-build:
	docker build -t horizon-games/ganache .

docker-run:
	docker run \
		-p 127.0.0.1:8545 \
		horizon-games/ganache \
		yarn ganache:verbose

deploy:
	ansible-playbook -i ../../SkyWeaver/ansible.hosts -v playbook.yml -e "host=$(DEPLOY_TARGET)"
