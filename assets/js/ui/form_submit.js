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
    //
    // @todo: Currently only targets the button in the toolbar. Should become more generic
    //   and target groups of outgoing links and buttons in the future.
    htmx.on("htmx:afterSettle", function(evt) {
        // Find .btn-save buttons on dataset/description and publication/description templates
        const buttons = document.querySelector(".btn-save");
        if (buttons !== null) {
            const publishToBiblio = document.querySelector('.bc-toolbar-right .btn-submit-description');
            if (publishToBiblio !== null) {
                publishToBiblio.setAttribute("data-toggle", "modal");
                publishToBiblio.setAttribute("data-target", "#confirmation-next-step");
            }
        }
        else {
            const publishToBiblio = document.querySelector('.bc-toolbar-right .btn-submit-description');
            if (publishToBiblio !== null) {
                publishToBiblio.removeAttribute("data-toggle");
                publishToBiblio.removeAttribute("data-target");
            }
        }

        BSN.initCallback();
    });

}