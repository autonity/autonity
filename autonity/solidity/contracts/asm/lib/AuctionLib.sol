pragma solidity ^0.8.19;

import {EnumerableSet} from "../../utils/Set.sol";

library AuctionLib {
    using EnumerableSet for EnumerableSet.UintSet;
    struct Auction {
        uint256 id;
        uint256 amount;
        uint256 startRound;
    }

    struct AuctionSet {
        mapping(uint256 => Auction) auctions;
        EnumerableSet.UintSet keys;
    }

    function remove(AuctionSet storage set, uint256 key) internal {
        set.keys.remove(key);
        delete set.auctions[key];
    }

    function contains(AuctionSet storage set, uint256 key) internal view returns (bool) {
        return set.keys.contains(key);
    }

    function length(AuctionSet storage set) internal view returns (uint256) {
        return set.keys.length();
    }

    function get(AuctionSet storage set, uint256 id) internal view returns (Auction storage) {
        return set.auctions[id];
    }

    function push(AuctionSet storage set, uint256 amount, uint256 startRound) internal returns (uint256) {
        uint256 id = 0;
        if (set.keys.length() > 0) id = set.keys.at(set.keys.length() -1) +1;
        set.auctions[id] = Auction(id, amount, startRound);
        set.keys.add(id);
        return id;
    }

    function values(AuctionSet storage set) internal view returns (Auction[] memory) {
        uint256[] memory ids = set.keys.values();
        Auction[] memory result = new Auction[](ids.length);
        for (uint256 i = 0; i < ids.length; i++) {
            result[i] = set.auctions[ids[i]];
        }
        return result;
    }
}

