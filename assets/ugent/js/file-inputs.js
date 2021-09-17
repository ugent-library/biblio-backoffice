import $ from 'jquery';

$(".custom-file input[type=file]").change(function () {

  var fieldVal = $(this).val();

  // Change the node's value by removing the fake path
  fieldVal = fieldVal.replace("C:\\fakepath\\", "");

  if (fieldVal != undefined || fieldVal != "") {
    $(this).next("label").html(fieldVal);
  }

});
