import dataclasses
import functools
import pathlib
from collections import namedtuple
from contextlib import contextmanager
from decimal import ROUND_FLOOR, Decimal

import ape
import pytest
from web3 import Web3

ATN_TOTAL_SUPPLY = int(100 * 1e18)
ERC20_CONTRACT = ape.project.path / pathlib.Path("autonity/solidity/test/asm/ERC20Basic.sol")
ORACLE_SCALE = 7  # oracle contract constant
PRB_MATH_E = Decimal(2.718281828459045235)  # match prb-math exactly
PRICE_NTN = int(1.234567 * 10**ORACLE_SCALE)
PRICE_NTN_18D = PRICE_NTN * 10 ** (18 - ORACLE_SCALE)
SCALE = 18
SCALING_FACTOR = int(10**SCALE)
SECONDS_IN_YEAR = Decimal(365 * 24 * 60 * 60)


# ┌──────────┐
# │ Fixtures │
# └──────────┘


@pytest.fixture
def basic_config():
    return Config(
        borrowInterestRate=int(0.05e18),
        liquidationRatio=int(1.5e18),
        minCollateralisationRatio=int(2.5e18),
        minDebtRequirement=int(1e6),  # megaton
        redemptionPrice=int(1e18),
    )


@pytest.fixture
def users(accounts):
    wallet = Accounts(
        deployer=accounts[0],
        operator=accounts[1],
        faucet=accounts[2],
        account1=accounts[3],
        account2=accounts[4],
        autonity=accounts[5],
        voter=accounts[6],
    )
    assert wallet.autonity.balance > ATN_TOTAL_SUPPLY
    return wallet


@pytest.fixture
def collateral_token(project, users):
    return users.faucet.deploy(erc20_contract())


@pytest.fixture
def oracle_factory(project, users, chain):
    def mkoracle(price, voting_period=1):
        price = int(price)
        contract = project.Oracle.deploy(
            [users.voter],
            users.autonity,
            users.autonity,
            ["NTN/ATN"],
            voting_period,
            sender=users.deployer,
        )
        # commit: (prices, salt, voter_address)
        commit = Web3.solidity_keccak(
            ["uint256[]", "uint256", "address"],
            [[price], 123, users.voter.address],
        )
        contract.vote(commit, [], 0, sender=users.voter)  # commit
        ape.chain.mine()
        contract.finalize(sender=users.autonity)
        contract.vote(
            Web3.solidity_keccak([], []), [price], 123, sender=users.voter
        )  # reveal
        ape.chain.mine(voting_period)
        contract.finalize(sender=users.autonity)
        return contract

    return mkoracle


@pytest.fixture
def supply_control(project, users):
    return project.SupplyControl.deploy(
        users.autonity,
        sender=users.deployer,
        value=ATN_TOTAL_SUPPLY,
    )


@pytest.fixture
def stabilisation(
    project, users, collateral_token, basic_config, oracle_factory, supply_control
):
    oracle = oracle_factory(PRICE_NTN)
    contract = project.Stabilisation.deploy(
        dataclasses.asdict(basic_config),
        users.operator,
        oracle.address,
        supply_control.address,
        collateral_token,
        sender=users.deployer,
    )
    supply_control.setOperator(contract.address, sender=users.autonity)
    return contract


@pytest.fixture
def funded_accounts(stabilisation, collateral_token, users):
    funded_amount = 100 * 10 ** collateral_token.decimals()
    collateral_token.transfer(users.account1, funded_amount, sender=users.faucet)
    collateral_token.transfer(users.account2, funded_amount, sender=users.faucet)
    collateral_token.approve(
        stabilisation.address, funded_amount, sender=users.account1
    )
    collateral_token.approve(
        stabilisation.address, funded_amount, sender=users.account2
    )
    return int(funded_amount)


@pytest.fixture
def deposit_scenario(stabilisation, funded_accounts, collateral_token, users):
    scenario = namedtuple(
        "DepositScenario",
        [
            "user",
            "funded_amount",
            "deposit_amount",
            "borrow_limit",
        ],
    )(
        user=users.account1,
        funded_amount=funded_accounts,
        deposit_amount=int(funded_accounts / 10),
        borrow_limit=int(Decimal("4.938268E+18")),
    )
    stabilisation.deposit(scenario.deposit_amount, sender=scenario.user)
    assert stabilisation.cdps(scenario.user).collateral == scenario.deposit_amount
    return scenario


