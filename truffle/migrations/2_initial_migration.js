var MoovCoin = artifacts.require("MoovCoin");

module.exports = function(deployer) {
  deployer.deploy(MoovCoin);
};