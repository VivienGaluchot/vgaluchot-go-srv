const local = (() => {
    function timeInS() {
        let x = Date.now() / 1000;
        return x;
    };

    function setCookie(name, value, daysToLive) {
        // Encode value in order to escape semicolons, commas, and whitespace
        var cookie = name + "=" + encodeURIComponent(value);
        if (typeof daysToLive === "number") {
            /* Sets the max-age attribute so that the cookie expires
            after the specified number of days */
            cookie += `; max-age=${daysToLive * 24 * 60 * 60}; SameSite=Strict`;
            document.cookie = cookie;
        }
    }

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
            let link = document.createElement("a");
            link.setAttribute("href", "#" + element.id);
            link.innerHTML = "<i class=\"fas fa-link\"></i>";
            link.classList.add("anchor");
            element.prepend(link);
        },

        switchLanguage(lang) {
            setCookie('lang', lang, 30);
            window.location.reload();
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