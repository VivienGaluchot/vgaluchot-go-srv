const pplane = function () {
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

    function makeUid(length) {
        var result = '';
        var characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
        var charactersLength = characters.length;
        for (var i = 0; i < length; i++) {
            result += characters.charAt(Math.floor(Math.random() * charactersLength));
        }
        return result;
    }


    class Conv {
        constructor() {
            // TODO fetch items from local storage
            // conversation id
            this.cid = makeUid(16);
            // local user id
            this.localUid = makeUid(16);
            // send counter
            this.localCtr = 0;
            // history of messages
            this.history = new History();

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
            try {
                this.socket = new WebSocket("ws://" + document.location.host + "/ws");
            } catch (exception) {
                console.error(exception);
            }

            let selfConv = this;
            this.socket.onopen = function (event) {
                selfConv.setState(CON_STATE_CONNECTED_TO_SERVER);
                // send local messages
                selfConv.history.forEachLocalMsg(function (msg) {
                    try {
                        selfConv.socket.send(msg.serialize());
                    } catch (exception) {
                        console.error(exception);
                    }
                });
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

                try {
                    selfConv.history.getMessage(obj.uid, obj.counter).setState(obj.state);
                } catch (exception) {
                    console.error(exception);
                }
            };
        }

        send(msg) {
            this.localCtr += 1;
            let message = new Message(this.localUid, MSG_STATE_LOCAL, msg, new Date(), this.localCtr);
            this.history.addMessage(message);
            try {
                this.socket.send(message.serialize());
            } catch (exception) {
                console.error(exception);
            }
            if (this.onMessage) {
                this.onMessage(message);
            }
        }

        setState(state) {
            assert(state == CON_STATE_NONE || state == CON_STATE_CONNECTED_TO_SERVER || state == CON_STATE_CONNECTED_TO_PAIR);
            if (this.state < state) {
                this.state = state;
                if (this.onStateChange) {
                    this.onStateChange(this);
                }
            }
        }
    };

    class History {
        constructor() {
            // map : (uid, counter) -> Message
            this.msgMap = new Map();
        }

        getMessage(uid, counter) {
            let key = (uid, counter);
            return this.msgMap.get(key);
        }

        // callback : void (Message)
        forEachLocalMsg(callback) {
            for (let msg of this.msgMap.values()) {
                if (msg.state == MSG_STATE_LOCAL) {
                    callback(msg);
                }
            }
        }

        addMessage(msg) {
            assert(msg instanceof Message);
            let key = (msg.uid, msg.counter);
            if (!this.msgMap.has(key)) {
                this.msgMap.set(key, msg);
            } else {
                console.error("existing message", key)
            }
        }
    };

    class Message {
        constructor(uid, state, data, timestamp, counter) {
            this.uid = uid;
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
            console.log("setState", this, this.state, state);
            if (this.state < state) {
                this.state = state;
                if (this.onChange) {
                    console.log("change");
                    this.onChange(this);
                }
            }
        }

        serialize() {
            let obj = { uid: this.uid, data: this.data, timestamp: this.timestamp, counter: this.counter };
            return JSON.stringify(obj);
        }
    };

    return {
        Conv: Conv,

        MSG_STATE_LOCAL: MSG_STATE_LOCAL,
        MSG_STATE_SENT_TO_SRV: MSG_STATE_SENT_TO_SRV,
        MSG_STATE_SENT_TO_PAIR: MSG_STATE_SENT_TO_PAIR,
        MSG_STATE_READ_BY_PAIR: MSG_STATE_READ_BY_PAIR,

        CON_STATE_NONE: CON_STATE_NONE,
        CON_STATE_CONNECTED_TO_SERVER: CON_STATE_CONNECTED_TO_SERVER,
        CON_STATE_CONNECTED_TO_PAIR: CON_STATE_CONNECTED_TO_PAIR
    }
}();

