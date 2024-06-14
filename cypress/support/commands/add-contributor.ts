import { logCommand } from "./helpers";

const NO_LOG = { log: false };

const CONTRIBUTOR_MAP = {
  author: "Authors",
  editor: "Editors",
  supervisor: "Supervisors",
  creator: "Creators",
} as const;

type ContributorType = keyof typeof CONTRIBUTOR_MAP;

export function addAuthor(
  firstName: string,
  lastName: string,
  external = false,
  role?: string,
): void {
  addContributor("author", firstName, lastName, external, role);
}

export function addEditor(
  firstName: string,
  lastName: string,
  external = false,
): void {
  addContributor("editor", firstName, lastName, external);
}

export function addSupervisor(
  firstName: string,
  lastName: string,
  external = false,
): void {
  addContributor("supervisor", firstName, lastName, external);
}

export function addCreator(
  firstName: string,
  lastName: string,
  external = false,
): void {
  addContributor("creator", firstName, lastName, external);
}

function addContributor(
  contributorType: ContributorType,
  firstName: string,
  lastName: string,
  external: boolean,
  role?: string,
): void {
  const consoleProps = {
    "Contributor type": contributorType,
    "First name": firstName,
    "Last name": lastName,
    "External contributor": external,
  };

  if (contributorType === "author") {
    consoleProps["Role"] = role;
  }

  logCommand("add" + Cypress._.capitalize(contributorType), consoleProps, [
    `${firstName} ${lastName} ${external ? "(external)" : ""}`.trim(),
  ]);

  cy.updateFields(
    CONTRIBUTOR_MAP[contributorType],
    () => {
      cy.intercept("/+(publication|dataset)/*/contributors/*/suggestions?*").as(
        "suggestContributor",
      );

      if (firstName) {
        cy.setFieldByLabel("First name", firstName);
        cy.wait("@suggestContributor", NO_LOG);
      }

      if (lastName) {
        cy.setFieldByLabel("Last name", lastName);
        cy.wait("@suggestContributor", NO_LOG);
      }

      cy.contains(
        ".btn",
        `Add ${external ? "external " : ""}${contributorType}`,
        NO_LOG,
      ).click(NO_LOG);
    },
    true,
  );
}

declare global {
  namespace Cypress {
    interface Chainable {
      addAuthor(
        firstName: string,
        lastName: string,
        external?: boolean,
        role?: string,
      ): Chainable<void>;

      addEditor(
        firstName: string,
        lastName: string,
        external?: boolean,
      ): Chainable<void>;

      addSupervisor(
        firstName: string,
        lastName: string,
        external?: boolean,
      ): Chainable<void>;

      addCreator(
        firstName: string,
        lastName: string,
        external?: boolean,
      ): Chainable<void>;
    }
  }
}
