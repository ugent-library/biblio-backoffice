import htmx from 'htmx.org';

export default function () {
    let addEvents = function(rootEl) {
        rootEl.querySelectorAll('.bc-toast .toast-dismiss').forEach(function(btn) {
            btn.addEventListener('click', () => {btn.closest('.bc-toast').remove()})
        })
    }

    htmx.onLoad(addEvents)
}