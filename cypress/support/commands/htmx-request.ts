import { extractCSRFTokenFromResponse } from "support/e2e";
import { getCSRFToken } from "support/util";

const NO_LOG = { log: false };

export default function htmxRequest(
  options: Partial<Cypress.RequestOptions>,
): Cypress.Chainable<Cypress.Response<string>> {
  if (!getCSRFToken()) {
    // Load home page once, this will capture the CSRF token
    cy.request("/", NO_LOG).then(extractCSRFTokenFromResponse);
  }

  // cy.then is necessary in case CSRFToken was only just loaded during cy.request("/")
  return cy
    .then(() => doRequest(options, false))
    .then((response) => {
      if (response.isOkStatusCode) {
        const $partial = Cypress.$(response.body);
        const $alert = $partial.find(".alert-danger");

        if ($alert.length) {
          throw new Error(
            `Error during backend request: ${options.method || "GET"} ${options.url}` +
              "\n\n" +
              $alert
                .find("li")
                .map((_, li) => `- ${li.textContent}`)
                .get()
                .join("\n"),
          );
        }

        return cy.wrap(response, NO_LOG);
      }

      if (
        typeof response.body === "string" &&
        response.body.includes("Forbidden - CSRF token") // ... invalid | ... not found
      ) {
        // Load home page to get fresh CSRFToken and try again
        return cy
          .request<string>("/", NO_LOG)
          .then(extractCSRFTokenFromResponse)
          .then(() => doRequest(options, true));
      }

      throw Error(
        `Error "${response.status} ${response.statusText}" during backend request: ${options.method || "GET"} ${options.url}`,
      );
    });
}

function doRequest(
  options: Partial<Cypress.RequestOptions>,
  failOnStatusCode: boolean,
) {
  // Deep clone options first so we don't alter it
  const clonedOptions = structuredClone(options);

  clonedOptions.failOnStatusCode = failOnStatusCode;

  if (!clonedOptions.headers) {
    clonedOptions.headers = {};
  }

  if (!clonedOptions.headers["X-CSRF-Token"]) {
    const csrfToken = getCSRFToken();
    if (!csrfToken) {
      throw new Error("CSRF token has not been set");
    }

    clonedOptions.headers["X-CSRF-Token"] = csrfToken;
  }

  if (typeof clonedOptions.log === "undefined") {
    // Do not log this request, unless specified otherwise
    clonedOptions.log = NO_LOG.log;
  }

  return cy.request<string>(clonedOptions);
}

declare global {
  namespace Cypress {
    interface Chainable {
      /**
       * cy.htmxRequest() is a convenience command that will deal with the CSRF token for making requests in several ways:
       * - If the CSRF token is not yet available, it will be loaded first via cy.request("/")
       * - If the CSRF token appears to be invalid (), it will be refreshed (also via cy.request("/"))
       */
      htmxRequest(
        options: Partial<RequestOptions>,
      ): Chainable<Response<string>>;
    }
  }
}
