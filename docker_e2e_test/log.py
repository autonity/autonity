import logging
import logging.config
import yaml

LOGGING_CONF = './etc/logconf.yaml'


def get_logger(name=None):
    try:
        with open(LOGGING_CONF, 'r') as f:
            config = yaml.safe_load(f.read())
            logging.config.dictConfig(config)
            return logging.getLogger(name)
    except IOError as e:
        print("Cannot find log conf file. %s", e)
        return None
    except ValueError as e:
        print("Wrong file of log conf. %s", e)
        return None
