import htmx from 'htmx.org';

export default function () {
    let dismissToast = function(evt) {
        evt.currentTarget.closest('.bc-toast').classList.add('d-none')
    }

    let addEvents = function() {
        document.querySelectorAll('.bc-toast .toast-close').forEach(btn =>
            btn.addEventListener('click', dismissToast)
        )
    }

    addEvents()

    // TODO don't use afterSettle
    htmx.on("htmx:afterSettle", function(evt) {
        addEvents()
    });
}