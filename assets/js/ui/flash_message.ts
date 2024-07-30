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

export default class FlashMessage {
  private toastEl: HTMLElement;
  private toast: Toast;

  constructor({
    level,
    isLoading = false,
    title,
    text,
    isDismissible = true,
    toastOptions,
  }: FlashMessageOptions) {
    this.initFlashMessage();

    this.toast = new Toast(this.toastEl, toastOptions);

    this.setLevel(level);
    this.setIsLoading(isLoading);
    this.setTitle(title);
    this.setText(text);
    this.setIsDismissible(isDismissible);
  }

  setLevel(level: Level | undefined) {
    this.toastEl.querySelectorAll(".toast-body i.if").forEach((el) => {
      el.classList.add("d-none");
    });

    if (level && Object.keys(levelMap).includes(level)) {
      this.query(`.if--${levelMap[level]}`).classList.remove("d-none");
    }
  }

  setIsLoading(isLoading: boolean) {
    this.query(".spinner-border").classList.toggle("d-none", !isLoading);
  }

  setTitle(title: string | undefined) {
    const titleEl = this.query(".alert-title");
    if (title) {
      titleEl.classList.remove("d-none");
      titleEl.textContent = title;
    } else {
      titleEl.classList.add("d-none");
    }
  }

  setText(text: string) {
    this.query(".toast-text").innerHTML = text;
  }

  setIsDismissible(isDismissible: boolean) {
    this.query(".btn-close").classList.toggle("d-none", !isDismissible);
  }

  setAutohide(autohide: boolean, delay = 5000) {
    this.toast.options.autohide = autohide;
    this.toast.options.delay = delay;
    this.toast.show();
  }

  show() {
    this.toast.show();

    this.toastEl.addEventListener("hidden.bs.toast", () => {
      this.toastEl.remove();
    });
  }

  hide() {
    // For some reason, BSN doesn't show the toast for 17ms, so we wait 20ms first before trying to hide to prevent a race condition.
    // https://github.com/thednp/bootstrap.native/blob/master/src/components/toast.ts#L140
    setTimeout(() => {
      this.toast.hide();
    }, 20);
  }

  private initFlashMessage() {
    const flashMessages = document.querySelector("#flash-messages");
    if (!flashMessages) {
      throw new Error("Container for flash messages not found.");
    }

    const templateFlashMessage = document.querySelector<HTMLTemplateElement>(
      "template#template-flash-message",
    );
    if (!templateFlashMessage) {
      throw new Error("Template for flash messages not found.");
    }

    const toastFragment = templateFlashMessage.content.cloneNode(
      true,
    ) as HTMLElement;
    const toastEl = toastFragment.querySelector<HTMLElement>(".toast");
    if (!toastEl) {
      throw new Error(
        "Template for flash messages does not contain a '.toast' element.",
      );
    }

    flashMessages.appendChild(toastEl);

    this.toastEl = toastEl;
  }

  private query(selector: string) {
    if (!this.toastEl) {
      throw new Error("FlashMessage is not yet initialized.");
    }

    const el = this.toastEl.querySelector(selector);

    if (!el) {
      throw new Error(`Element not found: ${selector}`);
    }

    return el;
  }
}
