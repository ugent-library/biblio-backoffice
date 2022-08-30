// jQuery
const $ = require('jquery');
window.$ = $;

// Bootstrap bundle
require('bootstrap/dist/js/bootstrap.bundle');

$(function () {

  // Custom javascript
  require('./clickable-table-row');
  require('./collapse');
  require('./collapsible-header');
  require('./collapsible-sub-sidebar');
  require('./editable-panels');
  require('./file-inputs');
  require('./fixed-table-cols');
  require('./flatpickr');
  require('./master-detail');
  require('./popover');
  require('./radio-card');
  require('./selects');
  require('./tablesort');
  require('./tooltips');
});


// Require JS to render prototype,
// styleguide and navigation.
// Remove this line when going to production.
//
require('../../core/js/index');

