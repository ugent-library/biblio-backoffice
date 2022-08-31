$(function () {
  $('[data-toggle="popover"]').popover({html:true})
})

$(function(){
  $("[data-toggle=popover-custom]").popover({
      html : true,
      content: function() {
        var content = $(this).attr("data-popover-content");
        return $(content).children(".popover-body").html();
      },
      title: function() {
        var content = $(this).attr("data-popover-content");
        if ($(content).children(".popover-heading").length > 0) {
          return $(content).children(".popover-heading").html();
        } else {
          return '';
        }
      }
  });
});
