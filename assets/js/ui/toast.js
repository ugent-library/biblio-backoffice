import htmx from "htmx.org/dist/htmx.esm.js";
import { Toast } from "bootstrap.native";

export default function () {
  htmx.onLoad((rootEl) => {
    rootEl.querySelectorAll(".toast").forEach((toastEl) => {
      // Bootstrap already initializes the toast but doesn't show it automatically
      Toast.getInstance(toastEl).show();
    });
  });
}
