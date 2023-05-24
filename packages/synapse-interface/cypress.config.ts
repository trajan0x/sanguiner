import { defineConfig } from 'cypress'
import synpressPlugins from '@synthetixio/synpress/plugins'

export default defineConfig({
  e2e: {
    baseUrl: 'http://localhost:3000/',
    setupNodeEvents: (on, config) => {
      synpressPlugins(on, config)
    },
    supportFile: 'cypress/support/e2e.ts',
  },
  video: false,
  trashAssetsBeforeRuns: true,
  screenshotOnRunFailure: true,
  screenshotsFolder: 'cypress/visual-states/current-screenshots',
  videosFolder: 'cypress/visual-states/current-videos',
});
