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


def extract_keys(output_string):
    keys = {}

    # Split the output string by lines
    lines = output_string.split('\n')

    # Iterate through the lines to find the lines containing the keys
    for line in lines:
        if "Node Address" in line:
            keys["Node Address"] = line.split(":")[1].strip()
        elif "Node Public Key" in line:
            keys["Node Public Key"] = line.split(":")[1].strip()
        elif "Consensus Public Key" in line:
            keys["Consensus Public Key"] = line.split(":")[1].strip()
        elif "Consensus Private Key" in line:
            keys["Consensus Private Key"] = line.split(":")[1].strip()
        elif "Node Private Key" in line:
            keys["Node Private Key"] = line.split(":")[1].strip()

    return keys

def gen_autonity_keys(autonity, key_inspector, key_file):
    # Generate a node-specific key file using autonity command
    autonity_command = f"{autonity} genAutonityKeys {key_file}"
    subprocess.run(autonity_command, shell=True)

    # Inspect the generated key file using key_inspector command with -private flag
    key_inspector_command = f"{key_inspector} autinspect {key_file} -private"
    process = subprocess.Popen(key_inspector_command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    output, error = process.communicate()

    if error:
        return None  # Return None if there was an error

    output_string = output.decode("utf-8")

    # Extract and return the keys as separate parameters
    keys = extract_keys(output_string)
    return keys.get("Node Address"), keys.get("Node Public Key"), keys.get("Consensus Public Key"), keys.get("Consensus Private Key"), keys.get("Node Private Key")
