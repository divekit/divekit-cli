import { clear } from 'console'
import * as figlet from 'figlet'
import chalk from 'chalk'

export class CliVisuals {
  static printHeader(): void {
    clear()
    const divekitASCII = figlet.textSync('divekit', {
      horizontalLayout: 'full'
    })
    console.log(chalk.green(divekitASCII))
  }

  static printGreenText(msg: string): void {
    console.log(chalk.green(msg))
  }

  static printRedText(msg: string): void {
    console.log(chalk.red(msg))
  }

  static printHelp(): void {
    console.log('No arguments given. Get help elsewhere pls.')
  }
}
