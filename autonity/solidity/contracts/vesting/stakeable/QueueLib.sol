// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

struct PendingStakingRequest {
    uint256 epochID;
    address validator;
    uint256 amount;
    uint256 requestID;
}

struct StakingRequestQueue {
    uint256 topIndex;
    PendingStakingRequest[] array;
}

// this library is not for general queue operation
// a general-purpose queue library may be implemented later
library QueueLib {
    function enqueue(StakingRequestQueue storage _queue, PendingStakingRequest memory _item) internal {
        _queue.array.push(_item);
    }

    function dequeue(StakingRequestQueue storage _queue, uint256 _deleteCount) internal {
        uint256 _topIndex = _queue.topIndex;
        // length of the queue is `_queue.array.length - _queue.topIndex`
        require(_deleteCount <= _queue.array.length - _topIndex, "not enough elements in the queue");
        PendingStakingRequest storage _item;
        for ( ; _deleteCount > 0; _deleteCount--) {
            _item = _queue.array[_topIndex];
            _item.amount = 0;
            _item.epochID = 0;
            _item.requestID = 0;
            _item.validator = address(0);
            _topIndex++;
        }
        _queue.topIndex = _topIndex;
    }
}