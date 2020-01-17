# How it works
The script will setup 4 clients network in local host, 
and run all the test cases from end to end point of view
by using web3 libs to test all the functions exposed by 
autonity contract.

# How to setup in your local
```
sudo apt-get update -y
sudo apt-get install -y tmux
sudo apt-get install python3
sudo apt-get install python3-pip
cd ~/go/src/github.com/clearmatics/autonity/contract_e2e_test/
pip3 install -r requirements.txt

```
# Integrate with CI script
CI script will add new task to trigger this automation test.
Makefile script will be updated too.
