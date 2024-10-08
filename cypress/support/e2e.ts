import "cypress-common";
import { type CyHttpMessages } from "cypress/types/net-stubbing";

import "./commands";
import "./queries";

Cypress.Keyboard.defaults({ keystrokeDelay: 0 });

before(() => {
  // Make sure to install the capture script at the earliest possible point,
  // otherwise requests that require the token within before() hooks won't work
  installCSRFTokenCaptureScript();
});

beforeEach(() => {
  // Make sure we don't install the capture script twice for any given test.
  if (!cy.state("aliases")?.captureCSRFTokenInstalled) {
    installCSRFTokenCaptureScript();
  }

  // Keep polling requests from Cypress command log
  cy.intercept({ method: "GET", url: /^\/static\// }, { log: false });
  cy.intercept({ method: "GET", url: /^\/\w+-icon/ }, { log: false });
  cy.intercept({ method: "GET", url: /^\/static\// }, { log: false });
});

function installCSRFTokenCaptureScript() {
  cy.intercept(
    /^\/(?!_templ\/|static\/|\w+-icon).*$/,
    { middleware: true },
    (req) => {
      req.on("after:response", (req) => {
        return extractCSRFTokenFromResponse(req);
      });
    },
  );

  const aliases = cy.state("aliases") || {};
  aliases.captureCSRFTokenInstalled = true;
  cy.state("aliases", aliases);
}

export function extractCSRFTokenFromResponse(
  res: Cypress.Response<unknown> | CyHttpMessages.IncomingResponse,
) {
  if (typeof res.body === "string") {
    const match = res.body.match(
      /<meta name="csrf-token" content="(?<csrfToken>[^"]+)"\/?>/,
    );
    if (match) {
      const { csrfToken } = match?.groups;

      if (csrfToken) {
        // Cannot easily store aliased value from an event handler so we keep it on the test context
        const ctx = cy.state("ctx");
        if (csrfToken && csrfToken !== ctx.CSRFToken) {
          ctx.CSRFToken = csrfToken;
        }
      }
    }
  }
}
