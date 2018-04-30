require('truffle-test-utils').init();
var MRM = artifacts.require("./MoovRideManager.sol");
var MoovCoin = artifacts.require("./MoovCoin.sol");

contract('MoovRideManager', function(accounts) {
  it("should check Moov Coin Address", async () => {
    let moovCoin = await MoovCoin.deployed();
    let mrm = await MRM.deployed();
    let expected = await mrm.getMoovCoinAddress();
    assert.equal(moovCoin.address, expected, "Moov Coin address is different");
  });


  it("should create new ride request correctly", async () => {
    var rideAmount = 1;
    var from = "Earth";
    var to = "Mars";
    var userAccount = web3.eth.coinbase;
    let moovCoin = await MoovCoin.deployed();
  	let mrm = await MRM.deployed();
    let userMoovBalance = await moovCoin.balanceOf(userAccount);
    let mrmBalance = await moovCoin.balanceOf(mrm.address);
    let userMRMBalance = await mrm.getBalance(userAccount);
    await moovCoin.approve(mrm.address, rideAmount, {from:userAccount});
  	let result = await mrm.newRideRequest(from, to, rideAmount);
  	// Moov Coin Assert Userbalance reduce and MRM balance increase
    let actual = await moovCoin.balanceOf(userAccount);
    let expected = userMoovBalance.sub(rideAmount);
    assert.ok(actual.eq(expected), "User Balance in Moov Coin is not "+expected.toString());
    actual = await moovCoin.balanceOf(mrm.address);
    expected = mrmBalance.plus(rideAmount);
    assert.ok(actual.eq(expected), "Moov coin Balance is not "+expected.toString());
    // MRM assert User balance increase
    actual = await mrm.getBalance(userAccount);
  	expected = userMRMBalance.plus(rideAmount);
  	assert.ok(actual.eq(expected), "User Balance in MRM is not "+expected.toString());
  });

  //   it("should emit bid event correctly", async () => {
  //   let ferrisToken = await FerrisToken.deployed();
  //   let ferris = await Ferris.new(ferrisToken.address);
  //   let actual = await ferris.getBid(accounts[0]);
  //   var bidAmount = 1;
  //   await ferrisToken.approve(ferris.address, bidAmount);
  //   let result = await ferris.bid(bidAmount);
  //   assert.web3Event(result, {
  //     event: 'NewBid',
  //       args: {
  //         eventId: 0,
  //         bidder: accounts[0],
  //         amount: bidAmount 
  //     }
  //   }, 'The event is emitted');
  // });

  // it("should accept correctly", async () => {
  //   let ferrisToken = await FerrisToken.deployed();
  // 	let ferris = await Ferris.deployed();
  //   let beneficiary = await ferris.beneficiary();
  //   var userAccount = web3.eth.coinbase;
  //   var bidAmount = 1;
  //   await ferrisToken.approve(ferris.address, bidAmount);
  // 	assert.isOk(await ferris.bid(bidAmount));
  //   let beneficiaryTokenBalance = await ferrisToken.balanceOf(beneficiary);
  //   let ferrisBalance = await ferrisToken.balanceOf(ferris.address);
  //   let userBidBalance = await ferris.getBid(userAccount);
  // 	let result = await ferris.accept(accounts[0], bidAmount);
  //   // Ferris Token Assert beneficiary balance increase and ferris balance decrease
  //   let actual = await ferrisToken.balanceOf(beneficiary);
  //   let expected = beneficiaryTokenBalance.plus(bidAmount);
  //   assert.ok(actual.eq(expected), "Beneficiary token Balance is not "+expected.toString());
  //   actual = await ferrisToken.balanceOf(ferris.address);
  //   expected = ferrisBalance.sub(bidAmount);
  //   assert.ok(actual.eq(expected), "Ferris token Balance is not "+expected.toString());
  //   // Ferris assert
  //   actual = await ferris.getBid(userAccount);
  //   expected = userBidBalance.sub(bidAmount);
  //   assert.ok(actual.eq(expected), "User Bid Balance is not "+expected.toString());
  // });

  //   it("should emit accept events correctly", async () => {
  //   let ferrisToken = await FerrisToken.deployed();
  //   let ferris = await Ferris.new(ferrisToken.address);
  //   var userAccount = web3.eth.coinbase;
  //   var bidAmount = 1;
  //   await ferrisToken.approve(ferris.address, bidAmount);
  //   assert.isOk(await ferris.bid(bidAmount));
  //   let result = await ferris.accept(userAccount, bidAmount);
  //   assert.web3Event(result, {
  //     event: 'AcceptedBid',
  //       args: {
  //         eventId: 1,
  //         bidder: userAccount,
  //         amount: bidAmount 
  //       }
  //   }, 'The event is emitted');
  // });

  // it("should not accept from non beneficiary", async () => {
  // 	let ferrisToken = await FerrisToken.deployed();
  //   let ferris = await Ferris.deployed();
  //   var userAccount = web3.eth.coinbase;
  //   var nonBeneficiary = accounts[1];
  //   var bidAmount = 1;
  //   await ferrisToken.approve(ferris.address, bidAmount);
  // 	await ferris.bid(bidAmount);
  // 	let err = null;
  // 	try {
  // 		await ferris.accept(userAccount, bidAmount, {from:nonBeneficiary});
  // 	} catch (error) {
  // 		assert(
  //         error.message.search('revert'),
  //         "Expected revert, got '" + error + "' instead",
  //       );
  //       return;
  // 	}
  // 	assert.fail('Expected throw not received');
  // });

  // it("should not accept if less money", async () => {
  //   let ferrisToken = await FerrisToken.deployed();
  // 	let ferris = await Ferris.new(ferrisToken.address);
  //   var userAccount = web3.eth.coinbase;
  //   var bidAmount = 1;
  // 	let err = null;
  // 	try {
  // 		await ferris.accept(userAccount, bidAmount, {from:userAccount});
  // 	} catch (error) {
  // 		assert(
  //         error.message.search('revert'),
  //         "Expected revert, got '" + error + "' instead",
  //       );
  //       return;
  // 	}
  // 	assert.fail('Expected throw not received');
  // });

  // it("should increase bid when user adds", async () => {
  //   let ferrisToken = await FerrisToken.deployed();
  //   let ferris = await Ferris.new(ferrisToken.address);
  //   var userAccount = web3.eth.coinbase;
  //   var bidAmount = 1;
  //   await ferrisToken.approve(ferris.address, bidAmount * 2);
  //   await ferris.bid(bidAmount);
  //   let actual = await ferris.getBid(userAccount);
  //   var expected = 1;
  //   assert.equal(actual, expected, "Value is not 1");
  //   await ferris.bid(bidAmount);
  //   actual = await ferris.getBid(userAccount);
  //   expected = bidAmount * 2;
  //   assert.equal(actual, expected, "Value is not 2");
  // });

  // it("should decrease bid when ben. accepts", async () => {
  //   let ferrisToken = await FerrisToken.deployed();
  //   let ferris = await Ferris.new(ferrisToken.address);
  //   var userAccount = web3.eth.coinbase;
  //   var bidAmount = 1;
  //   await ferrisToken.approve(ferris.address, bidAmount * 2);
  //   await ferris.bid(bidAmount * 2);
  //   let actual = await ferris.getBid(userAccount);
  //   var expected = 2;
  //   assert.equal(actual, expected, "Value is not 2");
  //   await ferris.accept(userAccount, bidAmount);
  //   actual = await ferris.getBid(userAccount);
  //   expected = bidAmount;
  //   assert.equal(actual, expected, "Value is not 1");
  // });

  // it("should withdraw correctly", async () => {
  //   let ferrisToken = await FerrisToken.deployed();
  //   let ferris = await Ferris.new(ferrisToken.address);
  //   var userAccount = web3.eth.coinbase;
  //   var bidAmount = 2;
  //   await ferrisToken.approve(ferris.address, bidAmount);
  //   await ferris.bid(bidAmount);

  //   let userTokenBalance = await ferrisToken.balanceOf(userAccount);
  //   let ferrisBalance = await ferrisToken.balanceOf(ferris.address);
  //   let userBidBalance = await ferris.getBid(userAccount);
  //   await ferris.withdraw();
  //   // Ferris Token Assert Userbalance reduce and ferris balance increase
  //   let actual = await ferrisToken.balanceOf(userAccount);
  //   let expected = userTokenBalance.plus(bidAmount);
  //   assert.ok(actual.eq(expected), "User token Balance is not "+expected.toString());
  //   actual = await ferrisToken.balanceOf(ferris.address);
  //   expected = ferrisBalance.sub(bidAmount);
  //   assert.ok(actual.eq(expected), "Ferris token Balance is not "+expected.toString());
  //   // Ferris assert User bid increase
  //   actual = await ferris.getBid(userAccount);
  //   expected = userBidBalance.sub(bidAmount);
  //   assert.ok(actual.eq(expected), "User Bid Balance is not "+expected.toString());
  // });

  // it("should emit withdraw event correctly", async () => {
  //   let ferrisToken = await FerrisToken.deployed();
  //   let ferris = await Ferris.new(ferrisToken.address);
  //   var userAccount = web3.eth.coinbase;
  //   var bidAmount = 2;
  //   await ferrisToken.approve(ferris.address, bidAmount);
  //   await ferris.bid(bidAmount);
  //   let result = await ferris.withdraw();
  //   assert.web3Event(result, {
  //     event: 'WithdrewBid',
  //       args: {
  //         eventId: 1,
  //         bidder: userAccount,
  //         amount: bidAmount 
  //     }
  //   }, 'The event is emitted');
  // });

  // it("should not withdraw if no money", async () => {
  //   let ferrisToken = await FerrisToken.deployed();
  //   let ferris = await Ferris.new(ferrisToken.address);
  //   let err = null;
  //   try {
  //     await ferris.withdraw();
  //   } catch (error) {
  //     assert(
  //         error.message.search('revert'),
  //         "Expected revert, got '" + error + "' instead",
  //       );
  //       return;
  //   }
  //   assert.fail('Expected throw not received');
  // });


});

//Console variables
//var ferris = Ferris.deployed();
//web3.eth.accounts
//Ferris.deployed().then( function(instance) { return instance.beneficiary.call() })
