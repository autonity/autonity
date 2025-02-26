// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;


struct UintQueue {
    uint256 topIndex;
    uint256[] array;
}

// this library is not for general queue operation
// a general-purpose queue library may be implemented later
library UintQueueLib {
    function enqueue(UintQueue storage _queue, uint256 _item) internal {
        _queue.array.push(_item);
    }

    function dequeue(UintQueue storage _queue, uint256 _deleteCount) internal {
        if (_deleteCount == 0) {
            return;
        }
        uint256 _topIndex = _queue.topIndex;
        // length of the queue is `_queue.array.length - _queue.topIndex`
        require(_deleteCount <= _queue.array.length - _topIndex, "not enough elements in the queue");
        for ( ; _deleteCount > 0; _deleteCount--) {
            delete _queue.array[_topIndex];
            _topIndex++;
        }
        _queue.topIndex = _topIndex;

        // reset everything if no element left
        if (_topIndex == _queue.array.length) {
            _queue.topIndex = 0;
            // use assembly because at this point the array maybe too large but with zero values
            uint256[] storage _array = _queue.array;
            assembly {
                // The length of `_array` is stored at `_array.slot`.
                // By setting 0 at `_array.slot` we effectively delete the array without iterating.
                // This is not a good idea in general as the elements are not be set to 0 and will
                // consume space. But in this particular case, all the elements are already set to 0.
                sstore(_array.slot, 0)
            }
        }
    }
}