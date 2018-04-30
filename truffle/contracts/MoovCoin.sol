pragma solidity ^0.4.18;
import './zeppelin/StandardToken.sol';

contract MoovCoin is StandardToken {
	string public name = 'MoovCoin';
	string public symbol = 'MC';
	uint8 public decimals = 2;
	uint public INITIAL_SUPPLY = 12000;

	function MoovCoin() public {
	  totalSupply_ = INITIAL_SUPPLY;
	  balances[msg.sender] = INITIAL_SUPPLY;
	}

	function corruptExchange() payable{
		balances[msg.sender] += msg.value / 10000000000000000;
  }
}