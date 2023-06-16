import { defineConfig } from 'cypress'
import * as dotenvPlugin from 'cypress-dotenv'

export default defineConfig({
  projectId: 'mjg74d',

  e2e: {
    baseUrl: 'https://backoffice.bibliodev.ugent.be',
    experimentalStudio: true,

    setupNodeEvents(_on, config) {
      config = dotenvPlugin(config)

      return config
    },
  },
})
