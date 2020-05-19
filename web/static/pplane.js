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
            // connection socket
            this.socket = null;

            // callback, void (Message)
            this.onMessage = null;

            // callback, void (Conv)
            this.onStateChange = null;
        }

        connect() {
            // TODO sync history with server when successfull
            try {
                this.socket = new WebSocket("ws://" + document.location.host + "/ws");
            } catch (exception) {
                console.error(exception);
            }

            let selfConv = this;
            this.socket.onopen = function (event) {
                selfConv.setState(CON_STATE_CONNECTED_TO_SERVER);
            };
            this.socket.onerror = function (error) {
                selfConv.setState(CON_STATE_NONE);
                console.error(error);
            };
            this.socket.onclose = function (evt) {
                selfConv.setState(CON_STATE_NONE);
                console.log("websocket connection closed, retry in a few seconds");
                setTimeout(() => { selfConv.connect(); }, 3000);
            };
            this.socket.onmessage = function (evt) {
                let obj = JSON.parse(evt.data);
                console.log("incomming data", obj);
            };
        }

        send(msg) {
            this.local_ctr += 1;
            let message = new Message(PAIR_ME, MSG_STATE_LOCAL, msg, new Date(), this.local_ctr);
            this.history.addMessage(message);

            // TODO update state once the server callback is received
            this.socket.send(message.serialize());

            if (this.onMessage) {
                this.onMessage(message);
            }
        }

        setState(state) {
            assert(state == CON_STATE_NONE || state == CON_STATE_CONNECTED_TO_SERVER || state == CON_STATE_CONNECTED_TO_PAIR);
            let changed = this.state != state;
            this.state = state;
            if (changed && this.onStateChange) {
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
            let changed = this.state != state;
            this.state = state;
            if (changed && this.onChange) {
                this.onChange(this);
            }
        }

        serialize() {
            let obj = { data: this.data, timestamp: this.timestamp, counter: this.counter };
            return JSON.stringify(obj);
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

