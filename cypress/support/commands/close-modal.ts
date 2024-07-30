import { logCommand, updateConsoleProps } from "./helpers";

const NO_LOG = { log: false };

type CloseModalOptions = {
  log?: boolean;
};

export default function closeModal(
  subject: undefined | JQuery<HTMLElement>,
  save: boolean | string | RegExp = false,
  options: CloseModalOptions = { log: true },
): void {
  const dismissButtonText =
    typeof save === "boolean" ? (save ? /^\s*Save\s*$/ : "Cancel") : save;

  let log: Cypress.Log | null = null;
  if (options.log === true) {
    log = logCommand(
      "closeModal",
      {
        subject: subject && subject.get(0),
        "Dismiss button text": dismissButtonText,
      },
      dismissButtonText,
    );
    log.set("type", !subject ? "parent" : "child");
  }

  const doCloseModal = () => {
    cy.contains(".modal-footer .btn", dismissButtonText, NO_LOG)
      .then(($el) => {
        if (options.log === true) {
          log.set("$el", $el);
          updateConsoleProps(log, (cp) => {
            cp["Button element"] = $el.get(0);
          });
        }

        return $el;
      })
      .click({ ...NO_LOG, animationDistanceThreshold: 1 });
  };

  if (subject) {
    cy.wrap(subject, NO_LOG).within(NO_LOG, doCloseModal);
  } else {
    doCloseModal();
  }
}

declare global {
  namespace Cypress {
    interface Chainable {
      closeModal(
        save: boolean | string | RegExp,
        options?: CloseModalOptions,
      ): Chainable<void>;
    }
  }
}
