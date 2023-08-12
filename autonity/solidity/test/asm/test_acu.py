import dataclasses
from collections import namedtuple
from decimal import Decimal

import ape
import pytest
from ape.api import AccountAPI
from web3 import Web3
from web3.constants import ADDRESS_ZERO

INT256_MAX = (
    57896044618658097711785492504343953926634992332820282019728792003956564819967
)
ORACLE_SCALE_FACTOR = int(10**7)

ACUTestData = namedtuple(
    "ACUTestData",
    ["contract", "oracle", "symbols", "quantities", "prices", "scale_factor"],
)

INVALID_BASKET_ERROR = "0x4ff799c5"


@dataclasses.dataclass
class Accounts:
    deployer: AccountAPI
    autonity: AccountAPI
    # ACU Contract
    operator: AccountAPI
    caller: AccountAPI
    # Oracle Contract
    voter: AccountAPI


@pytest.fixture
def users(accounts):
    return Accounts(
        deployer=accounts[0],
        autonity=accounts[1],
        operator=accounts[2],
        caller=accounts[3],
        voter=accounts[4],
    )


@pytest.fixture
def oracle_factory(project, users):
    def mkoracle(symbols, prices=(), voting_period=60):
        contract = project.Oracle.deploy(
            [users.voter],
            users.autonity,
            users.operator,
            symbols,
            voting_period,
            sender=users.deployer,
        )
        if prices:
            # commit: (prices, salt, voter_address)
            commit = Web3.solidity_keccak(
                ["int256[]", "uint256", "address"],
                [prices, 123, users.voter.address],
            )
            contract.vote(commit, [], 0, sender=users.voter)  # commit
            ape.chain.mine(voting_period)
            contract.finalize(sender=users.autonity)
            contract.vote(
                Web3.solidity_keccak([], []), prices, 123, sender=users.voter
            )  # reveal
            ape.chain.mine(voting_period)
            contract.finalize(sender=users.autonity)
        return contract

    return mkoracle


@pytest.fixture
def tracing(networks):
    provider = networks.active_provider
    if not provider.supports_tracing:
        pytest.skip(f"provider '{provider.name}' doesn't support tracing")


@pytest.fixture
def acu_basic(project, users, oracle_factory):
    symbols = ["AAA", "BBB", "CCC"]
    oracle = oracle_factory(symbols)
    scale = 5
    scale_factor = int(1e5)
    quantities = [i * scale_factor for i in range(len(symbols))]
    contract = project.ACU.deploy(
        symbols,
        quantities,
        scale,
        users.autonity,
        users.operator,
        oracle,
        sender=users.deployer,
    )
    return ACUTestData(
        contract=contract,
        oracle=oracle,
        symbols=symbols,
        quantities=quantities,
        prices=None,
        scale_factor=scale_factor,
    )


@pytest.fixture
def acu_primed(project, users, oracle_factory):
    scale = 5
    scale_factor = int(10**scale)
    symbols = [
        "AUD/USD",
        "CAD/USD",
        "EUR/USD",
        "GBP/USD",
        "JPY/USD",
        "USD/USD",
        "SEK/USD",
    ]
    prices = [
        int(price * ORACLE_SCALE_FACTOR)
        for price in (
            Decimal("0.6757"),
            Decimal("0.75694"),
            Decimal("1.1085"),
            Decimal("1.29403"),
            Decimal("0.00713"),
            Decimal(1),
            Decimal("0.09597"),
        )
    ]
    quantities = [
        int(quantity * Decimal(scale_factor))
        for quantity in (
            Decimal("0.213"),
            Decimal("0.187"),
            Decimal("0.143"),
            Decimal("0.104"),
            Decimal("17.6"),
            Decimal("0.180"),
            Decimal("1.41"),
        )
    ]
    oracle = oracle_factory(symbols, prices, voting_period=30)
    for symbol, price in zip(symbols, prices):
        assert oracle.latestRoundData(symbol).price == price
    contract = project.ACU.deploy(
        symbols,
        quantities,
        scale,
        users.autonity,
        users.operator,
        oracle,
        sender=users.deployer,
    )
    return ACUTestData(
        contract=contract,
        oracle=oracle,
        symbols=symbols,
        quantities=quantities,
        prices=prices,
        scale_factor=scale_factor,
    )


def test_constructor(acu_basic):
    assert acu_basic.contract.round() == 0
    assert acu_basic.contract.scaleFactor() == acu_basic.scale_factor


def test_constructor_invalid_size(project, users):
    with ape.reverts(INVALID_BASKET_ERROR):
        project.ACU.deploy(
            ["FOO"],
            [],
            1,
            ADDRESS_ZERO,
            ADDRESS_ZERO,
            ADDRESS_ZERO,
            sender=users.deployer,
        )
    with ape.reverts(INVALID_BASKET_ERROR):
        project.ACU.deploy(
            [],
            [1, 2, 3],
            1,
            ADDRESS_ZERO,
            ADDRESS_ZERO,
            ADDRESS_ZERO,
            sender=users.deployer,
        )


def test_constructor_bad_quantity(project, users):
    with ape.reverts(INVALID_BASKET_ERROR):
        project.ACU.deploy(
            ["foo"],
            [1 + INT256_MAX],
            1,
            ADDRESS_ZERO,
            ADDRESS_ZERO,
            ADDRESS_ZERO,
            sender=users.deployer,
        )


