// Required modules
import $ from 'jquery';
require('select2');

htmx.on("htmx:afterSwap", function(evt) {
  if (evt.detail.elt.classList.contains('card-collapsible')) {
    setTimeout(function() {
      function setMultiIconsText(state) {
        if (!state.id) { return state.text; }

        var icons = JSON.parse(state.element.dataset.icon);

        var iconTemplate = '';

        icons.forEach((icon) => {
          iconTemplate+= '<i class="ml-2 if--small text-muted if if-' + icon + '"></i>'
        })
        var $multiIcons = $('<div class="d-flex align-items-center">' + state.text + iconTemplate + '</div>');

        return $multiIcons;
      };

      function setIconText (state) {
        if (!state.id) { return state.text; }
        var $state = $('<div class="d-flex align-items-center"><i class="mr-2 if if-' + state.element.dataset.icon + '"/></i>' + state.text + '</div>');
        return $state;
      };

      $('select').select2({
        theme: 'bootstrap4',
        tags: true,
        minimumResultsForSearch: Infinity
      });

      $(".select-icon").select2({
        theme: 'bootstrap4',
        templateResult: setIconText,
        templateSelection: setIconText,
        minimumResultsForSearch: Infinity
      });

      $(".select-multi-icons").select2({
        theme: 'bootstrap4',
        templateResult: setMultiIconsText,
        templateSelection: setMultiIconsText,
        minimumResultsForSearch: Infinity
      });
    }, 30)
  }
})