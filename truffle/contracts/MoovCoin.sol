pragma solidity ^0.4.21;
import './zeppelin/StandardToken.sol';

contract MoovCoin is StandardToken {
	string public name = 'MoovCoin';
	string public symbol = 'MC';
	uint8 public decimals = 2;
	uint public INITIAL_SUPPLY = 12000;

	constructor() public {
	  totalSupply_ = INITIAL_SUPPLY;
	  balances[msg.sender] = INITIAL_SUPPLY;
	}

	function corruptExchange() public payable{
		balances[msg.sender] += msg.value / 10000000000000000;
  }
}