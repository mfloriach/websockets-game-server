var stage = document.getElementById('svs')
var ctx = stage.getContext('2d'); 

var player = {
   ID:  "ww",
   color: "black",
   X: 10,
   Y: 10
}

function drawPlayer(p) {
   ctx.fillStyle = "black"; 
   ctx.fillRect(p.X, p.Y, 20, 20); 
}

function drawEnemy(p) {
   ctx.fillStyle = "red"; 
   ctx.fillRect(p.X, p.Y, 20, 20); 
}

function drawBullet(p) {
   ctx.fillStyle = 'green';
   ctx.fillRect(p.X, p.Y, 10, 10); 
}

drawPlayer(player)

window.onkeydown = function(event) {
    var keyPr = event.keyCode; //Key code of key pressed

    var type = ""
    if(keyPr === 39 && player.X<=460){ 
      type = "right"
    }
    else if(keyPr === 37 && player.X>10){
      type = "left"
    }
    else if(keyPr === 38 && player.Y>10) {
      type = "down"
    }
    else if(keyPr === 40 && player.Y<=460){
      type = "up"
    }
    else if(keyPr === 65){
      type = "shoot"
   }

   ws.send(JSON.stringify({"type": type}));
};

