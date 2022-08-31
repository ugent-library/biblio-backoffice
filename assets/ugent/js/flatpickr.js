import flatpickr from "flatpickr";
flatpickr(".rangeDate", {
  mode: "range",
  dateFormat: "d-m-Y",
  maxDate: "today"
});
