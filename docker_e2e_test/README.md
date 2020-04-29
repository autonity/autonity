# ATAutonity
Python3 and docker based e2e testing framework for autonity.

# Dependencies
Run install_dep.sh to install all the dependencies. Or you can install them manually as below steps tells:
In most linux distribution, python3 and pip3 are included, follow below guide in case your linux need them:
## update repo source
`sudo apt-get update`
## Python3
`sudo apt-get install python3`
## pip3
`sudo apt-get install python3-pip`

## [required] 3rd party python libs
`pip3 install -r requirements_docker_test.txt`

## [required] docker
If your linux is on ubuntu-18.04, the script will auto install it for you.
`sudo apt-get install --yes docker.io`

# How to use it
## block chain genesis configuration
In the docker_e2e_test/planner/networkplanner.py + line: 122, there are default genesis configuration, update it on your demand.

## run it
You can run it via make:
`make docker-e2e-test`

or run the script with 1 parameter, the path of autonity binary.
For example:
`sudo python3 test_via_docker.py  ~/your_path_to/autonity`

# Outputs and reports.
The console will collect test report and it collects system logs of each autontiy client for per failed testcase.
The log will be compress in local dir name with: JOB_<job_id>.tar
