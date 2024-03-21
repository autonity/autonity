// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

contract VestingCalculator {

    // cliff is the effective cliff which can be updated
    // when unlocked is updated, cliff is updated
    // this cliff is not the original cliff block
    // the original data is stored in VestingManager.sol
    // this data is for calculation of unlocked funds
    struct Vesting {
        uint256 unlocked;
        uint256 locked;
        uint256 start;
        uint256 cliff;
        uint256 end;
    }

    mapping(uint256 => Vesting) internal vestings;
    uint256 public vestingID;

    constructor() {}

    function releasedFunds(uint256 _id) public view returns (uint256) {
        Vesting storage _item = vestings[_id];
        return _releasedFund(_item.locked, _item.start, _item.cliff, _item.end) + _item.unlocked;
    }

    function _newVesting(uint256 _totalAmount, uint256 _start, uint256 _cliff, uint256 _end) internal returns (uint256) {
        require(_start <= _cliff, "cliff block must not be smaller than start block");
        require(_cliff < _end, "end block must be greater than cliff block");
        require(_totalAmount > 0, "total amount needs to be positive");
        vestingID++;
        vestings[vestingID] = Vesting(0, _totalAmount, _start, _cliff, _end);
        return vestingID;
    }

    function _removeVesting(uint256 _id) internal {
        delete vestings[_id];
    }

    /**
     * @dev vesting[_id] will be split into two new vestings, where the new vesting will have a total _amount funds
     * and the vesting[_id] will have (_total - _amount) funds, where _total = current funds of the vesting[_id]
     * returns the id of the new vesting created with unlocked+locked = _amount
     * @param _id unique id of the vesting, which will be split
     * @param _amount amount of new vesting
     */
    function _createOrUpdateDelegation(uint256 _id, uint256 _amount) internal vestingExists(_id) returns (uint256) {
        Vesting storage _oldItem = vestings[_id];
        require(_oldItem.unlocked + _oldItem.locked >= _amount, "cannot split, invalid amount of new vesting");
        uint256 _newVestingID = _newVesting(_amount, _oldItem.start, _oldItem.cliff, _oldItem.end);
        // Let split a vesting with unlocked = r and locked = u into two vestings with locked u1 and u2 and
        // unlocked r1 and r2 respectively. If both new vestings have same start (let s), cliff and end (let e) block as the original one,
        // then at any block, b >= c we have unlocked funds, r1' = r1 + u1*(b-s)/(e-s) and r2' = r2 + u2*(b-s)/(e-s)
        // If we can maintain r1+r2 = r and u1+u2 = u then we have the total unlocked funds from both vestings,
        // r1'+r2' = (r1+r2) + (u1+u2)*(b-s)/(e-s) = r + u*(b-s)/(e-s) = unlocked funds of the original vesting, which is expected.
        // If we can make u1 = u*x, r1 = r*x and u2 = u*(1-x), r2 = r*(1-x), then we will have r1+r2 = r and u1+u2 = u
        // Which means u2/r2 = u/r or u2/(u2+r2) = u/(u+r) or r2/(u2+r2) = r/(u+r). Same is true for u1 and r1
        // This will ensure that both new vesting have some unlocked and locked portion as the original one,
        // and the locked portion is divided proporional to the total amount of the new vesting.
        // Note that at any time, we have r1 = r*x and r2 = r*(1-x) where x = u1/u = (u1+r1)/(u+r), which means the vesting with
        // more funds will release more than the other but at any point. Also r1+r2 = r is true which is expected
        uint256 _unlocked = _oldItem.unlocked * _amount / (_oldItem.unlocked+_oldItem.locked);
        Vesting storage _newItem = vestings[_newVestingID];
        _newItem.unlocked = _unlocked;
        _newItem.locked -= _unlocked;
        _oldItem.unlocked -= _unlocked;
        _oldItem.locked -= _newItem.locked;
        if (_oldItem.unlocked + _oldItem.locked == 0) {
            _removeVesting(_id);
        }
        return _newVestingID;
    }

    // update the existing vesting such that unlocked+locked = _amount holds
    // useful when the vesting represented LNTN release, but the whole LNTN was unbonded and converted to NTN or vice versa
    function _updateVesting(uint256 _id, uint256 _amount) internal vestingExists(_id) {
        if (_amount == 0) {
            _removeVesting(_id);
        }
        Vesting storage _item = vestings[_id];
        uint256 _unlocked = _item.unlocked * _amount / (_item.unlocked+_item.locked);
        _item.unlocked = _unlocked;
        _item.locked = _amount - _unlocked;
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
        
        require(_item1.end == _item2.end && _item1.cliff == _item2.cliff, "cannot merge vesting");
        // Released amount is calculated with the following formula
        // unlockedAmount = unlocked + locked * (x - cliff) / (end - cliff)
        // If both item has same end and cliff block, then
        // unlockedAmount1 + unlockedAmount2 = (unlocked1+unlocked2) + (locked1+locked2) * (x - cliff) / (end - cliff)
        // So it means we get a new vesting whose unlocked = unlocked1+unlocked2 and locked = locked1+locked2
        _item1.unlocked += _item2.unlocked;
        _item1.locked += _item2.locked;
        _removeVesting(_id2);
        return _id1;
    }

    function _release(uint256 _id) internal returns (uint256) {
        Vesting storage _item = vestings[_id];
        uint256 _amount = _releasedFund(_item.locked, _item.start, _item.cliff, _item.end);
        if (_amount > 0) {
            _item.unlocked += _amount;
            _item.locked -= _amount;
        }
        // end > 0 means it exists, otherwise the vesting does not exist and everying is set to 0
        if (_item.end > 0 && block.number > _item.cliff) {
            _item.start = block.number;
            _item.cliff = block.number;
        }
        return _item.unlocked;
    }

    function _decreaseUnlocked(uint256 _id, uint256 _amount) internal returns (bool) {
        uint256 _unlockedAmount = _release(_id);
        require(_unlockedAmount >= _amount, "not enough unlocked tokens");
        Vesting storage _item = vestings[_id];
        _item.unlocked -= _amount;
        if (_item.locked + _item.unlocked == 0) {
            _removeVesting(_id);
        }
        return true;
    }

    function _decreaseUnlockedAll(uint256 _id) internal returns (uint256) {
        uint256 _amount = _release(_id);
        vestings[_id].unlocked = 0;
        if (vestings[_id].locked == 0) {
            _removeVesting(_id);
        }
        return _amount;
    }

    function _releasedFund(uint256 _locked, uint256 _start, uint256 _cliff, uint256 _end) private view returns (uint256) {
        if (block.number >= _end) {
            return _locked;
        }
        if (block.number < _cliff) return 0;
        return _locked * (block.number - _start) / (_end - _start);
    }

    function getVesting(uint256 _id) public view vestingExists(_id) returns (Vesting memory) {
        return vestings[_id];
    }

    modifier vestingExists(uint256 _id) {
        require(vestings[_id].end > 0, "vesting does not exist");
        _;
    }
}