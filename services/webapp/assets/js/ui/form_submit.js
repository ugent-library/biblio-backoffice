import htmx from 'htmx.org';

export default function() {

    // Disable save/cancel buttons when a htmx request is being processed.
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

    // Generic spinner for buttons. Add .btn-display-indicator as a button class. Ensure
    // you've got a spinner added to the button itself e.g.
    //
    // <button type="submit" class="btn btn-primary btn-display-indicator">
    //     <div class="btn-text">Complete description</div>
    //     <i class="if if-arrow-right"></i>
    //     <div class="spinner-border">
    //         <span class="sr-only"></span>
    //     </div>
    // </button>
    function showIndicatorOnButton() {
        document.querySelectorAll('.btn-display-indicator').forEach(function(button) {
            button.addEventListener("click", function(evt) {
                const spinner = button.querySelector('.spinner-border');
                spinner.style.display = "inline-block";
                spinner.style.opacity = "1";
            })
        })
    }

    showIndicatorOnButton();

}