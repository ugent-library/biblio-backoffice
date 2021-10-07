import htmx from 'htmx.org';

// Set an event handler on any buttons that have the 'modal-close' class.
// This ensures that any and all modals can be closed via 'close', 'cancel', etc. buttons.
export default function() {
    htmx.on("htmx:afterSettle", function(evt) {
        let container = evt.detail.target

        if (container.classList.contains("modals")) {
            let modalClose = function(e) {
                let modal = container.querySelectorAll(".modal").item(0)
                let backdrop = container.querySelectorAll(".modal-backdrop").item(0)

                modal.classList.remove("show")
                backdrop.classList.remove("show")

                setTimeout(function() {
                    container.removeChild(backdrop)
                    container.removeChild(modal)
                }, 100)
            }

            container.querySelectorAll(".modal-close").forEach( el =>
                el.addEventListener("click", modalClose)
            )
        }

    });
}