import htmx from 'htmx.org';

export default function () {
    let toggleSelected = function(evt) {
        let group = evt.currentTarget.closest('.radio-card-group')
        let cards = group.querySelectorAll('.c-radio-card')
        cards.forEach(function(card) {
            card.setAttribute('aria-selected', 'false')
            card.classList.remove('c-radio-card--selected')
        })
        evt.currentTarget.setAttribute('aria-selected', 'true')
        evt.currentTarget.classList.add('c-radio-card--selected')
    }

    let addEvents = function(rootEl) {
        rootEl.querySelectorAll('.radio-card-group .c-radio-card').forEach(card =>
            card.addEventListener('click', toggleSelected)
        )
    }

    addEvents(document)

    htmx.onLoad(addEvents)
}