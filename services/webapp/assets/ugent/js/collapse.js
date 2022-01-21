import $ from 'jquery';

$('[data-panel-toggle]').on('click', function () {
  var id = $(this).attr("data-panel-toggle");
  var collapsibleContent = $('#'+id).parent();
  collapsibleContent.collapse('show');
})
