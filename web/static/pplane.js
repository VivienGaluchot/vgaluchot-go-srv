const pplane = function () {
    const PAIR_ME = 0;
    const PAIR_OTHER = 1;

    const STATE_LOCAL = 0;
    const STATE_SENT_TO_SRV = 1;
    const STATE_SENT_TO_PAIR = 2;
    const STATE_READ_BY_PAIR = 3;

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

            // callback, void (Message)
            this.onMessage = null;
        }

        connect() {
            // TODO connect or wait pair through websocket
            // sync history with pair when successfull
        }

        send(msg) {
            this.local_ctr += 1;
            let message = new Message(PAIR_ME, STATE_LOCAL, msg, new Date(), this.local_ctr);
            this.history.addMessage(message);
            // TODO sync history with pair
            if (this.onMessage) {
                this.onMessage(message);
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
            assert(state == STATE_LOCAL || state == STATE_SENT_TO_SRV || state == STATE_SENT_TO_PAIR || state == STATE_READ_BY_PAIR);
            this.state = state;
            this.data = data;
            this.timestamp = timestamp;
            this.counter = counter;

            // callback, void (Message)
            this.onChange = null;
        }

        setState(state) {
            assert(state == STATE_LOCAL || state == STATE_SENT_TO_SRV || state == STATE_SENT_TO_PAIR || state == STATE_READ_BY_PAIR);
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

        STATE_LOCAL: STATE_LOCAL,
        STATE_SENT_TO_SRV: STATE_SENT_TO_SRV,
        STATE_SENT_TO_PAIR: STATE_SENT_TO_PAIR,
        STATE_READ_BY_PAIR: STATE_READ_BY_PAIR
    }
}();

