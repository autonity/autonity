from typing import Any, Dict, Optional

import eth_utils


class Validator:
    """Interface for composable validators."""
    def validate(self, obj: Any):
        raise NotImplementedError()


class HexString(Validator):
    """Validator for a string field in a response in '0x123' format."""
    def __init__(self, length: Optional[int] = None):
        self.length = length

    def validate(self, obj: Any):
        try:
            eth_utils.to_int(hexstr=obj)
        except (TypeError, ValueError):
            assert False, f"'{obj}' is not a hex-string"
        if self.length is not None:
            assert len(eth_utils.to_bytes(hexstr=obj)) == self.length, \
                f"'{obj}' should be {self.length} bytes long"


class Array(Validator):
    """Validator for an array field in a response."""
    def __init__(
            self,
            elem_validator: Optional[Validator] = None,
            length: Optional[int] = None
        ):
        self.elem_validator = elem_validator
        self.length = length

    def validate(self, obj: Any):
        assert isinstance(obj, list)
        if self.length is not None:
            assert len(obj) == self.length
        if self.elem_validator is not None:
            for elem in obj:
                self.elem_validator.validate(elem)


class Object(Validator):
    """Validator for an object field in a response."""
    def __init__(self, fields: Dict[str, Validator]):
        self.fields = fields

    def validate(self, obj: Any):
        assert isinstance(obj, dict)
        for name, field_validator in self.fields.items():
            assert name in obj
            field_validator.validate(obj[name])


class Boolean(Validator):
    """Validator for a Boolean field in a response."""
    @staticmethod
    def validate(obj: Any):
        assert isinstance(obj, bool)


class Null(Validator):
    """Validator for a 'null' field in a response."""
    @staticmethod
    def validate(obj: Any):
        assert obj is None
