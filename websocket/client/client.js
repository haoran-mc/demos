window.onload = function () {
  var ws = new WebSocket("ws://localhost:8080");
  var oText = document.getElementById("message");
  var oSend = document.getElementById("send");
  var oUl = document.getElementsByTagName("ul")[0];
  ws.onopen = function () {
    ws.send("hello world");
    oSend.onclick = function () {
      if (!/^\s*$/.test(oText.value)) {
        ws.send(oText.value);
      }
    };
  };
  ws.onmessage = function (msg) {
    console.log("msg: " + msg);

    var str = "<li>" + msg.data + "</li>";
    oUl.innerHTML += str;
  };
  ws.onclose = function (e) {
    console.log("Disconnected from server.");
    ws.close();
  };
};
