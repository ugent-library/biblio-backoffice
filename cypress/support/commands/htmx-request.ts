const NO_LOG = { log: false };

export default function htmxRequest<T>(
  options: Partial<Cypress.RequestOptions>,
): Cypress.Chainable<Cypress.Response<T>> {
  if (!getCSRFToken()) {
    // Load home page once, this will capture the CSRF token
    cy.visit("/", { log: false });
  }

  // cy.then is necessary in case CSRFToken was only just loaded during cy.visit("/")
  return cy
    .then(() => doRequest<T>(options, false))
    .then((response) => {
      if (response.isOkStatusCode) {
        return cy.wrap(response, NO_LOG);
      }

      if (
        typeof response.body === "string" &&
        response.body.includes("Forbidden - CSRF token invalid")
      ) {
        // Load home page to get fresh CSRFToken and try again
        return cy
          .visit("/", { log: false })
          .then(() => doRequest<T>(options, true));
      }

      return cy.wrap(response, NO_LOG);
    });
}

function doRequest<T>(
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
    clonedOptions.headers["X-CSRF-Token"] = getCSRFToken();
  }

  if (typeof clonedOptions.log === "undefined") {
    // Do not log this request, unless specified otherwise
    clonedOptions.log = false;
  }

  return cy.request<T>(clonedOptions);
}

declare global {
  namespace Cypress {
    interface Chainable {
      /**
       * cy.htmxReques() is a convenience command that will deal with the CSRF token for making requests in several ways:
       * - If the CSRF token is not yet available, it will be loaded first via cy.visit("/")
       * - If the CSRF token appears to be invalid (), it will be refreshed (also via cy.visit("/"))
       */
      htmxRequest<T>(options: Partial<RequestOptions>): Chainable<Response<T>>;
    }
  }
}

function getCSRFToken() {
  const ctx = cy.state<{ CSRFToken?: string }>("ctx");

  return ctx.CSRFToken;
}