@pytest.fixture
def borrow_scenario(stabilisation, deposit_scenario):
    scenario = namedtuple(
        "BorrowScenario",
        deposit_scenario._fields + ("borrow_amount",),
    )(
        borrow_amount=int(deposit_scenario.borrow_limit / 2),
        **deposit_scenario._asdict(),
    )
    stabilisation.borrow(scenario.borrow_amount, sender=scenario.user)
    assert stabilisation.cdps(scenario.user).principal == scenario.borrow_amount
    assert stabilisation.debtAmount(scenario.user) == scenario.borrow_amount
    return scenario


# ┌─────────────┐
# │ Constructor │
# └─────────────┘


def test_constructor_zero_mcr(
    project, users, collateral_token, basic_config, oracle_factory, supply_control
):
    basic_config.minCollateralisationRatio = 0
    oracle = oracle_factory(PRICE_NTN)  # InvalidParameter
    with ape.reverts("0x613970e0"):
        project.Stabilisation.deploy(
            dataclasses.asdict(basic_config),
            collateral_token,
            users.operator,
            oracle.address,
            supply_control.address,
            sender=users.deployer,
        )


def test_constructor_invalid_ratios(
    project, users, collateral_token, basic_config, oracle_factory, supply_control
):
    oracle = oracle_factory(PRICE_NTN)
    # liquidationRatio == minCollateralisationRatio
    basic_config.liquidationRatio = int(1e18)
    basic_config.minCollateralisationRatio = int(1e18)
    with ape.reverts("0x613970e0"):  # InvalidParameter
        project.Stabilisation.deploy(
            dataclasses.asdict(basic_config),
            users.operator,
            oracle.address,
            supply_control.address,
            collateral_token,
            sender=users.deployer,
        )
    # liquidationRatio > minCollateralisationRatio
    basic_config.liquidationRatio += 1
    with ape.reverts("0x613970e0"):  # InvalidParameter
        project.Stabilisation.deploy(
            dataclasses.asdict(basic_config),
            users.operator,
            oracle.address,
            supply_control.address,
            collateral_token,
            sender=users.deployer,
        )


# ┌─────────┐
# │ Deposit │
# └─────────┘


def verify_deposit_event(stabilisation, receipt, account, amount):
    assert len(receipt.events) == 2
    transfer_event = receipt.events[0]
    assert transfer_event.event_name == "Transfer"
    assert transfer_event.__getattr__("from") == account
    assert transfer_event.to == stabilisation
    assert transfer_event.value == amount
    deposit_event = receipt.events[1]
    assert deposit_event.event_name == "Deposit"
    assert deposit_event.account == account
    assert deposit_event.amount == amount


def test_deposit_zero(stabilisation, users):
    with ape.reverts(stabilisation.InvalidAmount):
        stabilisation.deposit(0, sender=users.account1)


def test_deposit_initial(stabilisation, funded_accounts, collateral_token, users):
    funded_amount = funded_accounts
    deposit_amount = int(funded_amount / 2)
    assert stabilisation.accounts() == []
    with check_token_transfer(
        collateral_token, users.account1, stabilisation, deposit_amount
    ):
        receipt = stabilisation.deposit(deposit_amount, sender=users.account1)
    cdp = stabilisation.cdps(users.account1)
    assert collateral_token.balanceOf(users.account1) == funded_amount - deposit_amount
    assert cdp.timestamp > 0
    assert cdp.collateral == deposit_amount
    assert stabilisation.accounts() == [users.account1]
    verify_deposit_event(stabilisation, receipt, users.account1, deposit_amount)


def test_deposit_subsequent(stabilisation, deposit_scenario, collateral_token, users):
    with check_token_transfer(
        collateral_token, users.account1, stabilisation, deposit_scenario.deposit_amount
    ):
        receipt = stabilisation.deposit(
            deposit_scenario.deposit_amount, sender=deposit_scenario.user
        )
    cdp = stabilisation.cdps(deposit_scenario.user)
    assert (
        collateral_token.balanceOf(deposit_scenario.user)
        == 8 * deposit_scenario.deposit_amount
    )
    assert cdp.timestamp > 0
    assert cdp.collateral == 2 * deposit_scenario.deposit_amount
    assert stabilisation.accounts() == [deposit_scenario.user]  # not duplicated
    verify_deposit_event(
        stabilisation, receipt, deposit_scenario.user, deposit_scenario.deposit_amount
    )


