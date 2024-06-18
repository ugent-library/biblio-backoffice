import htmx from "htmx.org";
import Tagify from "@yaireo/tagify"; // see https://github.com/yairEO/tagify

export default function () {
  /*
   * Expected html layout:
   *
   * <div class="tags">
   *  <span class="tags-real-values"></span>
   *  <textarea class="tags-widget-values" data-input-name="keyword">
   *    ["tag1", "tag2"]
   *  </textarea>
   * </div>
   *
   * "tags-real-values" is filled with hidden input (by the widget) that are sent to the server
   * "tags-widget-values" is a textarea that should be prefilled with json array of tags, that should have NO name,
   * and should therefore not be sent to the server
   * */

  const fillRealValues = function (element, name, values) {
    element.innerHTML = "";
    for (let i = 0; i < values.length; i++) {
      const input = document.createElement("input");
      input.type = "hidden";
      input.name = `${name}[${i}]`;
      input.value = values[i];
      element.appendChild(input);
    }
  };

  htmx.onLoad((rootEl) => {
    rootEl.querySelectorAll(".tags").forEach((tag) => {
      const realValues = tag.querySelector(".tags-real-values");
      const widgetValues = tag.querySelector(".tags-widget-values");
      const { inputName } = widgetValues.dataset;

      // parse json, and fill tags-real-values
      {
        let val = widgetValues.value;
        try {
          val = JSON.parse(val);
          if (Array.isArray(val)) {
            fillRealValues(realValues, inputName, val);
          }
        } catch (err) {}
      }

      // load and configure tagify
      const tagify = new Tagify(widgetValues, {
        delimiters: ";|\n|\r",
        duplicates: false,
        pasteAsTags: true,

        // we have no dropdown, but setting "caseSensitive" is used by Tagify for duplicate check
        dropdown: {
          enabled: false,
          caseSensitive: true,
        },
      });

      tagify.on("change", function (evt) {
        fillRealValues(
          realValues,
          inputName,
          evt.detail.tagify.value.map(function (v) {
            return v.value;
          }),
        );
      });
    });
  });
}
