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
    if (msg.state == pplane.STATE_LOCAL) {
        state.innerText = "local";
    } else if (msg.state == pplane.STATE_SENT_TO_SRV) {
        state.innerText = "sent to server";
    } else if (msg.state == pplane.STATE_SENT_TO_PAIR) {
        state.innerText = "sent to pair";
    } else if (msg.state == pplane.STATE_READ_BY_PAIR) {
        state.innerText = "read";
    } else {
        state.innerText = "?";
    }

    if (msg.src_pair == pplane.PAIR_ME) {
        div.classList.add("me");
        pair.innerText = "me";
    } else {
        div.classList.add("other");
        pair.innerText = "pair";
    }

    div.appendChild(pair);
    div.innerHTML = div.innerHTML + " ";
    div.appendChild(timestamp);
    div.appendChild(data);
    div.appendChild(state);
    return div;
}

let msg = document.getElementById("msg");
let convEl = document.getElementById("conv");

msg.onkeyup = function (event) {
    if (event.keyCode == 13) {
        if (event.shiftKey) {
            // new line
        } else {
            msg.value = msg.value.substring(0, msg.value.length - 1);
            document.getElementById("sendform").onsubmit();
        }
        return false;
    }
}

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
    }
};

document.getElementById("sendform").onsubmit = function () {
    if (!msg.value) {
        return false;
    }
    pplaneConv.send(msg.value);
    msg.value = "";
    return false;
};