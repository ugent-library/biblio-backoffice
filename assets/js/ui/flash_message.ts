import { Toast, type ToastOptions } from "bootstrap.native";

const levelMap = {
  success: "success",
  info: "primary",
  warning: "warning",
  error: "error",
} as const;

type Level = keyof typeof levelMap;

type FlashMessageOptions = {
  level?: Level;
  isLoading?: boolean;
  title?: string;
  text: string;
  isDismissible?: boolean;
  toastOptions?: Partial<ToastOptions>;
};

export default function showFlashMessage({
  level,
  isLoading = false,
  title,
  text,
  isDismissible = true,
  toastOptions,
}: FlashMessageOptions) {
  let flashMessages = document.querySelector("#flash-messages");
  if (flashMessages) {
    /*
     * Expects somewhere in the document ..
     *
     * <template id="template-flash-message"></template>
     *
     * .. a template that encapsulates the toast
     * */
    const templateFlashMessage = document.querySelector<HTMLTemplateElement>(
      "template#template-flash-message",
    );

    if (templateFlashMessage) {
      const toastEl = (
        templateFlashMessage.content.cloneNode(true) as HTMLElement
      ).querySelector<HTMLElement>(".toast");

      if (toastEl) {
        const flashMessage = new FlashMessage(toastEl);

        flashMessage.setLevel(level);
        flashMessage.setIsLoading(isLoading);
        flashMessage.setTitle(title);
        flashMessage.setText(text);
        flashMessage.setIsDismissible(isDismissible);

        flashMessages.appendChild(toastEl);

        flashMessage.show(toastOptions);

        return flashMessage;
      }
    }
  }
}

class FlashMessage {
  #toastEl: HTMLElement;
  #toast: Toast | null;

  constructor(readonly toastEl: HTMLElement) {
    this.#toastEl = toastEl;
  }

  setLevel(level: Level | undefined) {
    this.#toastEl.querySelectorAll(".toast-body i.if").forEach((el) => {
      el.classList.add("d-none");
    });

    if (level && Object.keys(levelMap).includes(level)) {
      this.#toastEl
        .querySelector(`.if--${levelMap[level]}`)
        ?.classList.remove("d-none");
    }
  }

  setIsLoading(isLoading: boolean) {
    this.#toastEl
      .querySelector(".spinner-border")
      ?.classList.toggle("d-none", !isLoading);
  }

  setTitle(title: string | undefined) {
    const titleEl = this.#toastEl.querySelector(".alert-title");
    if (titleEl) {
      if (title) {
        titleEl.classList.remove("d-none");
        titleEl.textContent = title;
      } else {
        titleEl.classList.add("d-none");
      }
    }
  }

  setText(text: string) {
    const textEl = this.#toastEl.querySelector(".toast-text");
    if (textEl) {
      textEl.innerHTML = text;
    }
  }

  setIsDismissible(isDismissible: boolean) {
    const btnClose = this.#toastEl.querySelector(".btn-close");
    if (btnClose) {
      btnClose.classList.toggle("d-none", !isDismissible);
    }
  }

  setAutohide(autohide: boolean, delay = 5000) {
    if (this.#toast) {
      this.#toast.options.autohide = autohide;
      this.#toast.options.delay = delay;
      this.#toast.show();
    }
  }

  show(toastOptions: Partial<ToastOptions> = {}) {
    this.#toast =
      Toast.getInstance(this.#toastEl) ??
      new Toast(this.#toastEl, toastOptions);

    this.#toast.show();

    this.#toastEl.addEventListener("hidden.bs.toast", () => {
      this.#toastEl.remove();
    });
  }

  hide() {
    if (this.#toast) {
      this.#toast.hide();
    }
  }
}
