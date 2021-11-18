export default function() {
    document.body.addEventListener('htmx:configRequest', (evt) => {
        evt.detail.headers['X-CSRF-Token'] = document.querySelector('meta[name="csrf-token"]').content
    })
}