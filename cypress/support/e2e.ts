import "cypress-common";

import "./commands";
import "./queries";

Cypress.Keyboard.defaults({ keystrokeDelay: 0 });

before(() => {
  // Make sure to install the capture script at the earliest possible point,
  // otherwise requests that require the token within before() hooks won't work
  installCSRFTokenCaptureScript();
});

beforeEach(() => {
  // Make sure we don't install the catpure script twice for any given test.
  if (!cy.state("aliases")?.captureCSRFTokenInstalled) {
    installCSRFTokenCaptureScript();
  }

  // Keep dashboard-icon polling from Cypress command log
  cy.intercept({ method: "GET", url: "/dashboard-icon" }, { log: false });
  cy.intercept(
    { method: "GET", url: "/candidate-records-icon" },
    { log: false },
  );
});

function installCSRFTokenCaptureScript() {
  cy.intercept("**/*", { middleware: true }, (req) => {
    req.on("after:response", (req) => {
      return extractCSRFTokenFromResponse(req);
    });
  });

  const aliases = cy.state("aliases") || {};
  aliases.captureCSRFTokenInstalled = true;
  cy.state("aliases", aliases);
}

export function extractCSRFTokenFromResponse(res) {
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
