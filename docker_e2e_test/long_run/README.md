# Autonity long run e2e test pipe line script
The script create a virtual machine(VM) under the google cloud project: autonity-e2e-918063 on zone europe-west2-c by 
default with a VM template: test-engine-template attaching with user's ssh public key. After the creation, the script 
start to deploy autonity docker e2e test framework on remote VM, and launch a long running test.

# Prerequisites

`make install-prerequisites-[linux|mac]`
## gcloud SDK
Run `gcloud projects list` to double check if "autonity-e2e-918063" is presented at your gcloud SDK. Otherwise you need:
install gcloud SDK by referring from: https://cloud.google.com/sdk/docs/install and init your gcloud SDK with your google
 account and cloud project-id: autonity-e2e-918063 by referring from: https://cloud.google.com/sdk/docs/initializing
## Get gcloud API token
### Create service account
`gcloud iam service-accounts create somename`
### Bind service account to project
`gcloud projects add-iam-policy-binding autonity-e2e-918063 --member="serviceAccount:somename@autonity-e2e-918063.iam.gserviceaccount.com" --role="roles/owner"`
### Generate API token
`gcloud iam service-accounts keys create your_token_file.json --iam-account=somename@autonity-e2e-918063.iam.gserviceaccount.com`
## Set API token file path in environment variable
For example:
`export GOOGLE_APPLICATION_CREDENTIALS="/home/user/Downloads/your_token_file.json"`
The client api will read environment variable GOOGLE_APPLICATION_CREDENTIALS to complete the service finding and auth. 

## 3rd party python libs
Execute below command:
`pip3 install -r requirements.txt`
to install:
```
google-api-python-client==1.12.4
google-auth==1.22.1
google-auth-httplib2==0.0.4
fabric==2.5.0
```
# How to use it
## Run
`make autonity-long-tests`
It start the test by using master branch of autonity by default.
## Parameters and default values
```python
    parser.add_argument('--project_id', default='autonity-e2e-918063', help='Your Google Cloud project ID.')
    parser.add_argument('--zone', default='europe-west2-c', help='Compute Engine zone to deploy to.')
    parser.add_argument('--name',
                        default='test-engine-{}'.format(time.asctime().lower().replace(" ", "-").replace(":", "")),
                        help='New compute instance name.')
    parser.add_argument('--template', default='test-engine-template', help='Compute instance template name.')
    parser.add_argument('--ssh_key', default="{}/.ssh/id_rsa.pub".format(home_dir),
                        help='SSH public key for accessing remote VM.')
    parser.add_argument('--user', default=getpass.getuser(), help='SSH public key for accessing remote VM.')
    parser.add_argument('--rm_instance', default='', help='Name of compute instance to be removed.')
    parser.add_argument('--branch', default='master', help='Branch name of autonity to be tested by test engine.')
```

# Example outputs
```
****** Autonity long run test pipe line starting ******
Going to use below parameters to set up test.
	gcloud project id: autonity-e2e-918063
 	zone: europe-west2-c
 	vm_name: test-engine-fri-oct-23-005200-2020
 	vm_template: test-engine-template
 	user_ssh_key: /home/jason.chen/.ssh/id_rsa.pub
 	user: jason.chen
 	rm_instance: 
 	autonity_branch: master
Creating virtual machine ......
Waiting for operation to finish ......
Done.
Virtual machine created:
	project: autonity-e2e-918063
	zone: europe-west2-c
	vm: test-engine-fri-oct-23-005200-2020
	ip: 35.197.218.5
	user: jason.chen
Waiting for virtual machine to finish bootstrap ......
[STEP-1]: Remote installing python3, pip3, make, gcc and docker ......
[STEP-1]: Remote install dependencies done!
[STEP-2]: Remote build autonity binary .....
[STEP-2]: Remote build autonity binary done!
[STEP-3]: Remote Install requirements_docker_test.txt ......
[STEP-3]: Remote Install requirements_docker_test.txt done!
[STEP-4]: Launching test engine ......
[STEP-4]: Launch test engine done!
[TIP] Collect run time test report at 35.197.218.5:/home/jason.chen/autonity/docker-e2e-test/test_report.log
[TIP] Once test case failed, network system logs can be find at test_report.log or zipped at 35.197.218.5:/home/jason.chen/autonity/docker-e2e-test/JOB_<timestamp>.tar
[TIP] please use "ssh jason.chen@35.197.218.5" to access remote compute instance
```
