import $ from 'jquery';

let cards;

$(window).bind("load", function() {
  if (document.querySelector('.c-radio-card')) {
    cards = document.querySelectorAll('.c-radio-card');
    for(let i = 0; i < cards.length; i++) {
      cards[i].addEventListener('click', toggleAttribute)
    }
  }
});

const toggleAttribute = e => {
  for(let i = 0; i < cards.length; i++) {
    cards[i].setAttribute('aria-selected', 'false');
    cards[i].classList.remove('c-radio-card--selected');
  }
  e.currentTarget.setAttribute('aria-selected', 'true');
  e.currentTarget.classList.add('c-radio-card--selected');
}
