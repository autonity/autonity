import yaml
import log
import ipaddress

LOGGER = log.get_logger()

CONFIG_FILE = "./etc/e2e_test_conf.yml"
CONF = None


def write_yaml(file, data):
    try:
        with open(file, 'w') as f:
            yaml.dump(data, f, default_flow_style=False)
    except IOError as e:
        LOGGER.error('cannot write file: %s', e)
        return None
    return True


def get_engine_conf():
    return CONF


def get_testbed_conf_file_name():
    try:
        return CONF["generate_testbed_conf_at"]
    except KeyError as e:
        LOGGER.error("wrong config for generated test bed conf.", e)
        return None


def get_test_case_conf_file_name():
    try:
        return CONF["generate_testcase_conf_at"]
    except KeyError as e:
        LOGGER.error("wrong config for testcase conf file.", e)
        return None


def get_misbehaviour_test_case_conf_file_name():
    try:
        return CONF["generate_misbehaviour_testcase_conf_at"]
    except KeyError as e:
        LOGGER.error("wrong config for misbehaviour testcase conf file.", e)
        return None


def get_test_case_conf():
    file = get_test_case_conf_file_name()
    if file is None:
        return None
    return load_conf(file)


def get_misbehaviour_test_case_conf():
    file = get_misbehaviour_test_case_conf_file_name()
    if file is None:
        return None
    return load_conf(file)


def get_test_bed_conf():
    file = get_testbed_conf_file_name()
    if file is None:
        return None
    return load_conf(file)


def get_testbed_template():
    try:
        return CONF["test_bed_template"]
    except KeyError as e:
        LOGGER.error("wrong config for test bed template.", e)
        return None


def dump_test_bed_conf(data):
    file = get_testbed_conf_file_name()
    if file:
        return write_yaml(file, data)
    return None


def get_validator_host_name_by_ip(ip):
    try:
        validator_file = CONF["validator_ip_file"]
        with open(validator_file) as f:
            for line in f.readlines():
                parts = line.split(" ")
                if parts[0] == ip:
                    return parts[1][:-1]
    except Exception as e:
        LOGGER.error('Cannot parse validator host name from input file. %s', e)
        return None


def get_client_ips():
    """
    :return: ip list of validator, ip list of participant.
    """
    try:
        validator_file = CONF["validator_ip_file"]
        participant_file = CONF["participant_ip_file"]
        return parse_ip_from_text(validator_file), parse_ip_from_text(participant_file)
    except Exception as e:
        LOGGER.error("cannot read validator ips. %s", e)
        return None


def load_project_conf():
    global CONF
    CONF = load_conf(CONFIG_FILE)
    if CONF is None:
        LOGGER.error("Cannot load project config file.")
        exit(1)


def load_conf(file):
    try:
        with open(file) as f:
            conf = yaml.load(f, Loader=yaml.FullLoader)
            return conf
    except (IOError, OSError) as e:
        LOGGER.error("Cannot find conf file: %s, %s", file, e)
        return None
    except ValueError as e:
        LOGGER.error("Wrong file: %s, %s", file, e)
        return None


def parse_ip_from_text(file):
    """
    get public ip list from text.
    :param file:
    :return:
    """
    ip_set = set()
    try:
        with open(file) as f:
            for line in f.readlines():
                for part in line.split():
                    try:
                        a = ipaddress.ip_network(part)
                    except ValueError:
                        pass
                    else:
                        if a.is_private:
                            ip_set.add(str(a.network_address))
    except Exception as e:
        LOGGER.error('Cannot parse ip from input file. %s', e)
    finally:
        print('Get public IP from file, counted: ', len(ip_set), file)
        return sorted(ip_set)
