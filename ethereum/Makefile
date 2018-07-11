DEPLOY_TARGET ?= dev

deps:
	yarn

run:
	yarn ganache:verbose

docker-build:
	docker build -t horizon-games/ganache .

deploy:
	ansible-playbook -i ../../SkyWeaver/ansible.hosts -v playbook.yml -e "host=$(DEPLOY_TARGET)"
