// var url = "ws://" + document.location.host + "/api/v1/log/ns/mservice/pod/96143-helloworld-mservice-8644cb8dd7-sfkx4/container/application"
var url = "ws://" + document.location.host + "/api/v1/log/ns/mservice/pod/128330-aosnotice-mservice-56844b8c58-k7blk/container/application"
if (document.location.pathname !== "/") {
    url = "ws://" + document.location.host + "/api/v1" + document.location.pathname
}

document.getElementById('terminal').style.height = window.innerHeight + 'px';
const fitAddon = new FitAddon.FitAddon()
const term = new Terminal({
    "cursorBlink": true
});
term.loadAddon(fitAddon)
window.addEventListener('resize', () => {
    fitAddon.fit()
});

if (window["WebSocket"]) {
    term.open(document.getElementById("terminal"));

    term.onData(data => {
        conn.send(JSON.stringify({operation: "stdin", data: data}))
    });

    term.onResize(data => {
        conn.send(JSON.stringify({operation: "resize", cols: data.cols, rows: data.rows}))
    });

    conn = new WebSocket(url);
    conn.onopen = function (e) {
        conn.send(JSON.stringify({operation: "resize", cols: term.cols, rows: term.rows}))
        fitAddon.fit()
    };
    conn.onmessage = function (event) {
        var msg = JSON.parse(event.data)

        if (msg.operation === "stdout") {
            if (conn.url.includes("/api/v1/log")) {
                term.write(msg.data.replace(/\n/g, "\r\n"))
            } else {
                term.write(msg.data)
            }
        } else {
            console.log("invalid msg operation: " + msg)
        }
    };
    conn.onclose = function (event) {
        if (event.wasClean) {
            console.log(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
        } else {
            console.log('[close] Connection died');
            term.writeln("")
        }
        term.write('Connection Reset By Peer! Try Refresh.');
    };
    conn.onerror = function (error) {
        console.log('[error] Connection error: ', error);
        term.destroy();
    };
} else {
    var item = document.getElementById("terminal");
    item.innerHTML = "<h2>Your browser does not support WebSockets.</h2>";
}