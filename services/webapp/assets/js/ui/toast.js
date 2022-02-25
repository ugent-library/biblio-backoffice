import htmx from 'htmx.org';

export default function () {
    let addEvents = function(rootEl) {
        rootEl.querySelectorAll('.bc-toast .toast-dismiss').forEach(function(btn) {
            btn.addEventListener('click', () => {btn.closest('.bc-toast').remove()})
        })
        rootEl.querySelectorAll('.bc-toast[data-dismiss-after]').forEach(function(toast) {
            const t = parseInt(toast.dataset.dismissAfter, 10)
            setTimeout(() => {toast.remove()}, t)
        })
    }

    htmx.onLoad(function(rootEl) {
        addEvents(rootEl)
    })
}