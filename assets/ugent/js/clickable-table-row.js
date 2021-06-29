$(document).ready(function($) {
  $(".clickable-table-row").click(function() {
      window.document.location = $(this).data("href");
  });
});
