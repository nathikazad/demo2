


ws = new WebSocket('ws://' + window.location.host + '/ws');
ws.addEventListener('message', function(e) {
    msg = JSON.parse(e.data);
    if (msg.id =="1")  {
    document.getElementById('Cars').style.top = parseInt(msg.y)-55+"px"
    document.getElementById('Cars').style.left = parseInt(msg.x)-55+"px"
    document.getElementById('Car1').style.transform  = "rotate("+(parseInt(msg.orientation)+180)+"deg)";
  }
});

var example = document.getElementById('Map'); 
example.onclick = function(e) { 
    var x = e.pageX - this.offsetLeft; 
    var y = e.pageY - this.offsetTop; 
    console.log(x+" "+y)
}


