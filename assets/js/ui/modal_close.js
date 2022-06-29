import htmx from 'htmx.org';

// Set an event handler on any buttons that have the 'modal-close' class.
// This ensures that any and all modals can be closed via 'close', 'cancel', etc. buttons.
export default function() {
    let modalClose = function(evt) {
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

    htmx.onLoad(function(el) {
        el.querySelectorAll(".modal-close").forEach(function (btn) {
            btn.addEventListener("click", modalClose);
        });
    });
}