var ws = new WebSocket('ws://localhost:3001/ws/'+player.ID)

ws.onopen = () => {
    console.log('ws opened on browser')
}

ws.onmessage = (message) => {
    const obj = JSON.parse(message.data)

    const players = obj.Players
    const bullets = obj.Bullets

    ctx.clearRect(0,0, 500, 500);
    for (var i = 0; i < players.length; i++) {
        if (player.ID == players[i].ID) {
            player = players[i];
            drawPlayer(player);
        } else {
            drawEnemy(players[i])
        }
    }

    for (var i = 0; i < bullets.length; i++) {
        drawBullet(bullets[i])
    }
}