def test_deposit_insufficient_funds(stabilisation, users):
    with ape.reverts(stabilisation.InsufficientAllowance):
        stabilisation.deposit(1, sender=users.account1)


def test_deposit_insufficient_allowance(stabilisation, deposit_scenario):
    with ape.reverts(stabilisation.InsufficientAllowance):
        stabilisation.deposit(
            deposit_scenario.funded_amount - deposit_scenario.deposit_amount + 1,
            sender=deposit_scenario.user,
        )


def test_deposit_second_user(stabilisation, deposit_scenario, collateral_token, users):
    other_user = users.account2
    assert deposit_scenario.user != other_user
    assert stabilisation.cdps(other_user).collateral == 0
    with check_token_transfer(
        collateral_token, other_user, stabilisation, deposit_scenario.deposit_amount
    ):
        receipt = stabilisation.deposit(
            deposit_scenario.deposit_amount, sender=other_user
        )
    assert (
        collateral_token.balanceOf(other_user)
        == deposit_scenario.funded_amount - deposit_scenario.deposit_amount
    )
    cdp = stabilisation.cdps(other_user)
    assert cdp.timestamp > 0
    assert cdp.collateral == deposit_scenario.deposit_amount
    assert stabilisation.accounts() == [deposit_scenario.user, other_user]
    verify_deposit_event(
        stabilisation, receipt, other_user, deposit_scenario.deposit_amount
    )


# ┌──────────┐
# │ Withdraw │
# └──────────┘


def test_withdraw_zero(stabilisation, deposit_scenario):
    with ape.reverts(stabilisation.InvalidAmount):
        stabilisation.withdraw(0, sender=deposit_scenario.user)


def test_withdraw_full_deposit(stabilisation, collateral_token, deposit_scenario):
    with check_token_transfer(
        collateral_token,
        stabilisation,
        deposit_scenario.user,
        deposit_scenario.deposit_amount,
    ):
        receipt = stabilisation.withdraw(
            deposit_scenario.deposit_amount, sender=deposit_scenario.user
        )
    cdp = stabilisation.cdps(deposit_scenario.user)
    assert cdp.collateral == 0
    assert len(receipt.events) == 2
    transfer_event = receipt.events[0]
    assert transfer_event.event_name == "Transfer"
    assert transfer_event.__getattr__("from") == stabilisation.address
    assert transfer_event.to == deposit_scenario.user
    assert transfer_event.value == deposit_scenario.deposit_amount
    withdraw_event = receipt.events[1]
    assert withdraw_event.event_name == "Withdraw"
    assert withdraw_event.account == deposit_scenario.user
    assert withdraw_event.amount == deposit_scenario.deposit_amount


def test_withdraw_overdrawn(stabilisation, deposit_scenario):
    with ape.reverts(stabilisation.InvalidAmount):
        stabilisation.withdraw(
            1 + deposit_scenario.deposit_amount, sender=deposit_scenario.user
        )


def test_withdraw_liquidatable(stabilisation, users, deposit_scenario):
    stabilisation.borrow(deposit_scenario.borrow_limit, sender=deposit_scenario.user)
    assert not stabilisation.isLiquidatable(deposit_scenario.user)
    mcr = stabilisation.config().minCollateralisationRatio
    stabilisation.setMinCollateralisationRatio(
        mcr + SCALING_FACTOR, sender=users.operator
    )
    stabilisation.setLiquidationRatio(mcr, sender=users.operator)
    assert stabilisation.isLiquidatable(deposit_scenario.user)
    with ape.reverts(stabilisation.Liquidatable):
        stabilisation.withdraw(1, sender=deposit_scenario.user)


def test_withdraw_insufficient_collateral1(stabilisation, users, deposit_scenario):
    stabilisation.borrow(deposit_scenario.borrow_limit, sender=deposit_scenario.user)
    mcr = stabilisation.config().minCollateralisationRatio
    stabilisation.setMinCollateralisationRatio(
        mcr + SCALING_FACTOR, sender=users.operator
    )
    with ape.reverts(stabilisation.InsufficientCollateral):
        stabilisation.withdraw(1, sender=deposit_scenario.user)


