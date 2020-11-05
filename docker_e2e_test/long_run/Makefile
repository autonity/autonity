.PHONY: autonity_long_run_test clean install-prerequisites-linux install-prerequisites-mac

# Runs the autonity long run test for 7 days.
#
# By default each new line in a rule's commands is run in a separate shell. To
# run multiple commands in the same shell the commands need to be separated
# with ';' or '&&'. So any python commands that depend on a venv needs to be
# run in the same shell in which the venv is activated.
autonity-long-tests: ./venv/bin/activate
	. venv/bin/activate && sudo python3 long_run_test.py

# Creates venv if not present and then installs dependencies. Executes only if
# requirements.txt has been updated since the last time it ran.
#
# Make targets that start with a dot are considered hidden and are not provided
# as autocomplete options.
./venv/bin/activate: requirements.txt
	test -d venv || python3 -m venv venv
	. venv/bin/activate && sudo pip3 install -r requirements.txt
	touch ./venv/bin/activate

clean:
	rm -rf venv

# Installs the prerequisites for a Linux user
install-prerequisites-linux:
	sudo apt-get update
	sudo apt-get install python3 python3-venv python3-pip google-cloud-sdk

# Installs the prerequisites for a mac user
install-prerequisites-mac:
	brew update
	brew install python3 google-cloud-sdk
