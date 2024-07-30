import FlashMessage from "./flash_message";

export default function fileDownload(el) {
  el.querySelectorAll("a[download]").forEach((anchor) => {
    anchor.addEventListener("click", async (e) => {
      e.preventDefault();

      const flashMessage = new FlashMessage({
        isLoading: true,
        text: "Preparing download...",
        isDismissible: false,
        toastOptions: {
          autohide: false,
        },
      });
      flashMessage.show();

      try {
        const response = await fetch(anchor.href);
        if (!response.ok) {
          throw new Error(
            "An unexpected error occurred while downloading the file. Please try again.",
          );
        }

        const filename = extractFileName(response);
        const blob = await response.blob();

        const dummyAnchor = document.createElement("a");
        dummyAnchor.href = URL.createObjectURL(blob);
        dummyAnchor.setAttribute("class", "link-primary text-nowrap");
        dummyAnchor.setAttribute("download", filename);
        dummyAnchor.textContent = filename;

        flashMessage.setLevel("success");
        flashMessage.setIsLoading(false);
        flashMessage.setText("Download ready: " + dummyAnchor.outerHTML);
        flashMessage.setIsDismissible(true);

        // Trigger download (save-as window or auto-download, depending on browser settings)
        dummyAnchor.click();
      } catch (error) {
        flashMessage.hide();

        new FlashMessage({
          level: "error",
          text: error.message,
        }).show();
      }
    });
  });
}

function extractFileName(response: Response) {
  const FILENAME_REGEX = /filename\*?=(UTF-8'')?/;
  const contentDispositionHeader = response.headers.get("content-disposition");

  if (contentDispositionHeader) {
    const fileNamePart = contentDispositionHeader
      .split(";")
      .find((n) => n.match(FILENAME_REGEX));

    if (fileNamePart) {
      const fileName = fileNamePart.replace(FILENAME_REGEX, "").trim();
      return decodeURIComponent(fileName);
    }
  }

  return "";
}
