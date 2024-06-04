import htmx from "htmx.org";
import modalError from "./modal_error.js";

export default function () {
  htmx.onLoad(function (rootEl) {
    rootEl
      .querySelectorAll("input[type=file].upload-progress")
      .forEach((input) => {
        input.addEventListener("change", (evt) => {
          const files = Array.from(input.files);
          if (!files.length) return;

          const file = files[0];
          let form = input.closest("form");
          let target = document.querySelector(form.dataset.target);
          let headers = JSON.parse(form.dataset.headers);
          const maxSize = parseInt(input.dataset.maxSize);
          const maxSizeError = input.dataset.maxSizeError;

          if (!isNaN(maxSize) && maxSize > 0 && input.files[0].size > maxSize) {
            evt.preventDefault();
            evt.stopPropagation();
            modalError(maxSizeError);
            hideFormUpload(form);
            return;
          }

          // weird, but makes sure that middleware does not try to read _method from form
          headers["X-HTTP-Method-Override"] = "POST";
          //"Failed to execute 'setRequestHeader' on 'XMLHttpRequest': String contains non ISO-8859-1 code point"
          headers["X-Upload-Filename"] = encodeURIComponent(file.name);
          headers["Content-Type"] = file.type;

          let req = new XMLHttpRequest();

          // request aborted by user
          /*
        req.addEventListener(
          "abort",
          (e) => {
            e.preventDefault();
            e.stopPropagation();
            modalError(input.dataset.uploadMsgFileAborted);
            hideFormUpload(form);
          },
          false
        );*/

          req.upload.addEventListener(
            "progress",
            (e) => {
              if (!e.lengthComputable) return;
              setProgress(form, Math.floor((e.loaded / e.total) * 100));
            },
            false,
          );

          req.addEventListener("readystatechange", (e) => {
            if (req.readyState !== 4) return;

            hideFormUpload(form);

            // file created
            if (req.status == 200 || req.status == 201) {
              target.innerHTML = req.responseText;
              htmx.process(target);
            }
            // file too large (server side check)
            else if (req.status == 413) {
              modalError(input.dataset.uploadMsgFileTooLarge);
            }
            // publication has been removed in the meantime
            else if (req.status == 404) {
              modalError(input.dataset.uploadMsgRecordNotFound);
            }
            // undetermined errors
            else {
              modalError(input.dataset.uploadMsgUnexpected);
            }
          });

          // switch to form upload
          showFormUpload(form);

          // and start request
          req.open(form.method, form.action);
          for (let key in headers) {
            req.setRequestHeader(key, headers[key]);
          }
          req.send(file);
        });
      });
  });
}

function showFormUpload(form) {
  setProgress(form, 0);
  form.querySelectorAll("input").forEach((el) => {
    el.disabled = true;
  });
  form.querySelector(".file-upload-start").classList.add("d-none");
  form.querySelector(".file-upload-busy").classList.remove("d-none");
}

function hideFormUpload(form) {
  setProgress(form, 0);
  form.querySelectorAll("input").forEach((el) => {
    el.disabled = false;
    // important to retrigger "change" when someone enters the same file name again
    el.value = "";
  });
  form.querySelector(".file-upload-start").classList.remove("d-none");
  form.querySelector(".file-upload-busy").classList.add("d-none");
}

function setProgress(form, pct) {
  let pb = form.querySelector(".progress-bar");
  pb.setAttribute("style", "width: " + pct + "%");
  pb.setAttribute("aria-valuenow", pct);
  let pi = form.querySelector(".progress-bar-percent");
  pi.innerText = pct;
}
