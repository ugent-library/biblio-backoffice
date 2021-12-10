import htmx from 'htmx.org';

export default function() {

    // See: https://github.com/bigskysoftware/htmx/issues/394
    htmx.on("htmx:afterSettle", function(evt) {
        const buttons = document.querySelectorAll(".btn-save");
        if (buttons !== undefined) {
            buttons.forEach(function(el) {
                el.addEventListener("click", function (evt) {
                    htmx.on('htmx:beforeRequest', function(e) {
                        // Save button
                        el.setAttribute("disabled", "disabled");
                        // Cancel button, if any
                        if (el.previousElementSibling !== undefined) {
                            el.previousElementSibling.setAttribute("disabled", "disabled");
                        }
                    })
                })
            });
        }
    });
}