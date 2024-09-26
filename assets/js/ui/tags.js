import htmx from "htmx.org/dist/htmx.esm.js";
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
    rootEl.querySelectorAll(".tags").forEach((tags) => {
      const realValues = tags.querySelector(".tags-real-values");
      const widgetValues = tags.querySelector(".tags-widget-values");
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
        pasteAsTags: true, //automatically converted pasted text into tags (using delimiters)

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
          evt.detail.tagify.value.map((v) => v.value),
        );
      });

      const { originalInput, input } = tagify.DOM;
      if (originalInput && input) {
        handleLabelFocus(originalInput, input);

        copyAriaAttributes(originalInput, input);
      }
    });
  });
}

function handleLabelFocus(originalInput, input) {
  if (originalInput.id) {
    const label = document.querySelector(`label[for="${originalInput.id}"]`);

    if (label) {
      label.addEventListener("click", () => {
        // for some reason focus is lost again immediately if you don't set it with timeout
        window.setTimeout(() => input.focus(), 0);
      });
    }
  }
}

function copyAriaAttributes(originalInput, input) {
  if (originalInput && input) {
    for (const attr of originalInput.attributes) {
      if (attr.name.startsWith("aria-") && !input.hasAttribute(attr.name)) {
        input.setAttribute(attr.name, attr.value);
      }
    }
  }
}