def test_withdraw_insufficient_collateral2(stabilisation, deposit_scenario):
    borrow_amount = int(deposit_scenario.borrow_limit / 2)
    collateral_required = stabilisation.minimumCollateral(
        borrow_amount, PRICE_NTN_18D, stabilisation.config().minCollateralisationRatio
    )
    withdraw_max = deposit_scenario.deposit_amount - collateral_required
    stabilisation.borrow(borrow_amount, sender=deposit_scenario.user)
    with ape.reverts(stabilisation.InsufficientCollateral):
        stabilisation.withdraw(1 + withdraw_max, sender=deposit_scenario.user)


# ┌────────┐
# │ Borrow │
# └────────┘


def verify_borrow_event(deposit_scenario, receipt, borrow_amount):
    assert len(receipt.events) == 2
    mint_event = receipt.events[0]
    assert mint_event.event_name == "Mint"
    assert mint_event.recipient == deposit_scenario.user
    assert mint_event.amount == borrow_amount
    borrow_event = receipt.events[1]
    assert borrow_event.event_name == "Borrow"
    assert borrow_event.account == deposit_scenario.user
    assert borrow_event.amount == borrow_amount


def test_borrow_zero(stabilisation, deposit_scenario):
    with ape.reverts(stabilisation.InvalidAmount):
        stabilisation.borrow(0, sender=deposit_scenario.user)


def test_borrow_to_limit(stabilisation, supply_control, deposit_scenario, chain):
    cdp_before = stabilisation.cdps(deposit_scenario.user)
    with auton_transfer_checker(
        chain, supply_control, deposit_scenario.user, deposit_scenario.borrow_limit
    ) as check:
        receipt = stabilisation.borrow(
            deposit_scenario.borrow_limit, sender=deposit_scenario.user
        )
        check(receipt)
    cdp_after = stabilisation.cdps(deposit_scenario.user)
    assert cdp_after.timestamp > cdp_before.timestamp
    assert cdp_after.principal == deposit_scenario.borrow_limit
    assert cdp_after.interest == 0
    verify_borrow_event(deposit_scenario, receipt, deposit_scenario.borrow_limit)


def test_borrow_subsequent(stabilisation, supply_control, borrow_scenario, chain):
    cdp_before = stabilisation.cdps(borrow_scenario.user)
    timestamp = chain.pending_timestamp
    debt = stabilisation.debtAmount(borrow_scenario.user, timestamp)
    interest = debt - borrow_scenario.borrow_amount
    amount = int((borrow_scenario.borrow_limit - debt) / 2)
    with auton_transfer_checker(
        chain, supply_control, borrow_scenario.user, amount
    ) as check:
        receipt = stabilisation.borrow(amount, sender=borrow_scenario.user)
        check(receipt)
    cdp_after = stabilisation.cdps(borrow_scenario.user)
    assert cdp_after.timestamp > cdp_before.timestamp
    assert cdp_after.principal == cdp_before.principal + amount
    assert cdp_after.interest == interest


def test_borrow_minimum(stabilisation, supply_control, deposit_scenario, chain):
    min_debt = stabilisation.config().minDebtRequirement
    with auton_transfer_checker(
        chain, supply_control, deposit_scenario.user, min_debt
    ) as check:
        receipt = stabilisation.borrow(min_debt, sender=deposit_scenario.user)
        check(receipt)
    assert stabilisation.cdps(deposit_scenario.user).principal == min_debt
    verify_borrow_event(deposit_scenario, receipt, min_debt)


def test_borrow_too_little(stabilisation, deposit_scenario):
    min_debt = stabilisation.config().minDebtRequirement
    with ape.reverts(stabilisation.InvalidDebtPosition):
        stabilisation.borrow(min_debt - 1, sender=deposit_scenario.user)


def test_borrow_liquidatable(stabilisation, deposit_scenario, users):
    stabilisation.borrow(
        deposit_scenario.borrow_limit - 1, sender=deposit_scenario.user
    )
    assert not stabilisation.isLiquidatable(deposit_scenario.user)
    mcr = stabilisation.config().minCollateralisationRatio
    stabilisation.setMinCollateralisationRatio(
        mcr + SCALING_FACTOR, sender=users.operator
    )
    stabilisation.setLiquidationRatio(mcr, sender=users.operator)
    assert stabilisation.isLiquidatable(deposit_scenario.user)
    with ape.reverts(stabilisation.Liquidatable):
        stabilisation.borrow(1, sender=deposit_scenario.user)


