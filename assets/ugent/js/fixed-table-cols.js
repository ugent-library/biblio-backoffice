$('.table-col-sm-fixed .dropdown').on('show.bs.dropdown', (event) => {
  let tableCol = event.target.closest('.table-col-sm-fixed');
  $(tableCol).addClass('dropdown-active')
});


$('.table-col-sm-fixed .dropdown').on('hidden.bs.dropdown', (event) => {
  let tableCol = event.target.closest('.table-col-sm-fixed');
  $(tableCol).removeClass('dropdown-active')
});
