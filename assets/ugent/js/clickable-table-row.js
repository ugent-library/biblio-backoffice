$('.clickable-table-row').on('click',function (e){
    if (!$(e.target).closest('button, a, input, .form-check label').not(this).length) {
      window.document.location = $(this).data("href");
    }
});
