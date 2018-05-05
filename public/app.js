
$("[type='number']").keypress(function (evt) {
    evt.preventDefault();
});

ws = new WebSocket('ws://' + window.location.host + '/ws');
ws.addEventListener('message', function(e) {
    msg = JSON.parse(e.data);
    document.getElementById('Car'+msg.id).style.top = parseInt(msg.y)+"px"
    document.getElementById('Car'+msg.id).style.left = parseInt(msg.x)+"px"
    document.getElementById('Car'+msg.id).style.transform  = "rotate("+(parseInt(msg.orientation)+180)+"deg)";
});

window.addEventListener('load', function() {
    // Check if Web3 has been injected by the browser:
  if (typeof web3 !== 'undefined' ) {
    // You have a web3 browser! Continue below!
    startApp(web3);
  } else {
    alert("Get METAMASK!");
     // Warn the user that they need to get a web3 browser
     // Or install MetaMask, maybe with a nice graphic.
  }

  document.getElementById("getMC").onclick = getMCs;
  document.getElementById("approve-mc-button").onclick = approveMC;
  document.getElementById("get-ride-button").onclick = getRide;
  document.getElementById("set-start-point-button").onclick =
    function (e) {
      document.getElementById('Map').onclick = setStartPoint;
      document.getElementById("get-ride-debug").innerHTML = " Click anywhere on map to select start point";
    };
  document.getElementById("set-end-point-button").onclick =
    function (e) {
      document.getElementById('Map').onclick = setEndPoint;
      document.getElementById("get-ride-debug").innerHTML = " Click anywhere on map to select end point";
    };
})

const mrmAddress = '0x6e43523c1081dffae69ff1cdac27f3bf178ba852';
const mrmContractABI =  [{"constant":true,"inputs":[{"name":"addr","type":"address"}],"name":"getRideStatus","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"finishRide","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"chosenRider","type":"address"}],"name":"acceptRideRequest","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"addr","type":"address"}],"name":"getCarAddress","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"from","type":"string"},{"name":"to","type":"string"},{"name":"amount","type":"uint256"}],"name":"newRideRequest","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getAvailableRides","outputs":[{"name":"","type":"address[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"moovCoin","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"cancelRideRequest","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"addr","type":"address"}],"name":"getBalance","outputs":[{"name":"balance","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"addr","type":"address"}],"name":"getDestinations","outputs":[{"name":"from","type":"string"},{"name":"to","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"moovCoinAddress","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"name":"rider","type":"address"},{"indexed":false,"name":"from","type":"string"},{"indexed":false,"name":"to","type":"string"},{"indexed":false,"name":"amount","type":"uint256"}],"name":"NewRideRequest","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"rider","type":"address"},{"indexed":false,"name":"car","type":"address"}],"name":"RideAccepted","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"rider","type":"address"},{"indexed":false,"name":"car","type":"address"}],"name":"RideFinished","type":"event"}];
const moovCoinABI = [{"constant":false,"inputs":[],"name":"corruptExchange","outputs":[],"payable":true,"stateMutability":"payable","type":"function"},{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"INITIAL_SUPPLY","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_subtractedValue","type":"uint256"}],"name":"decreaseApproval","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_addedValue","type":"uint256"}],"name":"increaseApproval","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}];


async function startApp(web3) {
    eth = new Eth(web3.currentProvider);
    var version = await eth.net_version();
    if(version == "3") {
      mrm = eth.contract(mrmContractABI).at(mrmAddress);
      const moovCoinAddress = await mrm.moovCoin();
      moovCoin = eth.contract(moovCoinABI).at(moovCoinAddress[0]);
      coinbase = await eth.coinbase();
      updateView();
    } else {
      document.getElementById("ui").innerHTML = "Change network to Ropsten";
    }
}

