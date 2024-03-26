// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

contract VestingCalculator {

    // start is the effective start which can be updated when unlocked is updated.
    // this start is not the original start block.
    // the original data is stored in VestingManager.sol
    // this data is for calculation of unlocked funds
    struct Vesting {
        uint256 unlocked;
        uint256 locked;
        uint256 start;
        uint256 cliff;
        uint256 end;
    }

    // Stores all the Vesting object. Each object represent a vesting of some schedule or
    // some liquidNTN of a pair (id, v) where id = unique schedule id and v = validator address
    mapping(uint256 => Vesting) private vestings;
    uint256 private vestingID;

    constructor() {}

    function releasedFunds(uint256 _id) internal view returns (uint256) {
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

    /** @dev calculates the amount of total unlocked amount and updates the unlocked and locked amount.
     * If current block = x > cliff, then it also updates start = cliff = x. start is updated for the following reasons:
     * 
     *      1.  The release equation can be written: u' = (t+w) * (x-s) / (e-s) where s = start, e = end, x = current block
     *          and t = total amount remaining and w = withdrawn amount, u' = total unlocked amount at current block.
     *          But keeping w in equation creates some issues in some particular scenario. If w = 0 we can write
     *          u' = t * (x-s) / (e-s). Consider t = (u+l), where u = unlocked and l = locked amount at some block, b <= x.
     *          New unlocked = new_u = u'-u, so u and l can be updated: u += new_u and l -= new_u
     *          If w remains 0, u+l remains constant as u increases and l decreases by the same amount.
     *          So the equation becomes u' = (l+u) * (x-s) / (e-s) only if w remains 0.
     *          As funds can only be withdrawn from unlocked amount, if we can remove u from this
     *          equation, we don't need to worry about w.
     * 
     *      2.  By rewriting the equation u' = l * (x-s) / (e-s), it works for the first time when u = 0 and w = 0. After that u and l
     *          is updated as u += new_u and l -= new_u. Now in the above equation, the slope is changed and the equation is incorrect with u > 0.
     *          But this can be mitigated by updaing s = x. Cliff is updated for the sake of the logic that start <= Cliff should be true.
     * 
     *      How updaing s = x let us remove u from the equation is explained here:
     *      Note: The following math can be easily understood and visulaized with a graph instead
     *          Lets denote the following variables:
     *          u' = unlocked funds at block x'
     *          u'' = unlocked funds at block x'' where x'' > x' and so u'' > u'
     *          l = total funds = locked amount before start or cliff block
     *          we can write u' = l * (x'-s) / (e-s) and u'' = l * (x''-s) / (e-s)
     *          which can be rewritten l / (e-s) = u' / (x'-s) and l / (e-s) = u'' / (x''-s)
     *          Applying the following formula in the above equations: [Can't remember the name; TODO: put the name of the following formula]
     *          if a/b = c/d then a/b = c/d = (ma+nc)/(mb+nd) is true where m != 0 and n != 0
     *          So we can write l / (e-s) = u' / (x'-s) = (l-u') / (e-x') and l / (e-s) = u' / (x'-s) = u'' / (x''-s) = (u''-u') / (x''-x')
     *          So we have (l-u') / (e-x') = (u''-u') / (x''-x') or (u''-u') = (l-u') * (x''-x') / (e-x')
     *          Here we can write, after block x', locked amount = l' = l-u' and unlocked amount = u'
     *          And after block x'', new unlocked amount = u''-u' = l' * (x''-x') / (e-x'), x' is sitting where we should put start block.
     *          So we can calculated new unlocked amount with the formula, new_u = l * (x-s) / (e-s), where e = end and x = current block,
     *          s = start block or the last block when unlocked amount was updated, l = locked amount after block s.
     *          And then we can update locked amount, l -= new_u, unlocked amount, u += new_u and start block, s = x
     * 
     * @param _id unique id of the vesting
     */
    function _release(uint256 _id) internal returns (uint256) {
        Vesting storage _item = vestings[_id];
        uint256 _newUnlockedamount = _releasedFund(_item.locked, _item.start, _item.cliff, _item.end);
        if (_newUnlockedamount > 0) {
            _item.unlocked += _newUnlockedamount;
            _item.locked -= _newUnlockedamount;
        }
        // end > 0 means it exists, otherwise the vesting does not exist and everying is set to 0
        if (_item.end > 0 && block.number > _item.cliff) {
            _item.start = block.number;
            _item.cliff = block.number;
        }
        return _item.unlocked;
    }

    function _removeVesting(uint256 _id) internal {
        delete vestings[_id];
    }

    /** @dev creates a new vesting for delegation of _amount NTN or LNTN from existing vesting with id = _id.
     * Two scenarios can happen:
     * 
     *      1.  A new bonding of _amount NTN is requested from some schedule. In that case vesting with id = _id represents that schedule.
     *          NTN of _amount will be deducted from vesting with id = _id and the new created vesting will represent the new delegation.
     *          In other words, the new vesting will represent the LNTN created from this bonding request. In case LNTN:NTN != 1,
     *          the amount of this new vesting will be updated when we get the LNTN amount after bonding is applied (see _updateVesting).
     * 
     *      2.  A new unbonding of _amount LNTN is applied. In that case vesting with id = _id represent the liquid amount from  which the unbonding
     *          was requested. So LNTN of _amount will be deducted from vesting with id = _id and the new created vesting will represent the new
     *          unbonding of delegation. In other words the new vesting will represent the NTN created from this unbonding request. In case
     *          LNTN:NTN != 1, the amount of this new vesting will be updated when we get the NTN amount when unbonding is released (see _updateVesting).
     * 
     * returns the id of the new vesting created with unlocked+locked = _amount
     * @param _id unique id of the existing vesting
     * @param _amount amount of new vesting
     */
    function _createDelegation(uint256 _id, uint256 _amount) internal vestingExists(_id) returns (uint256) {
        Vesting storage _oldItem = vestings[_id];
        require(_oldItem.unlocked + _oldItem.locked >= _amount, "cannot split, invalid amount of new vesting");
        uint256 _newVestingID = _newVesting(_amount, _oldItem.start, _oldItem.cliff, _oldItem.end);
        // Let split a vesting with unlocked = u and locked = l into two vestings with locked l1 and l2 and
        // unlocked u1 and u2 respectively. If both new vestings have same start (let s), cliff and end (let e) block as the original one,
        // then at any block, b >= c we have unlocked funds, u1' = u1 + l1*(b-s)/(e-s) and u2' = u2 + l2*(b-s)/(e-s)
        // If we can maintain u1+u2 = u and l1+l2 = l then we have the total unlocked funds from both vestings,
        // u1'+u2' = (u1+u2) + (l1+l2)*(b-s)/(e-s) = u + l*(b-s)/(e-s) = unlocked funds of the original vesting, which is expected.
        // If we can make l1 = l*x, u1 = u*x and l2 = l*(1-x), u2 = u*(1-x), then we will have l1+l2 = l and u1+u2 = u
        // Which means u2/l2 = u/l or u2/(u2+l2) = u/(u+l) or l2/(u2+l2) = l/(u+l). Same is true for u1 and l1
        // This will ensure that both new vesting have some unlocked and locked portion as the original one,
        // and the locked portion is divided proporional to the total amount of the new vesting.
        // Note that at any time, we have u1 = u*x and u2 = u*(1-x) where x = u1/u = l1/l = (u1+l1)/(u+l), which means the vesting with
        // more funds will release more than the other one at any point. Also u1+u2 = u is true which is expected
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

    /** @dev merges the two vestings with id _id1 and _id2 and creates a new vesting so that the total amount is the sum of 
     * the total amount of these two vestings. And also at any block = x the total unlocked amount of the new vesting is
     * sum of unlocked amount of these two vestings at block = x. Returns the id of the new vesting
     * @param _id1 unique id of the first vesting
     * @param _id2 unique id of the second vesting
     */
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
        
        require(_item1.end == _item2.end, "end block not equal");
        require(_item1.cliff == _item2.cliff, "cliff block not equal");
        require(_item1.start == _item2.start, "start block not equal");
        // Released amount is calculated with the following formula
        // unlocked' = unlocked + locked * (x - start) / (end - start)
        // If both item has same end and start block, then
        // unlocked1' + unlocked2' = (unlocked1+unlocked2) + (locked1+locked2) * (x - start) / (end - start)
        // Weh can say, totalUnlocked = unlocked1' + unlocked2'
        // It is equivalent to having a new vesting whose unlocked = unlocked1+unlocked2 and locked = locked1+locked2
        // So we have totalUnlocked = unlocked + locked * (x - start) / (end - start) from the new vesting
        _item1.unlocked += _item2.unlocked;
        _item1.locked += _item2.locked;
        _removeVesting(_id2);
        return _id1;
    }

    function _decreaseUnlocked(uint256 _id, uint256 _newUnlockedamount) internal returns (bool) {
        uint256 _unlockedAmount = _release(_id);
        require(_unlockedAmount >= _newUnlockedamount, "not enough unlocked tokens");
        Vesting storage _item = vestings[_id];
        _item.unlocked -= _newUnlockedamount;
        if (_item.locked + _item.unlocked == 0) {
            _removeVesting(_id);
        }
        return true;
    }

    function _decreaseAllUnlocked(uint256 _id) internal returns (uint256) {
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

    function getVesting(uint256 _id) internal view vestingExists(_id) returns (Vesting memory) {
        return vestings[_id];
    }

    modifier vestingExists(uint256 _id) {
        require(vestings[_id].end > 0, "vesting does not exist");
        _;
    }
}