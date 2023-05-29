// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

contract Records {
    mapping(string => bool) public cidSet;

    function addCID(string memory cid) public {
        cidSet[cid] = true;
    }

    function getCID() public pure returns (string memory) {
        return "test";
     }
}