import htmx from 'htmx.org';
import BSN from "bootstrap.native/dist/bootstrap-native-v4";

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

    // Show a warning / confirmation when user navigates away from form with
    // potential unsaved data.
    htmx.on("htmx:afterSettle", function(evt) {
        // Find .btn-save buttons on dataset/description and publication/description templates
        const buttons = document.querySelector(".btn-save");

        function confirmationHandler(el) {
            let callback = function (evt) {
                const submit = document.querySelector('#confirmation-next-step a.btn-primary');
                submit.href = el.href;
            }

            if (buttons !== null) {
                el.setAttribute("data-toggle", "modal");
                el.setAttribute("data-target", "#confirmation-next-step");

                // Set URL from the link in the modal
                el.addEventListener("click", callback)
            }
            else {
                el.removeAttribute("data-toggle");
                el.removeAttribute("data-target");
            }
        }

        document.querySelectorAll('.btn-confirmation-next-step').forEach(function(el) {
            confirmationHandler(el);
        });

        BSN.initCallback();
    });

}