def test_borrow_over_limit(stabilisation, deposit_scenario):
    with ape.reverts(stabilisation.InsufficientCollateral):
        stabilisation.borrow(
            1 + deposit_scenario.borrow_limit, sender=deposit_scenario.user
        )


# ┌───────┐
# │ Repay │
# └───────┘


def verify_repay_event(stabilisation, receipt, account, payment, interest, surplus=0):
    events = list(receipt.events)
    repay_event = events.pop()
    assert repay_event.account == account
    assert repay_event.amount == payment
    if interest > 0:
        burn_event = events.pop()
        assert burn_event.amount == payment - interest - surplus


def test_repay_zero(stabilisation, borrow_scenario):
    with ape.reverts(stabilisation.ZeroValue):
        stabilisation.repay(value=0, sender=borrow_scenario.user)


def test_repay_invalid_position(stabilisation, borrow_scenario, chain):
    timestamp = chain.pending_timestamp
    too_much = (
        1
        + stabilisation.debtAmount(borrow_scenario.user, timestamp)
        - stabilisation.config().minDebtRequirement
    )
    with ape.reverts(stabilisation.InvalidDebtPosition):
        stabilisation.repay(value=too_much, sender=borrow_scenario.user)


def test_repay_to_minimum_debt(stabilisation, borrow_scenario, chain):
    timestamp = chain.pending_timestamp
    debt = stabilisation.debtAmount(borrow_scenario.user, timestamp)
    interest = debt - borrow_scenario.borrow_amount
    payment = debt - stabilisation.config().minDebtRequirement
    receipt = stabilisation.repay(value=payment, sender=borrow_scenario.user)
    cdp = stabilisation.cdps(borrow_scenario.user)
    assert cdp.interest == 0
    assert cdp.principal == stabilisation.config().minDebtRequirement
    verify_repay_event(stabilisation, receipt, borrow_scenario.user, payment, interest)


def test_repay_interest(stabilisation, borrow_scenario, chain):
    timestamp = chain.pending_timestamp
    interest = (
        stabilisation.debtAmount(borrow_scenario.user, timestamp)
        - borrow_scenario.borrow_amount
    )
    receipt = stabilisation.repay(value=interest, sender=borrow_scenario.user)
    cdp = stabilisation.cdps(borrow_scenario.user)
    assert cdp.interest == 0
    assert cdp.principal == borrow_scenario.borrow_amount
    verify_repay_event(stabilisation, receipt, borrow_scenario.user, interest, 0)


def test_repay_full(stabilisation, borrow_scenario, chain):
    timestamp = chain.pending_timestamp
    debt = stabilisation.debtAmount(borrow_scenario.user, timestamp)
    interest = debt - borrow_scenario.borrow_amount
    receipt = stabilisation.repay(value=debt, sender=borrow_scenario.user)
    cdp = stabilisation.cdps(borrow_scenario.user)
    assert cdp.interest == 0
    assert cdp.principal == 0
    verify_repay_event(stabilisation, receipt, borrow_scenario.user, debt, interest)


def test_repay_surplus(stabilisation, borrow_scenario, chain):
    surplus = 1
    timestamp = chain.pending_timestamp
    debt = stabilisation.debtAmount(borrow_scenario.user, timestamp)
    interest = debt - borrow_scenario.borrow_amount
    user_balance = borrow_scenario.user.balance
    receipt = stabilisation.repay(value=surplus + debt, sender=borrow_scenario.user)
    cdp = stabilisation.cdps(borrow_scenario.user)
    assert cdp.interest == 0
    assert cdp.principal == 0
    if chain.provider.name == "hardhat":
        expected_balance = user_balance - debt
    else:
        expected_balance = user_balance - debt - receipt.total_fees_paid
    assert borrow_scenario.user.balance == expected_balance  # - surplus + surplus
    verify_repay_event(
        stabilisation, receipt, borrow_scenario.user, surplus + debt, interest, surplus
    )


# ┌──────────────┐
# │ Calculations │
# └──────────────┘


def test_borrow_limit(stabilisation):
    tests = [
        (Decimal(100e18), Decimal("1.2e18"), Decimal("1.5e18"), 80000000000000000000),
        (Decimal(100e18), Decimal("0.8e18"), Decimal("1.5e18"), 53333333333333333333),
        (Decimal(100e18), Decimal("1.2e18"), Decimal("1.2e18"), 100000000000000000000),
        (Decimal(100e18), Decimal("0.8e18"), Decimal("1.2e18"), 66666666666666666666),
    ]
    redemption_price = int(1e18)
    for collateral, price, mcr, expected in tests:
        result = stabilisation.borrowLimit(
            int(collateral), int(price), redemption_price, int(mcr)
        )
        assert result == expected
        calculated = quantize(collateral * price / mcr)
        assert result == calculated


