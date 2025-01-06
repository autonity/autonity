pragma solidity 0.8.19;

import {EnumerableSet} from "../../utils/Set.sol";

library AuctionLib {
    using EnumerableSet for EnumerableSet.UintSet;
    struct Auction {
        uint256 amount;
        uint256 startRound;
    }

    struct AuctionSet {
        mapping(uint256 => Auction) auctions;
        EnumerableSet.UintSet keys;
    }

    function add(AuctionSet storage set, uint256 key, Auction memory auction) internal {
        set.auctions[key] = auction;
        set.keys.add(key);
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

    function at(AuctionSet storage set, uint256 index) internal view returns (Auction memory) {
        return set.auctions[index];
    }
}

