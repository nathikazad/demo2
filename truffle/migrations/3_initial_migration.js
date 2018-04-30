var MoovCoin = artifacts.require("MoovCoin");
var MoovRideManager = artifacts.require("MoovRideManager");

module.exports = function(deployer) {
	deployer.deploy(MoovRideManager, MoovCoin.address);
};
