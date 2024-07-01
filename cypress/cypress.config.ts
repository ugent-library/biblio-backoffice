import { defineConfig } from "cypress";

export default defineConfig({
  projectId: "mjg74d",

  env: {
    OIDC_ORIGIN: "http://localhost:3002",
    ELASTICSEARCH_ORIGIN: "http://localhost:3061",
    RESEARCHER_USER_NAME: "biblio-researcher",
    RESEARCHER_NAME: "Biblio Researcher",
    LIBRARIAN_USER_NAME: "biblio-librarian",
    LIBRARIAN_NAME: "Biblio Librarian",
  },

  e2e: {
    baseUrl: "http://localhost:3001",

    experimentalStudio: true,
    experimentalRunAllSpecs: true,

    // Increase viewport width because GitHub Actions may render a wider font which
    // may cause button clicks to be prevented by overlaying elements.
    viewportWidth: 1200,
  },
});
