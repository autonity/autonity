.PHONY: docker-e2e-tests clean install-prerequisites-linux install-prerequisites-mac

# Runs the docker e2e tests.
#
# By default each new line in a rule's commands is run in a separate shell. To
# run multiple commands in the same shell the commands need to be separated
# with ';' or '&&'. So any python commands that depend on a venv needs to be
# run in the same shell in which the venv is activated.
docker-e2e-tests: ./venv/bin/activate
	. venv/bin/activate && sudo python3 test_via_docker.py .. -id=$(ID) -hash=$(HASH)

# Creates venv if not present and then installs dependencies. Executes only if
# requirements_docker_test.txt has been updated since the last time it ran.
#
# Make targets that start with a dot are considered hidden and are not provided
# as autocomplete options.
./venv/bin/activate: requirements_docker_test.txt
	test -d venv || python3 -m venv venv
	python3 --version
	pip3 --version
	. venv/bin/activate && sudo pip3 install -r requirements_docker_test.txt
	touch ./venv/bin/activate

clean:
	rm -rf venv

# Installs the prerequisites for a Linux user
install-prerequisites-linux:
	gpg --recv-keys "B26995E310250568"
	wget https://www.python.org/ftp/python/3.8.2/Python-3.8.2.tgz
	wget https://www.python.org/ftp/python/3.8.2/Python-3.8.2.tgz.asc
	gpgv --keyring ~/.gnupg/pubring.kbx Python-3.8.2.tgz.asc Python-3.8.2.tgz
	tar -xf Python-3.8.2.tgz
	cd Python-3.8.2 && ./configure --with-ensurepip=install && make && sudo make install
	python3 --version
	sudo apt install golang-docker-credential-helpers

# Installs the prerequisites for a mac user
install-prerequisites-mac:
	brew install python3