def test_minimum_collateral(stabilisation):
    tests = [
        (Decimal(80e18), Decimal("1.2e18"), Decimal("1.5e18"), 100e18),
        (
            Decimal(53333333333333333333),
            Decimal("0.8e18"),
            Decimal("1.5e18"),
            99999999999999999999,
        ),
        (Decimal(100e18), Decimal("1.2e18"), Decimal("1.2e18"), 100e18),
        (
            Decimal(66666666666666666666),
            Decimal("0.8e18"),
            Decimal("1.2e18"),
            99999999999999999999,
        ),
    ]
    for principal, price, mcr, expected in tests:
        result = stabilisation.minimumCollateral(int(principal), int(price), int(mcr))
        assert result == expected
        calculated = quantize(principal * mcr / price)
        assert result == calculated


def test_interest_due(stabilisation):
    tests = [
        (Decimal(100e18), Decimal(0.05e18), 0, 2628000, 417535929111852800),
        (Decimal(100e18), Decimal(0.05e18), 2628000, 2628000, 0),
    ]

    for principal, rate, tstart, tend, expected in tests:
        result = stabilisation.interestDue(int(principal), int(rate), tstart, tend)
        assert result == expected
        t = quantize(SCALING_FACTOR * Decimal(tend - tstart) / SECONDS_IN_YEAR)
        rt = quantize(SCALING_FACTOR * rate / SCALING_FACTOR * t / SCALING_FACTOR)
        exp = quantize(SCALING_FACTOR * PRB_MATH_E ** (rt / SCALING_FACTOR))
        calculated = quantize(principal / SCALING_FACTOR * (exp - SCALING_FACTOR))
        assert result == calculated


def test_collateral_price_is_scaled(
    project,
    users,
    collateral_token,
    basic_config,
    oracle_factory,
    supply_control,
):
    oracle = oracle_factory(PRICE_NTN)
    contract = project.Stabilisation.deploy(
        dataclasses.asdict(basic_config),
        users.operator,
        oracle.address,
        supply_control,
        collateral_token,
        sender=users.deployer,
    )
    assert contract.collateralPrice() == PRICE_NTN_18D


# ┌──────────────┐
# │ Test Helpers │
# └──────────────┘


@dataclasses.dataclass
class Accounts:
    deployer: ape.api.TestAccountAPI
    operator: ape.api.TestAccountAPI
    faucet: ape.api.TestAccountAPI
    account1: ape.api.TestAccountAPI
    account2: ape.api.TestAccountAPI
    autonity: ape.api.TestAccountAPI
    # Oracle
    voter: ape.api.TestAccountAPI


@dataclasses.dataclass
class Config:
    minDebtRequirement: int
    minCollateralisationRatio: int
    borrowInterestRate: int
    liquidationRatio: int
    redemptionPrice: int


@contextmanager
def auton_transfer_checker(chain, from_account, to_account, amount):
    balance_from, balance_to = from_account.balance, to_account.balance
    checked = False

    def checker(receipt):
        nonlocal checked
        if chain.provider.name == "hardhat":
            assert to_account.balance == balance_to + amount
        else:
            assert to_account.balance == balance_to + amount - receipt.total_fees_paid
        assert from_account.balance == balance_from - amount
        checked = True

    yield checker
    assert checked is True  # make sure check() was called


@contextmanager
def check_token_transfer(token, from_account, to_account, amount):
    balance_from, balance_to = token.balanceOf(from_account), token.balanceOf(
        to_account
    )
    yield
    assert token.balanceOf(from_account) == balance_from - amount
    assert token.balanceOf(to_account) == balance_to + amount


@functools.lru_cache
def erc20_contract():
    ape.project.config_manager.contracts_folder = ERC20_CONTRACT.parent.parent.parent
    contracts = ape.compilers.compile([ERC20_CONTRACT])
    return ape.contracts.base.ContractContainer(contracts["ERC20Basic"])


# Solidity arithmetic helper
def quantize(dec):
    return dec.quantize(1, rounding=ROUND_FLOOR)
