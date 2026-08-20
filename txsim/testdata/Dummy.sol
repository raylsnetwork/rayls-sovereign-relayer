// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.10;

contract DummyContract {
    error CustomError(string);

    function hitGenericError() public{
        revert("Hit generic error revert!");
    }

    function hitCustomError() public {
        revert CustomError("Hit custom error revert!");
    }
}
