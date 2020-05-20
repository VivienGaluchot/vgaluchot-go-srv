let pplaneConv = new pplane.Conv();

function messageToDiv(msg, div) {
    div.innerHTML = "";
    div.classList.add("message");
    var pair = document.createElement("span");
    pair.classList.add("pair");
    var timestamp = document.createElement("span");
    timestamp.classList.add("timestamp");
    timestamp.innerText = msg.timestamp.toLocaleString();
    var data = document.createElement("span");
    data.classList.add("data");
    data.innerText = msg.data;
    var state = document.createElement("span");
    state.classList.add("state");
    if (msg.state == pplane.MSG_STATE_LOCAL) {
        state.innerText = "local";
    } else if (msg.state == pplane.MSG_STATE_SENT_TO_SRV) {
        state.innerText = "sent to server";
    } else if (msg.state == pplane.MSG_STATE_SENT_TO_PAIR) {
        state.innerText = "sent to pair";
    } else if (msg.state == pplane.MSG_STATE_READ_BY_PAIR) {
        state.innerText = "read";
    } else {
        convStateEl.innerText = `? ${msg.state}`;
    }

    if (msg.uid == pplaneConv.localUid) {
        div.classList.add("me");
        pair.innerText = `me @${msg.uid}`;
    } else {
        div.classList.add("other");
        pair.innerText = `other @${msg.uid}`;
    }

    div.appendChild(pair);
    div.innerHTML = div.innerHTML + " ";
    div.appendChild(timestamp);
    div.appendChild(data);
    div.appendChild(state);
    return div;
}

let msgEl = document.getElementById("i-msg");
let chanEl = document.getElementById("i-chan");
let convEl = document.getElementById("o-conv");
let convStateEl = document.getElementById("o-conv-state");

msgEl.onkeypress = function (event) {
    if (event.keyCode == 13) {
        if (event.shiftKey) {
            msgEl.value = msgEl.value + "\n";
        } else {
            document.getElementById("sendform").onsubmit();
        }
        return false;
    }
};

pplaneConv.onMessage = function (msg) {
    console.log("onMessage", msg);
    var div = document.createElement("div");
    messageToDiv(msg, div);
    var doScroll = convEl.scrollTop > convEl.scrollHeight - convEl.clientHeight - 1;
    convEl.appendChild(div);
    if (doScroll) {
        convEl.scrollTop = convEl.scrollHeight - convEl.clientHeight;
    }
    msg.onChange = function (msg) {
        messageToDiv(msg, div);
    };
};
pplaneConv.onStateChange = function (conv) {
    if (conv.state == pplane.CON_STATE_NONE) {
        convStateEl.classList.remove("ok");
        convStateEl.innerText = "no connection";
    } else if (conv.state == pplane.CON_STATE_CONNECTED_TO_SERVER) {
        convStateEl.classList.add("ok");
        convStateEl.innerText = "connected to server";
    } else if (conv.state == pplane.CON_STATE_CONNECTED_TO_PAIR) {
        convStateEl.classList.add("ok");
        convStateEl.innerText = "connected to pair";
    } else {
        convStateEl.classList.remove("ok");
        convStateEl.innerText = `? ${conv.state}`;
    }
};
pplaneConv.connect();

document.getElementById("sendform").onsubmit = function () {
    if (!msgEl.value) {
        return false;
    }
    pplaneConv.send(msgEl.value);
    msgEl.value = "";
    return false;
};