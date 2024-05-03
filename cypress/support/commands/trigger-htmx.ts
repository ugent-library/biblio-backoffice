import { logCommand } from "./helpers";

const NO_LOG = { log: false };

const METHODS = {
  "hx-get": "GET",
  "hx-post": "POST",
  "hx-put": "PUT",
  "hx-delete": "DELETE",
} as const;

type HtmxMethod = keyof typeof METHODS;

export default function triggerHtmx<T = unknown>(
  subject: JQuery<HTMLElement>,
  method: HtmxMethod,
): Cypress.Chainable<Cypress.Response<T>> {
  const url = subject.attr(method);
  if (!url) {
    throw new Error(`Could not find '${method}' attribute on subject`);
  }

  const hxHeaders = JSON.parse(subject.attr(`hx-headers`) ?? "null");

  logCommand("triggerHtmx", { method, url, "hx-headers": hxHeaders }, method);

  return cy
    .document(NO_LOG)
    .find<HTMLMetaElement>("meta[name=csrf-token]", NO_LOG)
    .then((csrfToken) => {
      return cy.request<T>({
        url,
        method: METHODS[method],
        headers: {
          ...hxHeaders,
          "X-CSRF-Token": csrfToken.prop("content"),
        },
      });
    });
}

declare global {
  namespace Cypress {
    interface Chainable {
      triggerHtmx<T = unknown>(method: HtmxMethod): Chainable<Response<T>>;
    }
  }
}
