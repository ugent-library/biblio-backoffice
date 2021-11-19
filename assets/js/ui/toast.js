import htmx from 'htmx.org';

export default function () {
    let dismissToast = function(evt) {
        evt.currentTarget.closest('.bc-toast').classList.add('d-none')
    }

    let addEvents = function() {
        document.querySelectorAll('.bc-toast .toast-dismiss').forEach(btn =>
            btn.addEventListener('click', dismissToast)
        )
        document.querySelectorAll('.bc-toast[data-dismiss-after]').forEach(function(toast) {
            const t = parseInt(toast.dataset.dismissAfter, 10)
            setTimeout(() => {toast.classList.add('d-none')}, t)
        })
    }

    addEvents()

    // TODO don't use afterSettle
    htmx.on("htmx:afterSettle", function(evt) {
        addEvents()
    });
}