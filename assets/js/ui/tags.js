import htmx from 'htmx.org';
import Tagify from '@yaireo/tagify';// see https://github.com/yairEO/tagify

export default function () {

    /*
     *
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

    let fillRealValues = function(element, name, values){
      element.innerHTML = "";
      for(let i = 0;i < values.length;i++) {
        let input = document.createElement("input");
        input.type = "hidden"
        input.name = name + "[" + i + "]";
        input.value = values[i];
        element.appendChild(input);
      }
    }

    let configureTags = function(element){

      let realValues = element.querySelector(".tags-real-values");
      let widgetValues = element.querySelector(".tags-widget-values");
      let inputName = widgetValues.dataset.inputName;

      // parse json, and fill tags-real-values
      {
        let val = widgetValues.value
        try {
          val = JSON.parse(val);
          if (val instanceof Array) {
            fillRealValues(realValues, inputName, val);
          }
        } catch(err) {}
      }

      //load and configure tagify
      let t = new Tagify(widgetValues, {
        delimiters: ";|\n|\r",
        duplicates: false,
        pasteAsTags: true,//automatically converted pasted text into tags (using delimiters)
        // we have no dropdown, but setting "caseSensitive" is used by tagify for duplicate check
        dropdown: {
          enabled: false,
          caseSensitive: true
        }
      });
      t.on("change", function(evt) {
        fillRealValues(
          realValues,
          inputName,
          evt.detail.tagify.value.map(function(v){ return v.value; })
        );
      });

    };

    let loadTagify = function(rootEl) {
      let allTags = rootEl.querySelectorAll(".tags");
      for(let i = 0;i < allTags.length;i++){
        configureTags(allTags[i]);
      }
    }

    htmx.onLoad(loadTagify)
}
