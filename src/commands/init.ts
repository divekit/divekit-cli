import type { Arguments, CommandBuilder } from 'yargs'
import * as fs from 'fs'
import { GitCli } from '../util/GitCli'
import { CliVisuals } from '../util/CliVisual'
import { rootLogger } from '../util/Logger'

type Options = {
  advanced: boolean | undefined
}

export const command = 'init'
export const desc = 'initialize divekit stuff'

export const builder: CommandBuilder<Options, Options> = yargs =>
  yargs.options({
    advanced: { type: 'boolean' }
  })

export const handler = (argv: Arguments<Options>): void => {
  if (fs.existsSync('.divekit')) {
    rootLogger.error('.divekit directory already existing.')
    CliVisuals.printRedText('Could not initialize Divekit.')
    process.exit(1)
  }

  const { advanced } = argv
  if (!advanced) rootLogger.info('Initialize with standard configuration')

  fs.mkdirSync('.divekit')
  fs.mkdirSync('.divekit/tools')

  const git = new GitCli()
  git.clone('divekit-config', 'config')
  git.clone('Access-Manager-2.0', 'access-manager', true)

  if (git.errorsOccurred) {
    rootLogger.error('Errors occurred while initializing pls contact your local admin.')
    process.exit(1)
  }

  CliVisuals.printGreenText('Divekit successfully initialized')
  process.exit(0)
}
