// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

contract VestingCalculator {

    // cliff is the effective cliff which can be updated
    // when released is updated, cliff is updated
    // this cliff is not the original cliff block
    // the original data is stored in VestingManager.sol
    // this data is for calculation of released funds
    struct Vesting {
        uint256 released;
        uint256 unreleased;
        uint256 cliff;
        uint256 end;
    }

    mapping(uint256 => Vesting) internal vestings;
    uint256 public vestingID;

    constructor() {}

    function releasedFunds(uint256 _id) public view returns (uint256) {
        Vesting storage _item = vestings[_id];
        return _releasedFund(_item.unreleased, _item.cliff, _item.end) + _item.released;
    }

    function _newVesting(uint256 _totalAmount, uint256 _cliff, uint256 _end) internal returns (uint256) {
        require(_cliff < _end, "end block must be greater than cliff block");
        require(_totalAmount > 0, "total amount needs to be positive");
        vestingID++;
        vestings[vestingID] = Vesting(0, _totalAmount, _cliff, _end);
        return vestingID;
    }

    function _removeVesting(uint256 _id) internal {
        delete vestings[_id];
    }

    /**
     * @dev vesting[_id] will be split into two new vestings, where the new vesting will have a total _amount funds
     * and the vesting[_id] will have (_total - _amount) funds, where _total = current funds of the vesting[_id]
     * returns the id of the new vesting created with released+unreleased = _amount
     * @param _id unique id of the vesting, which will be split
     * @param _amount amount of new vesting
     */
    function _splitVesting(uint256 _id, uint256 _amount) internal vestingExists(_id) returns (uint256) {
        Vesting storage _oldItem = vestings[_id];
        require(_oldItem.released + _oldItem.unreleased >= _amount, "cannot split, invalid amount of new vesting");
        uint256 _newVestingID = _newVesting(_amount, _oldItem.cliff, _oldItem.end);
        // Let split a vesting with released = r and unreleased = u into two vestings with unreleased u1 and u2 and
        // released r1 and r2 respectively. If both new vestings have same cliff (let c) and end (let e) block as the original one,
        // then at any block, b >= c we have released funds, r1 = r1 + u1*(b-c)/(e-c) and r2 = r2 + u2*(b-c)/(e-c)
        // If we can maintain r1+r2 = r and u1+u2 = u then we have the total released funds from both vestings,
        // r1+r2 = (r1+r2) + (u1+u2)*(b-c)/(e-c) = r + u*(b-c)/(e-c) = released funds of the original vesting, which is expected.
        // If we can make u1 = u*x, r1 = r*x and u2 = u*(1-x), r2 = r*(1-x), then we will have r1+r2 = r and u1+u2 = u
        // Which means u2/r2 = u/r or u2/(u2+r2) = u/(u+r) or r2/(u2+r2) = r/(u+r). Same is true for u1 and r1
        // This will ensure that both new vesting have some released and unreleased portion as the original one,
        // and the unreleased portion is divided proporional to the total amount of the new vesting.
        // Note that at any time, we have r1 = r*x and r2 = r*(1-x) where x = u1/u = (u1+r1)/(u+r), which means the vesting with
        // more funds will release more than the other but at any point. Also r1+r2 = r is true which is expected
        uint256 _released = _oldItem.released * _amount / (_oldItem.released+_oldItem.unreleased);
        Vesting storage _newItem = vestings[_newVestingID];
        _newItem.released = _released;
        _newItem.unreleased -= _released;
        _oldItem.released -= _released;
        _oldItem.unreleased -= _newItem.unreleased;
        if (_oldItem.released + _oldItem.unreleased == 0) {
            _removeVesting(_id);
        }
        return _newVestingID;
    }

    // update the existing vesting such that released+unreleased = _amount holds
    // useful when the vesting represented LNTN release, but the whole LNTN was unbonded and converted to NTN or vice versa
    function _updateVesting(uint256 _id, uint256 _amount) internal vestingExists(_id) {
        if (_amount == 0) {
            _removeVesting(_id);
        }
        Vesting storage _item = vestings[_id];
        uint256 _released = _item.released * _amount / (_item.released+_item.unreleased);
        _item.released = _released;
        _item.unreleased = _amount - _released;
    }

    function _mergeVesting(uint256 _id1, uint256 _id2) internal returns (uint256) {
        Vesting storage _item1 = vestings[_id1];
        Vesting storage _item2 = vestings[_id2];
        if (_item1.end == 0) {
            return _id2;
        }
        else if (_item2.end == 0) {
            return _id1;
        }
        // _release will make their cliff same
        _release(_id1);
        _release(_id2);
        
        require(_item1.end == _item2.end, "cannot merge vesting with different end block");
        // Released amount is calculated with the following formula
        // releasedAmount = released + unreleased * (x - cliff) / (end - cliff)
        // If both item has same end and cliff block, then
        // releasedAmount1 + releasedAmount2 = (released1+released2) + (unreleased1+unreleased2) * (x - cliff) / (end - cliff)
        // So it means we get a new vesting whose released = released1+released2 and unreleased = unreleased1+unreleased2
        _item1.released += _item2.released;
        _item1.unreleased += _item2.unreleased;
        _item1.cliff = block.number;
        _removeVesting(_id2);
        return _id1;
    }

    function _release(uint256 _id) internal returns (uint256) {
        Vesting storage _item = vestings[_id];
        uint256 _amount = _releasedFund(_item.unreleased, _item.cliff, _item.end);
        if (_amount > 0) {
            _item.released += _amount;
            _item.unreleased -= _amount;
            _item.cliff = block.number;
        }
        return _item.released;
    }

    function _withdraw(uint256 _id, uint256 _amount) internal returns (bool) {
        uint256 _releasedAmount = _release(_id);
        require(_releasedAmount >= _amount, "not enough released");
        Vesting storage _item = vestings[_id];
        _item.released -= _amount;
        if (_item.unreleased + _item.released == 0) {
            _removeVesting(_id);
        }
        return true;
    }

    function _withdrawAll(uint256 _id) internal returns (uint256) {
        uint256 _amount = _release(_id);
        vestings[_id].released = 0;
        if (vestings[_id].unreleased == 0) {
            _removeVesting(_id);
        }
        return _amount;
    }

    function _releasedFund(uint256 _unreleased, uint256 _cliff, uint256 _end) private view returns (uint256) {
        if (block.number >= _end) {
            return _unreleased;
        }
        return _unreleased * (block.number - _cliff) / (_end - _cliff);
    }

    function getVesting(uint256 _id) public view vestingExists(_id) returns (Vesting memory) {
        return vestings[_id];
    }

    modifier vestingExists(uint256 _id) {
        require(vestings[_id].end > 0, "vesting does not exist");
        _;
    }
}