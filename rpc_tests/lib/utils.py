import time
from typing import Any, Callable


def repeat_until(
        callback: Callable,
        timeout: float,
        interval: float = 0.1
    ) -> Any:
    """Executes the callback repeatedly with the given interval until the
    callback returns a truthy or the timer expires."""
    start_time = time.time()
    while time.time() - start_time <= timeout:
        result = callback()
        if result:
            return result
        time.sleep(interval)
    return None
