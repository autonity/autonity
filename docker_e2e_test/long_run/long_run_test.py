#!/usr/bin/env python3

import argparse
import copy
import time
import os
import getpass

import googleapiclient.discovery

from fabric import Connection


def get_user_ssh_pub_key(path):
    with open(path) as pub_key_reader:
        key = pub_key_reader.read()
    return key.rstrip("\n")


def create_instance_template(compute, project, user, ssh_key_path, default_template):
    key = get_user_ssh_pub_key(ssh_key_path)
    time_string = time.asctime().lower().replace(" ", "-").replace(":", "")
    new_template_name = 'template-{}'.format(time_string)

    # get default template from cloud.
    default = compute.instanceTemplates().get(project=project, instanceTemplate=default_template).execute()

    # set template with new name and user ssh pub key.
    new = copy.deepcopy(default)
    new['name'] = new_template_name
    for item in new['properties']['metadata']['items']:
        if item['key'] == 'ssh-keys':
            item['value'] = "{}\n{}:{}".format(item['value'], user, key)

    # create new template.
    compute.instanceTemplates().insert(project=project, body=new).execute()
    return new_template_name


def create_instance(compute, project, zone, name, template_name):
    config = {'name': name}
    template_uri = 'projects/{}/global/instanceTemplates/{}'.format(project, template_name)
    operation = compute.instances().insert(project=project, zone=zone, sourceInstanceTemplate=template_uri,
                                           body=config).execute()
    wait_for_operation(compute, project, zone, operation['name'])


def wait_for_operation(compute, project, zone, operation):
    print('Waiting for operation to finish ......')
    while True:
        result = compute.zoneOperations().get(
            project=project,
            zone=zone,
            operation=operation).execute()

        if result['status'] == 'DONE':
            print("Done.")
            if 'error' in result:
                raise Exception(result['error'])
            return result

        time.sleep(1)


def create_virtual_machine(compute, project, zone, instance_name, template, ssh_key, user):
    # create instance template.
    template_name = create_instance_template(compute, project, user, ssh_key, template)
    create_instance(compute, project, zone, instance_name, template_name)
    compute.instanceTemplates().delete(project=project, instanceTemplate=template_name).execute()

    # get instance meta data.
    response = compute.instances().get(project=project, zone=zone, instance=instance_name).execute()
    return response


def launch_test_engine(ip, ssh_key, user, branch):
    print("Waiting for virtual machine to finish bootstrap ......")
    time.sleep(60 * 5)
    i = 0
    while i < 10:
        i += 1
        try:
            with Connection(ip, user=user, connect_kwargs={
                "key_filename": ssh_key,
            }) as c:
                print("[STEP-1]: Remote installing python3, pip3, make, gcc and docker ......")
                r = c.run(
                    "sudo apt-get upgrade -y && sudo apt-get install -y python3 python3-pip make gcc docker.io"
                    " && sudo snap install go --classic", pty=True, warn=True, hide=True, asynchronous=False)

                if r and r.exited == 0 and r.ok:
                    print("[STEP-1]: Remote install dependencies done!")

                print("[STEP-2]: Remote build autonity binary .....")
                r = c.run("git clone https://github.com/clearmatics/autonity.git && cd autonity "
                          "&& git checkout {} && make all".format(branch), pty=True, warn=True, hide=True,
                          asynchronous=False)
                if r and r.exited == 0 and r.ok:
                    print("[STEP-2]: Remote build autonity binary done!")

                print("[STEP-3]: Remote Install requirements_docker_test.txt ......")
                r = c.run(
                    "sudo pip3 install -r /home/{}/autonity/docker_e2e_test/requirements_docker_test.txt".format(user),
                    pty=True, warn=True, hide=True, asynchronous=False)
                if r and r.exited == 0 and r.ok:
                    print("[STEP-3]: Remote Install requirements_docker_test.txt done!")

                print("[STEP-4]: Launching test engine ......")
                c.run("echo \"#!/usr/bin/env bash\" >> run.sh")
                c.run("echo \"cd /home/{}/autonity/docker_e2e_test/\" >> run.sh".format(user))
                c.run(
                    "echo \"sudo nohup python3 test_via_docker.py .. > test_report.log 2> stdout.er < /dev/null &\" >> run.sh")
                with c.cd("/home/{}/".format(user)):
                    r = c.run("sudo bash run.sh".format(user))
                    if r and r.exited == 0 and r.ok:
                        print("[STEP-4]: Launch test engine done!")
                break

        except Exception as e:
            if str(e).find("Unable to connect to port 22") != -1:
                print("Compute instance is starting, retry connecting...")
                time.sleep(5)
                continue


def main(compute, project, zone, instance_name, template_name, ssh_key, user, branch):
    # create compute instance.
    print("Creating virtual machine ......")
    instance = create_virtual_machine(compute, project, zone, instance_name, template_name, ssh_key, user)
    instance_ip = ""
    # deploy test framework on remote compute instance and run test.
    for item in instance["networkInterfaces"]:
        for i in item["accessConfigs"]:
            instance_ip = i["natIP"]
            break
        break

    print("Virtual machine created:\n\tproject: {}\n\tzone: {}\n\tvm: {}\n\tip: {}\n\tuser: {}".format(project, zone,
                                                                                                       instance_name,
                                                                                                       instance_ip,
                                                                                                       user))
    # launch test engine.
    launch_test_engine(instance_ip, ssh_key, user, branch)

    print("[TIP] Collect run time test report at {}:/home/{}/autonity/docker-e2e-test/test_report.log".format(
        instance_ip, user))
    print("[TIP] Once test case failed, the context and system logs are zipped at "
          "{}:/home/{}/autonity/docker-e2e-test/JOB_<timestamp>.tar".format(instance_ip, user))
    print("[TIP] please use \"ssh {}@{}\" to access remote compute instance".format(user, instance_ip))


if __name__ == '__main__':
    home_dir = os.path.expanduser("~")

    parser = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)

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

    args = parser.parse_args()

    # google cloud service finding and auth, it need user to install gcloud sdk, and launch gcloud login before running.
    service = googleapiclient.discovery.build('compute', 'v1')
    if args.rm_instance != '':
        print("Removing instance: {} from cloud".format(args.rm_instance))
        op = service.instances().delete(project=args.project_id, zone=args.zone, instance=args.rm_instance).execute()
        wait_for_operation(service, args.project_id, args.zone, op['name'])
        print("Instance removed.")
    else:
        print("****** Autonity long run test pipe line starting ******")
        print("Going to use below parameters to set up test.")
        print(
            "\tgcloud project id: {}\n \tzone: {}\n \tvm_name: {}\n \tvm_template: {}\n \tuser_ssh_key: {}\n \tuser: {}\n"
            " \trm_instance: {}\n \tautonity_branch: {}".format(args.project_id, args.zone, args.name, args.template,
                                                                args.ssh_key, args.user, args.rm_instance, args.branch))
        main(service, args.project_id, args.zone, args.name, args.template, args.ssh_key, args.user, args.branch)
