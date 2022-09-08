#!/usr/bin/env node

import { CliVisuals } from './util/CliVisual'
import { hideBin } from 'yargs/helpers'
import yargs from 'yargs'

CliVisuals.printHeader()

if (!process.argv.slice(2).length) CliVisuals.printHelp()

// todo add -debug flag for log level configuration

yargs(hideBin(process.argv))
  .commandDir('commands') // Use the 'commands' directory to scaffold.
  .strict() // Enable strict mode.
  .alias({ h: 'help' }).argv // Useful aliases.
