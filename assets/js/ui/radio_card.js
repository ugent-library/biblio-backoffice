import htmx from 'htmx.org';

export default function () {
    let addEvents = function(parentEl) {
        // TODO why is this fired twice for each click?
        let toggleSelected = function(evt) {
            let group = evt.currentTarget.closest('.radio-card-group')
            let cards = group.querySelectorAll('.c-radio-card')
            cards.forEach(function(card) {
                card.setAttribute('aria-selected', 'false');
                card.classList.remove('c-radio-card--selected');
            })
            evt.currentTarget.setAttribute('aria-selected', 'true');
            evt.currentTarget.classList.add('c-radio-card--selected');
        }

        parentEl.querySelectorAll('.c-radio-card').forEach(el =>
            el.addEventListener('click', toggleSelected)
        )
    }

    addEvents(document)

    htmx.on("htmx:afterSettle", function(evt) {
        addEvents(evt.detail.target)
    });
}