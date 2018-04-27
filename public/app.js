


ws = new WebSocket('ws://' + window.location.host + '/ws');
ws.addEventListener('message', function(e) {
    msg = JSON.parse(e.data);
    document.getElementById('Car'+msg.id).style.top = parseInt(msg.y)+"px"
    document.getElementById('Car'+msg.id).style.left = parseInt(msg.x)+"px"
    document.getElementById('Car'+msg.id).style.transform  = "rotate("+(parseInt(msg.orientation)+180)+"deg)";
});

var example = document.getElementById('Map'); 
example.onclick = function(e) { 
    var x = e.pageX - this.offsetLeft; 
    var y = e.pageY - this.offsetTop; 
    console.log(x+" "+y)
}


