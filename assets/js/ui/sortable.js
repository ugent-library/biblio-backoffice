import htmx from "htmx.org";
import Sortable from "sortablejs";

export default function () {
  const addEvents = (content) => {
    var sortables = content.querySelectorAll(".sortable");
    for (var i = 0; i < sortables.length; i++) {
      var sortable = sortables[i];
      new Sortable(sortable, {
        handle: ".sortable-handle",
        animation: 150,
        revertOnSpill: true,
      });
    }
  };

  htmx.onLoad(addEvents);
}
