require('../ugent/js/index');
htmx = require('htmx.org');

htmx.on("htmx:afterSwap", function(evt) {
    callback = function() {
        card = 'publication-details';
        document
            .getElementsByClassName(card)[0]
            .getElementsByClassName('collapse')[0]
            .classList.add('show')

        document
            .getElementsByClassName(card)[0]
            .getElementsByClassName('collapse-trigger')[0]
            .setAttribute("aria-expanded", "true");
            setTimeout(callback, 40);
    }

    setTimeout(callback, 40)
});