require('truffle-test-utils').init();
var MRM = artifacts.require("./MoovRideManager.sol");
var MoovCoin = artifacts.require("./MoovCoin.sol");

contract('MoovRideManager', function(accounts) {
  it("should check Moov Coin Address", async () => {
    let moovCoin = await MoovCoin.deployed();
    let mrm = await MRM.deployed();
    let expected = await mrm.moovCoin();
    assert.equal(moovCoin.address, expected, "Moov Coin address is different");
  });


  it("should create new ride request correctly", async () => {
    var rideAmount = 1;
    var from = "Earth";
    var to = "Mars";
    var riderAccount = web3.eth.coinbase;
    let moovCoin = await MoovCoin.deployed();
  	let mrm = await MRM.deployed();
    let riderMoovBalance = await moovCoin.balanceOf(riderAccount);
    let mrmBalance = await moovCoin.balanceOf(mrm.address);
    let riderMRMBalance = await mrm.getBalance(riderAccount);
    await moovCoin.approve(mrm.address, rideAmount, {from:riderAccount});
  	let result = await mrm.newRideRequest(from, to, rideAmount, {from:riderAccount});
  	// Moov Coin Assert Userbalance reduce and MRM balance increase
    let actual = await moovCoin.balanceOf(riderAccount);
    let expected = riderMoovBalance.sub(rideAmount);
    assert.ok(actual.eq(expected), "Rider Moov Coin balance is not "+expected.toString());
    actual = await moovCoin.balanceOf(mrm.address);
    expected = mrmBalance.plus(rideAmount);
    assert.ok(actual.eq(expected), "MRM Moov coin balance is not "+expected.toString());
    // MRM assert User balance increase
    actual = await mrm.getBalance(riderAccount);
  	expected = riderMRMBalance.plus(rideAmount);
  	assert.ok(actual.eq(expected), "Rider Balance in MRM is not "+expected.toString());

    // Assert Ride status
    actual = await mrm.getRideStatus(riderAccount);
    expected = 1;
    assert.ok(actual.eq(expected), "Ride status is not in requesting");

    // Assert destinations
    actual = await mrm.getDestinations(riderAccount)
    let expected_from = "Earth";
    let expected_to = "Mars";
    assert.equal(actual[0], expected_from, "from destination does not match");
    assert.equal(actual[1], expected_to, "to destination does not match");

    // Asset user account got added to available rides
    let availableRides = await mrm.getAvailableRides();
    assert.include(availableRides, riderAccount, "Available ride address does not include user account");
  });

  it("should emit new ride request event correctly", async () => {
    let moovCoin = await MoovCoin.deployed();
    let mrm = await MRM.new(moovCoin.address);
    var rideAmount = 1;
    var from = "Earth";
    var to = "Mars";
    var riderAccount = web3.eth.coinbase;
    await moovCoin.approve(mrm.address, rideAmount, {from:riderAccount});
    let result = await mrm.newRideRequest(from, to, rideAmount, {from:riderAccount});
    assert.web3Event(result, {
      event: 'NewRideRequest',
        args: {
          rider: riderAccount,
          from: "Earth",
          to: "Mars",
          amount: rideAmount 
      }
    }, 'The event is emitted');
  });

  it("should cancel ride request correctly", async () => {
    let moovCoin = await MoovCoin.deployed();
    let mrm = await MRM.new(moovCoin.address);
    var rideAmount = 1;
    var from = "Earth";
    var to = "Mars";
    var riderAccount = web3.eth.coinbase;
    await moovCoin.approve(mrm.address, rideAmount, {from:riderAccount});
    let result = await mrm.newRideRequest(from, to, rideAmount, {from:riderAccount});

    let riderMoovBalance = await moovCoin.balanceOf(riderAccount);
    let mrmBalance = await moovCoin.balanceOf(mrm.address);
    let riderMRMBalance = await mrm.getBalance(riderAccount);
    await mrm.cancelRideRequest();
    // Moov Coin Assert User balance increase and MRM balance decrease
    let actual = await moovCoin.balanceOf(riderAccount);
    let expected = riderMoovBalance.plus(rideAmount);
    assert.ok(actual.eq(expected), "Rider Moov Coin Balance is not "+expected.toString());
    actual = await moovCoin.balanceOf(mrm.address);
    expected = mrmBalance.sub(rideAmount);
    assert.ok(actual.eq(expected), "MRM Moov Coin Balance is not "+expected.toString());
    // MRM assert User account balance decrease
    actual = await mrm.getBalance(riderAccount);
    expected = 0;
    assert.ok(actual.eq(expected), "Rider MRM account Balance is not "+expected.toString());

    // Asset user account got added to available rides
    let availableRides = await mrm.getAvailableRides();
    assert.notInclude(availableRides, riderAccount, "Available ride address does not include user account");

    // Assert Ride status
    actual = await mrm.getRideStatus(riderAccount);
    expected = 0;
    assert.ok(actual.eq(expected), "Ride status is not in available");
  });


  it("should accept correctly", async () => {
    let moovCoin = await MoovCoin.deployed();
    let mrm = await MRM.new(moovCoin.address);
    var rideAmount = 1;
    var from = "Earth";
    var to = "Mars";
    var riderAccount = web3.eth.coinbase;
    var carAccount = web3.eth.accounts[1]
    await moovCoin.approve(mrm.address, rideAmount, {from:riderAccount});
    let result = await mrm.newRideRequest(from, to, rideAmount, {from:riderAccount});
    result = await mrm.acceptRideRequest(riderAccount, {from:carAccount});  

    // Assert Car address
    actual = await mrm.getCarAddress(riderAccount);
    expected = carAccount;
    assert.equal(actual, expected, "Car Address does not match");

    // Asset user account got added to available rides
    let availableRides = await mrm.getAvailableRides();
    assert.notInclude(availableRides, riderAccount, "Available ride address does not include user account");

    // Assert Ride status
    actual = await mrm.getRideStatus(riderAccount);
    expected = 2;
    assert.ok(actual.eq(expected), "Ride status is not in progress");
  });



  it("should emit accept events correctly", async () => {
    let moovCoin = await MoovCoin.deployed();
    let mrm = await MRM.new(moovCoin.address);
    var rideAmount = 1;
    var from = "Earth";
    var to = "Mars";
    var riderAccount = web3.eth.coinbase;
    var carAccount = web3.eth.accounts[1]
    await moovCoin.approve(mrm.address, rideAmount, {from:riderAccount});
    let result = await mrm.newRideRequest(from, to, rideAmount, {from:riderAccount});
    result = await mrm.acceptRideRequest(riderAccount, {from:carAccount});
    assert.web3Event(result, {
      event: 'RideAccepted',
        args: {
          rider: riderAccount,
          car: carAccount 
        }
    }, 'The event is emitted');
  });


  it("should finish correctly", async () => {
    let moovCoin = await MoovCoin.deployed();
    let mrm = await MRM.new(moovCoin.address);
    var rideAmount = 1;
    var from = "Earth";
    var to = "Mars";
    var riderAccount = web3.eth.coinbase;
    var carAccount = web3.eth.accounts[1]
    await moovCoin.approve(mrm.address, rideAmount, {from:riderAccount});
    let result = await mrm.newRideRequest(from, to, rideAmount, {from:riderAccount});
    result = await mrm.acceptRideRequest(riderAccount, {from:carAccount});

    let carMoovBalance = await moovCoin.balanceOf(carAccount);
    let mrmBalance = await moovCoin.balanceOf(mrm.address);
    let riderMRMBalance = await mrm.getBalance(riderAccount);
    await mrm.finishRide({from:riderAccount});
    // Moov Coin Assert MRM car balance increase and mrm balance decrease
    let actual = await moovCoin.balanceOf(carAccount);
    let expected = carMoovBalance.plus(rideAmount);
    assert.ok(actual.eq(expected), "Rider Moov Coin Balance is not "+expected.toString());
    actual = await moovCoin.balanceOf(mrm.address);
    expected = mrmBalance.sub(rideAmount);
    assert.ok(actual.eq(expected), "MRM Moov Coin Balance is not "+expected.toString());
    // MRM assert rider balance decrease
    actual = await mrm.getBalance(riderAccount);
    expected = riderMRMBalance.sub(rideAmount);
    assert.ok(actual.eq(expected), "Rider Balance is not "+expected.toString());

    // Assert Ride status
    actual = await mrm.getRideStatus(riderAccount);
    expected = 0;
    assert.ok(actual.eq(expected), "Ride status is not in available");
  });


  it("should emit finish event correctly", async () => {
    let moovCoin = await MoovCoin.deployed();
    let mrm = await MRM.new(moovCoin.address);
    var rideAmount = 1;
    var from = "Earth";
    var to = "Mars";
    var riderAccount = web3.eth.coinbase;
    var carAccount = web3.eth.accounts[1]
    await moovCoin.approve(mrm.address, rideAmount, {from:riderAccount});
    await mrm.newRideRequest(from, to, rideAmount, {from:riderAccount});
    await mrm.acceptRideRequest(riderAccount, {from:carAccount});
    let result = await mrm.finishRide({from:riderAccount});
    assert.web3Event(result, {
      event: 'RideFinished',
        args: {
          rider: riderAccount,
          car: carAccount 
        }
    }, 'The event is emitted');
  });
});

// TODO: Test for State Change 
// When AVAILABLE, should not accept or finish
// When REQUESTING, should not create new or finish
// When INPROGRESS, should not create new or accept

//Console variables
//var ferris = Ferris.deployed();
//web3.eth.accounts
//Ferris.deployed().then( function(instance) { return instance.beneficiary.call() })
