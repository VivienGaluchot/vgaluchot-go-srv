const pplane = function () {
    const PAIR_ME = 0;
    const PAIR_OTHER = 1;

    const MSG_STATE_LOCAL = 0;
    const MSG_STATE_SENT_TO_SRV = 1;
    const MSG_STATE_SENT_TO_PAIR = 2;
    const MSG_STATE_READ_BY_PAIR = 3;

    const CON_STATE_NONE = 0;
    const CON_STATE_CONNECTED_TO_SERVER = 1;
    const CON_STATE_CONNECTED_TO_PAIR = 2;

    function assert(pred) {
        if (pred != true)
            throw new Error("assert error");
    }

    class Conv {
        constructor() {
            // TODO fetch cid from local storage or create a new random one
            this.cid = null;
            // TODO fetch history from local storage or create a new one
            this.history = new History();
            // TODO fetch local counter from local storage
            this.local_ctr = 0;
            // connection state
            this.state = CON_STATE_NONE;

            // callback, void (Message)
            this.onMessage = null;

            // callback, void (Conv)
            this.onStateChange = null;
        }

        connect() {
            // TODO connect or wait pair through websocket
            // sync history with pair when successfull
            let socket = null;
            try {
                socket = new WebSocket("ws://" + document.location.host + "/ws");
            } catch (exception) {
                console.error(exception);
            }

            let selfConv = this;
            socket.onopen = function (event) {
                selfConv.setState(CON_STATE_CONNECTED_TO_SERVER);
            };
            socket.onerror = function (error) {
                selfConv.setState(CON_STATE_NONE);
                console.error(error);
                console.log("websocket connection failed, retry in a few seconds");
                setTimeout(() => { selfConv.connect(); }, 3000);
            };
            socket.onclose = function (evt) {
                selfConv.setState(CON_STATE_NONE);
                console.log("websocket connection closed, retry in a few seconds");
                setTimeout(() => { selfConv.connect(); }, 3000);
            };

            socket.onmessage = function (evt) {
                console.log("incomming data", evt.data);
                let messages = evt.data.split('\n');
                for (let i = 0; i < messages.length; i++) {
                    let item = document.createElement("div");
                    item.innerText = messages[i];
                    appendLog(item);
                }
            };
        }

        send(msg) {
            this.local_ctr += 1;
            let message = new Message(PAIR_ME, MSG_STATE_LOCAL, msg, new Date(), this.local_ctr);
            this.history.addMessage(message);
            // TODO sync history with pair
            if (this.onMessage) {
                this.onMessage(message);
            }
        }

        setState(state) {
            assert(state == CON_STATE_NONE || state == CON_STATE_CONNECTED_TO_SERVER || state == CON_STATE_CONNECTED_TO_PAIR);
            this.state = state;
            if (this.onStateChange) {
                this.onStateChange(this);
            }
        }
    };

    class History {
        constructor() {
            this.messages = [];
        }

        addMessage(msg) {
            assert(msg instanceof Message);
            this.messages.push(msg);
        }
    };

    class Message {
        constructor(src_pair, state, data, timestamp, counter) {
            assert(src_pair == PAIR_ME || src_pair == PAIR_OTHER);
            this.src_pair = src_pair;
            assert(state == MSG_STATE_LOCAL || state == MSG_STATE_SENT_TO_SRV || state == MSG_STATE_SENT_TO_PAIR || state == MSG_STATE_READ_BY_PAIR);
            this.state = state;
            this.data = data;
            this.timestamp = timestamp;
            this.counter = counter;

            // callback, void (Message)
            this.onChange = null;
        }

        setState(state) {
            assert(state == MSG_STATE_LOCAL || state == MSG_STATE_SENT_TO_SRV || state == MSG_STATE_SENT_TO_PAIR || state == MSG_STATE_READ_BY_PAIR);
            this.state = state;

            if (this.onChange) {
                this.onChange(this);
            }
        }
    };

    return {
        Conv: Conv,

        PAIR_ME: PAIR_ME,
        PAIR_OTHER: PAIR_OTHER,

        MSG_STATE_LOCAL: MSG_STATE_LOCAL,
        MSG_STATE_SENT_TO_SRV: MSG_STATE_SENT_TO_SRV,
        MSG_STATE_SENT_TO_PAIR: MSG_STATE_SENT_TO_PAIR,
        MSG_STATE_READ_BY_PAIR: MSG_STATE_READ_BY_PAIR,

        CON_STATE_NONE: CON_STATE_NONE,
        CON_STATE_CONNECTED_TO_SERVER: CON_STATE_CONNECTED_TO_SERVER,
        CON_STATE_CONNECTED_TO_PAIR: CON_STATE_CONNECTED_TO_PAIR
    }
}();