def test_modify_basket(acu_basic, users):
    new_basket = ["X", "Y"]
    new_quantities = [1, 2]
    new_scale = 18
    new_scale_factor = int(10**new_scale)
    receipt = acu_basic.contract.modifyBasket(
        new_basket, new_quantities, new_scale, sender=users.operator
    )
    assert acu_basic.contract.symbols() == new_basket
    assert acu_basic.contract.quantities() == new_quantities
    assert acu_basic.contract.scaleFactor() == new_scale_factor
    assert len(receipt.events) == 1
    event = receipt.events[0]
    assert event.event_name == "BasketModified"
    assert event.symbols == tuple(new_basket)
    assert event.quantities == tuple(new_quantities)
    assert event.scale == new_scale


def test_modify_basket_unauthorized(acu_basic, users):
    for unauth_user in [users.deployer, users.autonity, users.caller]:
        with ape.reverts(acu_basic.contract.Unauthorized):
            acu_basic.contract.modifyBasket(["X", "Y"], [1, 2], 1, sender=unauth_user)


def test_modify_basket_invalid_size(acu_basic, users):
    with ape.reverts(acu_basic.contract.InvalidBasket):
        acu_basic.contract.modifyBasket(["FOO"], [], 1, sender=users.operator)
    with ape.reverts(acu_basic.contract.InvalidBasket):
        acu_basic.contract.modifyBasket([], [1], 1, sender=users.operator)


def test_modify_basket_bad_quantity(acu_basic, users):
    with ape.reverts(acu_basic.contract.InvalidBasket):
        acu_basic.contract.modifyBasket(
            ["FOO"], [1 + INT256_MAX], 1, sender=users.operator
        )


def test_set_operator(acu_basic, accounts, users):
    new_operator = accounts.generate_test_account()
    acu_basic.contract.setOperator(new_operator, sender=users.autonity)
    users.deployer.transfer(new_operator, int(1e18))
    acu_basic.contract.modifyBasket(["A"], [1], 1, sender=new_operator)


def test_set_operator_unauthorized(acu_basic, accounts, users):
    new_operator = accounts.generate_test_account()
    for unauth_user in [users.deployer, users.operator, users.caller]:
        with ape.reverts(acu_basic.contract.Unauthorized):
            acu_basic.contract.setOperator(new_operator, sender=unauth_user)


def test_set_oracle(acu_basic, oracle_factory, users):
    symbols = ["FOO"]
    prices = [123]
    voting_period = 1
    new_oracle = oracle_factory(symbols, prices, voting_period)
    for symbol, price in zip(symbols, prices):
        assert new_oracle.latestRoundData(symbol).price == price
    acu_basic.contract.setOracle(new_oracle, sender=users.autonity)


def test_set_oracle_unauthorized(acu_basic, oracle_factory, users):
    new_oracle = oracle_factory(["FOO"])
    for unauth_user in [users.deployer, users.operator, users.caller]:
        with ape.reverts(acu_basic.contract.Unauthorized):
            acu_basic.contract.setOracle(new_oracle, sender=unauth_user)


def test_value_novalue(acu_basic):
    with ape.reverts(acu_basic.contract.NoACUValue):
        acu_basic.contract.value()


def test_update(acu_primed, users, chain, tracing):
    receipt = acu_primed.contract.update(sender=users.autonity)
    assert receipt.return_value is True
    value = compute_acu(
        acu_primed.symbols,
        acu_primed.prices,
        acu_primed.quantities,
    )
    oracle_round = acu_primed.oracle.getRound() - 1
    assert acu_primed.contract.round() == oracle_round
    assert acu_primed.contract.value() == value
    assert len(receipt.events) == 1
    modified = receipt.events[0]
    assert modified.event_name == "Updated"
    assert modified.height == chain.blocks.head.number
    assert modified.timestamp == chain.blocks.head.timestamp
    assert modified.value == value
    assert modified.round == oracle_round


def test_update_unauthorized(acu_primed, users):
    for unauth_user in [users.deployer, users.operator, users.caller]:
        with ape.reverts(acu_primed.contract.Unauthorized):
            acu_primed.contract.update(sender=unauth_user)


def test_update_same_round(acu_primed, accounts, users, tracing):
    with accounts.use_sender(users.autonity):
        receipt = acu_primed.contract.update()
        assert receipt.return_value is True
        round_before = acu_primed.contract.round()
        receipt = acu_primed.contract.update()
        assert receipt.return_value is False
        assert acu_primed.contract.round() == round_before


def test_update_missing_price(acu_basic, users, tracing):
    vote_period = acu_basic.oracle.votePeriod()
    ape.chain.mine(vote_period)
    acu_basic.oracle.finalize(sender=users.autonity)
    for symbol in acu_basic.symbols:
        assert acu_basic.oracle.latestRoundData(symbol).status == 1
    receipt = acu_basic.contract.update(sender=users.autonity)
    assert receipt.return_value is False
    with ape.reverts(acu_basic.contract.NoACUValue):
        acu_basic.contract.value(sender=users.caller)


def test_view_functions(acu_basic):
    assert acu_basic.contract.symbols() == acu_basic.symbols
    assert acu_basic.contract.quantities() == acu_basic.quantities
    assert acu_basic.contract.scaleFactor() == acu_basic.scale_factor


def compute_acu(symbols, prices, quantities):
    value = 0
    for s, p, q in zip(symbols, prices, quantities):
        if s == "USD/USD":
            p = ORACLE_SCALE_FACTOR
        pscaled = int(Decimal(p) * Decimal(q))
        value += pscaled
    value = int(value / ORACLE_SCALE_FACTOR)
    return value
