import { logCommand } from "./helpers";

const NO_LOG = { log: false };

const METHODS = {
  "hx-get": "GET",
  "hx-post": "POST",
  "hx-put": "PUT",
  "hx-delete": "DELETE",
} as const;

type HtmxMethod = keyof typeof METHODS;

export default function triggerHtmx(
  subject: JQuery<HTMLElement>,
  method: HtmxMethod,
): Cypress.Chainable<Cypress.Response<string>> {
  const url = subject.attr(method);
  if (!url) {
    throw new Error(`Could not find '${method}' attribute on subject`);
  }

  const hxHeaders = JSON.parse(subject.attr(`hx-headers`) ?? "null");

  const log = logCommand(
    "triggerHtmx",
    { method, url, "hx-headers": hxHeaders },
    method,
  );

  return cy
    .htmxRequest({
      url,
      method: METHODS[method],
      headers: {
        ...hxHeaders,
      },
      ...NO_LOG,
    })
    .finishLog(log);
}

declare global {
  namespace Cypress {
    interface Chainable {
      triggerHtmx(method: HtmxMethod): Chainable<Response<string>>;
    }
  }
}
