import { extractHtmxJsonAttribute, extractSnapshotId } from "support/util";
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
  options?: AddAuthorOptions,
): void {
  addContributor("publication", "author", firstName, lastName, options);
}

export function addEditor(
  firstName: string,
  lastName: string,
  options?: AddAuthorOptions,
): void {
  addContributor("publication", "editor", firstName, lastName, options);
}

export function addSupervisor(
  firstName: string,
  lastName: string,
  options?: AddAuthorOptions,
): void {
  addContributor("publication", "supervisor", firstName, lastName, options);
}

export function addCreator(
  firstName: string,
  lastName: string,
  options?: AddAuthorOptions,
): void {
  addContributor("dataset", "creator", firstName, lastName, options);
}

type PostBody = (
  | {
      first_name: string;
      last_name: string;
    }
  | {
      id: string;
    }
) & {
  credit_role?: string[];
};

function addContributor(
  scope: "publication" | "dataset",
  contributorType: ContributorType,
  firstName: string,
  lastName: string,
  { external, role, biblioIdAlias = "@biblioId" }: AddAuthorOptions = {},
): void {
  cy.get<string>(biblioIdAlias, NO_LOG).then((biblioId) => {
    const log = prepareLog({
      biblioIdAlias,
      biblioId,
      contributorType,
      firstName,
      lastName,
      external,
      role,
    });

    // Dataset creators are in fact authors
    contributorType =
      contributorType === "creator" ? "author" : contributorType;

    const qs = {
      first_name: firstName,
      last_name: lastName,
    };

    let postBody = cy.wrap<PostBody>(qs, NO_LOG);

    if (!external) {
      // For UGent contributors, we need the person ID from the suggestions API
      postBody = cy
        .htmxRequest({
          url: `/${scope}/${biblioId}/contributors/${contributorType}/suggestions`,
          qs,
        })
        .then(extractHxValues);
    } else {
      // For external contributors, we can skip this and continue working with first name & last name as parameters
    }

    postBody.then((postBody) => {
      // First GET the .../confirm-create route to extract the snapshot ID for the POST in the next step
      cy.htmxRequest({
        url: `/${scope}/${biblioId}/contributors/${contributorType}/confirm-create`,
        qs: postBody,
      })
        .then(extractSnapshotId)
        .then((snapshotId) => {
          // Add the role (if applicable - only for publication authors)
          if (role) {
            postBody.credit_role = [role];
          }

          cy.htmxRequest({
            method: "POST",
            url: `/${scope}/${biblioId}/contributors/${contributorType}`,
            headers: {
              "If-Match": snapshotId,
            },
            form: true,
            body: postBody,
          });
        });
    });
  });
}

function prepareLog({
  biblioIdAlias,
  biblioId,
  contributorType,
  firstName,
  lastName,
  external,
  role,
}: {
  biblioIdAlias: Cypress.Alias;
  biblioId: string;
  contributorType: string;
  firstName: string;
  lastName: string;
  external: boolean;
  role: string;
}) {
  const consoleProps = {
    "Biblio ID alias": biblioIdAlias,
    "Biblio ID": biblioId,
    "Contributor type": contributorType,
    "First name": firstName,
    "Last name": lastName,
    "External contributor": external,
  };

  if (contributorType === "author") {
    consoleProps["Role"] = role;
  }

  const typeName = Cypress._.capitalize(contributorType);

  return logCommand(`add${typeName}`, consoleProps, [
    `${firstName || "[missing]"} ${lastName || "[missing]"} ${external ? "(external)" : ""}`.trim(),
  ]);
}

function extractHxValues(response: Cypress.Response<string>) {
  return extractHtmxJsonAttribute<{ id: string }>(
    response,
    "button.btn-primary",
    "hx-vals",
  );
}

type AddContributorOptions = {
  external?: boolean;
  biblioIdAlias?: Cypress.Alias;
};

type AddAuthorOptions = AddContributorOptions & {
  role?: string;
};

declare global {
  namespace Cypress {
    interface Chainable {
      addAuthor(
        firstName: string,
        lastName: string,
        options?: AddAuthorOptions,
      ): Chainable<void>;

      addEditor(
        firstName: string,
        lastName: string,
        options?: AddContributorOptions,
      ): Chainable<void>;

      addSupervisor(
        firstName: string,
        lastName: string,
        options?: AddContributorOptions,
      ): Chainable<void>;

      addCreator(
        firstName: string,
        lastName: string,
        options?: AddContributorOptions,
      ): Chainable<void>;
    }
  }
}