async function updateView() {

    document.getElementById("user-address").innerHTML = coinbase;

    const userBalance = await moovCoin.balanceOf(coinbase);
    document.getElementById("user-balance").innerHTML = userBalance[0].toNumber();

    const userApprovedMCs = await(moovCoin.allowance(coinbase, mrmAddress))
    document.getElementById("user-approved-mcs").innerHTML = userApprovedMCs[0].toNumber();
    document.getElementById("approve-mc-field").setAttribute("max", userBalance[0].toNumber() - userApprovedMCs[0].toNumber());

    document.getElementById("get-ride-field").setAttribute("max", userApprovedMCs[0].toNumber());
}

async function waitForTxToBeMined (txHash) {
  $(':button').prop('disabled', true);
  let txReceipt
  while (!txReceipt) {
    try {
      txReceipt = await eth.getTransactionReceipt(txHash)
    } catch (err) {
      console.log("Tx Mine failed");
    }
  }
  console.log("Tx Mined");
  $(':button').prop('disabled', false);
  updateView();
}

function getMCs(){
    console.log("trying to get MC");
    const amount = parseInt(document.getElementById("get-mc-field").value);
    if(amount <= 0) {
      console.log("Come on dude, ask for more than 0");
      return;
    }
    document.getElementById("get-mc-field").value = 0;

    moovCoin.corruptExchange({from:coinbase, value:10000000000000000*amount}).then(function (txHash) {
      console.log('Transaction sent');
      console.dir(txHash);
      waitForTxToBeMined(txHash);
    });;
}

async function approveMC(){
    console.log("trying to approve");
    const amount = parseInt(document.getElementById("approve-mc-field").value);
    if(amount <= 0) {
      console.log("Come on dude, ask for more than 0");
      return;
    }
    document.getElementById("approve-mc-field").value = 0;

    moovCoin.increaseApproval(mrmAddress, amount, { from: coinbase }).then(function (txHash) {
      console.log('Transaction sent');
      console.dir(txHash);
      waitForTxToBeMined(txHash);
    });
}

async function getRide(){
    console.log("trying to increase bid");
    var debugElement = document.getElementById("get-ride-debug");
    debugElement.innerHTML = "";
    const start = getLocations('start-point')
    const end = getLocations('end-point')
    const amount = parseInt(document.getElementById("get-ride-amount-field").value);
    if (start == "location not set"){
      debugElement.innerHTML = "Set Start Point before getting ride";
      return;
    }else if (end == "location not set"){
      debugElement.innerHTML = "Set End Point before getting ride";
      return;
    }else if (amount <= 0 ) {
      document.getElementById("get-ride-debug").innerHTML = "Enter a ride amount greater than 0";
      return;
    }
    console.log(start+" "+end);
    mrm.newRideRequest(start, end, amount, { from: coinbase }).then(function (txHash) {
      console.log('Transaction sent');
      console.dir(txHash);
      waitForTxToBeMined(txHash);
      document.getElementById("get-ride-amount-field").value = 0;
    });

}

function getLocations(locString) {
  if (document.getElementById(locString).style.visibility == "visible") {
    var y = parseInt(document.getElementById(locString).style.top) + 10;
    y -= document.getElementById('Map').offsetTop;
    var x = parseInt(document.getElementById(locString).style.left) + 5;
    x -= document.getElementById('Map').offsetLeft;
    return ""+x+","+y;
  } else {
    return "location not set";
  }
}


function setStartPoint(e){
  document.getElementById('start-point').style.visibility = "visible";
  document.getElementById('start-point').style.top = (e.pageY-10)+"px";
  document.getElementById('start-point').style.left = (e.pageX-5)+"px";
  document.getElementById("get-ride-debug").innerHTML = "";
  document.getElementById('Map').onclick = null;
}



function setEndPoint(e){
  document.getElementById('end-point').style.visibility = "visible";
  document.getElementById('end-point').style.top = (e.pageY-10)+"px";
  document.getElementById('end-point').style.left = (e.pageX-5)+"px";
  document.getElementById("get-ride-debug").innerHTML = "";
  document.getElementById('Map').onclick = null;
}
