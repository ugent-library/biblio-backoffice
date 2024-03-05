import { defineConfig } from "cypress";
import * as dotenvPlugin from "cypress-dotenv";

export default defineConfig({
  projectId: "mjg74d",

  e2e: {
    baseUrl: "https://backoffice.bibliodev.ugent.be",
    experimentalStudio: true,
    experimentalRunAllSpecs: true,

    // Increase viewport width because GitHub Actions may render a wider font which
    // may cause button clicks to be prevented by overlaying elements.
    viewportWidth: 1200,

    setupNodeEvents(_on, config) {
      config = dotenvPlugin(config);

      return config;
    },
  },
});
