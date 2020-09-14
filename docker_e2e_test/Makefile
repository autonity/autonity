.PHONY: docker-e2e-tests clean install-prerequisites-linux install-prerequisites-mac

# Runs the docker e2e tests.
#
# By default each new line in a rule's commands is run in a separate shell. To
# run multiple commands in the same shell the commands need to be separated
# with ';' or '&&'. So any python commands that depend on a venv needs to be
# run in the same shell in which the venv is activated.
docker-e2e-tests: ./venv/bin/activate
	. venv/bin/activate && sudo python3 test_via_docker.py ..

# Run a long term testcase on docker e2e test framework.
long-term-e2e-tests: ./venv/bin/activate
	. venv/bin/activate && sudo python3 test_via_docker.py .. -l True

# Creates venv if not present and then installs dependencies. Executes only if
# requirements_docker_test.txt has been updated since the last time it ran.
#
# Make targets that start with a dot are considered hidden and are not provided
# as autocomplete options.
./venv/bin/activate: requirements_docker_test.txt
	test -d venv || python3 -m venv venv
	. venv/bin/activate && sudo pip3 install -r requirements_docker_test.txt
	touch ./venv/bin/activate

clean:
	rm -rf venv

# Installs the prerequisites for a Linux user
install-prerequisites-linux:
	sudo apt-get install -y python3 python3-venv python3-pip golang-docker-credential-helpers

# Installs the prerequisites for a mac user
install-prerequisites-mac:
	brew install python3
