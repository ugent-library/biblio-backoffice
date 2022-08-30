/* ==========================================================================
    Collapsible header code
   ========================================================================== */
$('[data-scroll-area]').on("scroll", function(e) {
  if($('[data-scroll-area]').scrollTop() > 0){
    $('.c-header-collapsible').addClass('collapsed');
  }
  else{
    $('.c-header-collapsible').removeClass('collapsed');
  }
  });
