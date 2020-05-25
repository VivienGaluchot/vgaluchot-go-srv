const local = (() => {
    function timeInS() {
        let x = Date.now() / 1000;
        return x;
    };

    return {
        installTypingAnimation(element) {
            let text = element.innerHTML;

            let start = timeInS();
            function typing(text, index) {
                if (Math.round(timeInS() - start) % 2 != 0) {
                    hider = (_) => '&nbsp;';
                } else {
                    hider = (_, i) => i == 0 ? '_' : '&nbsp;';
                }
                return text.slice(0, index) + Array.from(text.slice(index)).map(hider).join('');
            };

            let i = 0;
            function update() {
                element.innerHTML = typing(text, i);
                i++;
                if (i > text.length)
                    return;
                window.setTimeout(update, 10 + 200 * Math.random());
            };

            element.innerHTML = typing(text, 0)
            window.setTimeout(update, 1000);
        },

        installAnchor(element) {
            console.log(element);
            let link = document.createElement("a");
            link.setAttribute("href", "#" + element.id);
            link.innerHTML = "<i class=\"fas fa-link\"></i>";
            link.classList.add("anchor");
            element.prepend(link);
        }
    };
})();

window.addEventListener("DOMContentLoaded", (event) => {
    for (let el of document.getElementsByClassName("quote")) {
        local.installTypingAnimation(el);
    }

    for (let el of document.getElementsByClassName("anchorable")) {
        local.installAnchor(el);
    }
});