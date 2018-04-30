pragma solidity ^0.4.18;

import './MoovCoin.sol';

contract MoovRideManager {
  // MoovCoin Contract variable to conduct financial transactions. 
  MoovCoin moovCoin;

  // Allowed statuses for each ride
  enum RideStatus {AVAILABLE, REQUESTING, INPROGRESS}
  
  // Struct defining attributes for each ride
  struct ride {
    string from;
    string to;
    uint amount;
    RideStatus rideStatus;
    address carAddress;
  }

  // Mapping from each rider to a ride
  mapping(address => ride) rides;

  // Array to keep track of rides that are available to be picked up by cars
  address[] public availableRides;

  // Events that will be fired on changes.
  event NewRideRequest(address rider, string from, string to, uint amount);
  event RideAccepted(address rider, address car);
  event RideFinished(address rider, address car);

  /// Create a Moov Ride Manager pointing to current coin address
  function MoovRideManager(address moovCoinAddress) public {
    moovCoin = MoovCoin(moovCoinAddress);
  }

  /// Create a new Ride Request  FROM:Rider
  function newRideRequest(string from, string to, uint amount) public returns (bool){
    address riderAddress = msg.sender;
    require (rides[riderAddress].rideStatus == RideStatus.AVAILABLE);
    require(moovCoin.transferFrom(riderAddress, this, amount));  
    rides[riderAddress].from = from;
    rides[riderAddress].to = to;
    rides[riderAddress].amount = amount;
    rides[riderAddress].rideStatus = RideStatus.REQUESTING;
    addToAvailableRides(riderAddress);
    NewRideRequest(riderAddress, from, to, amount);
  }

  //// cancelRideRequest  FROM:Rider
  function cancelRideRequest() public {
    address riderAddress = msg.sender;
    require (rides[riderAddress].rideStatus == RideStatus.REQUESTING);
    require(moovCoin.transfer(riderAddress, rides[riderAddress].amount));
    rides[riderAddress].rideStatus = RideStatus.AVAILABLE;
    removeFromAvailableRides(riderAddress);
  }

  /// Accept a ride request  FROM:Car
  function acceptRideRequest(address chosenRider) public {
    address carAddress = msg.sender;
    require (rides[chosenRider].rideStatus == RideStatus.REQUESTING);
    rides[chosenRider].rideStatus = RideStatus.INPROGRESS;
    rides[chosenRider].carAddress = carAddress;
    removeFromAvailableRides(chosenRider);
    RideAccepted(chosenRider, carAddress);
  }

  /// Finish ride and transfer Moov coins to car   FROM:Rider
  function finishRide() public{
    address riderAddress = msg.sender;
    require (rides[riderAddress].rideStatus == RideStatus.INPROGRESS);
    require(moovCoin.transfer(rides[riderAddress].carAddress, rides[riderAddress].amount));
    rides[riderAddress].rideStatus = RideStatus.AVAILABLE;
    rides[riderAddress].amount = 0;
    RideFinished(msg.sender, rides[riderAddress].carAddress);
  }

  function getMoovCoinAddress() public view returns(address) {
    return address(moovCoin);
  }

  function getBalance(address addr) public view returns(uint) {
    return rides[msg.sender].amount;
  }

  function addToAvailableRides(address newAddress) private{
    availableRides.push(newAddress);
  }

  function removeFromAvailableRides(address newAddress) private{
    for(uint i=0; i < availableRides.length; i++){
      if (availableRides[i] == newAddress){
        // foung address so push the last address to the found address's position
        availableRides[i] = availableRides[availableRides.length-1];
        // delete last entry in the array
        delete availableRides[availableRides.length-1];
        availableRides.length--;
        break;
      }
    }
  }

}


  