.PHONY: serve clean install-prerequisites-linux install-prerequisites-mac

serve: build
	# By default each new line in a rule's commands is run in a separate shell.
	# To run multiple commands in the same shell the commands need to be
	# separated with a semicolon. So any python commands that depend on a venv
	# need to be run in the same shell in which the venv is activated.
	. venv/bin/activate; cd build/html && python3 -m http.server 4567


# Builds the docs, executes if python modules or any source file have been
# updated.
build: ./venv/bin/activate $(shell find source)
	make clean
	. venv/bin/activate; sphinx-build -M html source build

# Creates venv if not present and then installs dependencies. Executes if
# requirements.txt has been updtaed.
#
# Make targets that start with a dot are considered hidden and are not provided
# as autocomplete options.
./venv/bin/activate: requirements.txt
	test -d venv || python3 -m venv venv
	. venv/bin/activate; pip3 install -r requirements.txt
	touch ./venv/bin/activate

# Deletes the build
clean:
	rm -rf build

# Installs the prerequisites for a linux user
install-prerequisites-linux:
	sudo apt-get install python3 python3-venv python3-pip

# Installs the prerequisites for a mac user
install-prerequisites-mac:
	sudo apt-get install python3 python3-venv python3-pip
