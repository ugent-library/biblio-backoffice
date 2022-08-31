$('[data-detail-toggle]').click(function () {
  $(this).closest('.c-master-detail, .c-master-detail-responsive').find('[data-detail]').addClass('active');
  $(this).closest('.c-master-detail, .c-master-detail-responsive').find('[data-master]').removeClass('active');
});

$('[data-master-toggle]').click(function () {
  $(this).closest('.c-master-detail, .c-master-detail-responsive').find('[data-detail]').removeClass('active');
  $(this).closest('.c-master-detail, .c-master-detail-responsive').find('[data-master]').addClass('active');
});
