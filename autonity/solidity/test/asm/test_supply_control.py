import dataclasses

import ape
import pytest
from ape.api import AccountAPI
from web3.constants import ADDRESS_ZERO

TOTAL_SUPPLY = int(1e8)


@dataclasses.dataclass
class Accounts:
    deployer: AccountAPI
    admin: AccountAPI
    operator: AccountAPI
    unauthorized: AccountAPI
    unfunded: AccountAPI


@pytest.fixture
def users(accounts):
    assert accounts[0].balance > TOTAL_SUPPLY
    return Accounts(
        deployer=accounts[0],
        admin=accounts[1],
        operator=accounts[2],
        unauthorized=accounts[3],
        unfunded=accounts.generate_test_account(),
    )


@pytest.fixture
def supply_control(project, users):
    return project.SupplyControl.deploy(
        users.admin,
        users.operator,
        sender=users.deployer,
        value=TOTAL_SUPPLY,
    )


def test_deploy(supply_control, users):
    assert supply_control.operator() == users.operator
    assert supply_control.totalSupply() == TOTAL_SUPPLY
    assert supply_control.availableSupply() == TOTAL_SUPPLY


def test_deploy_zero_value(project, users):
    with ape.reverts("0x7c946ed7"):
        project.SupplyControl.deploy(
            users.deployer,
            users.operator,
            sender=users.admin,
            value=0,
        )


def test_set_operator(supply_control, users):
    assert supply_control.operator() == users.operator


def test_set_operator_unauthorized(supply_control, users):
    with ape.reverts(supply_control.Unauthorized):
        supply_control.setOperator(users.operator, sender=users.unauthorized)


def test_mint_authorized(supply_control, users):
    supply_control.mint(users.unfunded, 1, sender=users.operator)


def test_mint_unauthorized(supply_control, users):
    for sender in (users.deployer, users.admin, users.unauthorized):
        with ape.reverts(supply_control.Unauthorized):
            supply_control.mint(users.unfunded, 1, sender=sender)


def test_mint_invalid_recipient(supply_control, users):
    with ape.reverts(supply_control.InvalidRecipient):
        supply_control.mint(ADDRESS_ZERO, 1, sender=users.operator)
    with ape.reverts(supply_control.InvalidRecipient):
        supply_control.mint(users.operator, 1, sender=users.operator)


def test_mint_valid_amount(supply_control, users):
    supply_control.mint(users.unfunded, 1, sender=users.operator)
    assert users.unfunded.balance == 1
    supply_control.mint(users.unfunded, TOTAL_SUPPLY - 1, sender=users.operator)
    assert users.unfunded.balance == TOTAL_SUPPLY


def test_mint_invalid_amount(supply_control, users):
    with ape.reverts(supply_control.InvalidAmount):
        supply_control.mint(users.unfunded, 0, sender=users.operator)
    with ape.reverts(supply_control.InvalidAmount):
        supply_control.mint(users.unfunded, TOTAL_SUPPLY + 1, sender=users.operator)


def test_mint_event(supply_control, users):
    receipt = supply_control.mint(users.unfunded, 1, sender=users.operator)
    assert len(receipt.events) == 1
    mint_event = receipt.events[0]
    assert mint_event.event_name == "Mint"
    assert mint_event.recipient == users.unfunded
    assert mint_event.amount == 1


def test_burn(supply_control, users):
    supply_control.burn(sender=users.operator, value=1)


def test_burn_unauthorized(supply_control, users):
    for sender in (users.deployer, users.admin, users.unauthorized):
        with ape.reverts(supply_control.Unauthorized):
            supply_control.burn(sender=sender, value=1)


def test_burn_event(supply_control, users):
    receipt = supply_control.burn(sender=users.operator, value=1)
    assert len(receipt.events) == 1
    burn_event = receipt.events[0]
    assert burn_event.event_name == "Burn"
    assert burn_event.amount == 1
