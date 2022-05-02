import htmx from 'htmx.org';

export default function () {
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

    let goToURL = function(evt) {
        window.location = this.dataset.url
    }

    let addEvents = function() {
        document.querySelectorAll('.radio-card-group .c-radio-card').forEach(card =>
            card.addEventListener('click', toggleSelected)
        )
        document.querySelectorAll('.c-radio-card[data-url]').forEach(card =>
            card.addEventListener('click', goToURL)
        )
    }

    addEvents()

    // TODO don't use afterSettle
    htmx.on("htmx:afterSettle", function(evt) {
        addEvents()
    });
}