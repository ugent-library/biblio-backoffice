import htmx from 'htmx.org';

export default function () {
    let addEvents = function(parentEl) {
        let dismissToast = function(evt) {
            evt.currentTarget.closest('.bc-toast').forEach(function(toast) {
                toast.classList.add('d-none');
            })
        }

        parentEl.querySelectorAll('.bc-toast .toast-close').forEach(btn =>
            btn.addEventListener('click', dismissToast)
        )
    }

    addEvents(document)

    // TODO don't use afterSettle
    htmx.on("htmx:afterSettle", function(evt) {
        addEvents(evt.detail.target)
    });
}