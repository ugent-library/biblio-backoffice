// Tooltips
const $ = require('jquery');

// button tooltip from html element
$('.fromElement').each(function(){
  let title = $(this).find('.tooltip-text').html();
  $(this).attr('title', title);
 });

//alert($(window).width());
if ($(window).width() > 991) {
  $('[data-toggle="tooltip"]').tooltip();
}
