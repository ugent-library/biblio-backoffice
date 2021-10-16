import htmx from 'htmx.org';

// Set an event handler on any buttons that have the 'modal-close' class.
// This ensures that any and all modals can be closed via 'close', 'cancel', etc. buttons.
export default function() {
    let modalClose = function(e) {
        let modal = document.querySelectorAll(".modal").item(0)
        let backdrop = document.querySelectorAll(".modal-backdrop").item(0)

        if (modal) {
            modal.classList.remove("show")
        }

        if (backdrop) {
            backdrop.classList.remove("show")
        }

        // Timeout gives us a fluid animation
        setTimeout(function() {
            if (backdrop) {
                backdrop.remove();
            }

            if (modal) {
                modal.remove();
            }
        }, 100)
    }

    // Close the modal after the item was deleted from the backend.
    //
    // If we tried to add the modal-close class directly to the "confirm" button
    // on a confirmation button, the modal would be removed from the DOM before
    // the itemDeleted event could be triggered. Without the modal nodes present
    // in the DOM, the event won't be registered correctly by HTMX. As a result,
    // other triggers listening for the event won't execute. Instead, we
    // use the event itself as a trigger for the modal to close.
    htmx.on("itemDeleted", function(evt) {
        modalClose();
    });

    htmx.on("htmx:afterSettle", function(evt) {
        let container = evt.detail.target

        if (container.classList.contains("modals")) {
            container.querySelectorAll(".modal-close").forEach( el =>
                el.addEventListener("click", modalClose)
            )
        }
    });
}