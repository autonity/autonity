import subprocess


def execute(cmd):
    print("[CMD] {}".format(cmd))
    process = subprocess.Popen(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, encoding="utf-8", shell=True)
    return process.communicate(input='\n')


def create_dir(dir_name):
    execute("mkdir -p {}".format(dir_name))


def remove_dir(dir_name):
    execute("rm -rf {}".format(dir_name))


def create_network_dir(ip_list):
    remove_dir("./network-data")
    create_dir("./network-data")
    for ip in ip_list:
        create_dir("./network-data/{}".format(ip))
