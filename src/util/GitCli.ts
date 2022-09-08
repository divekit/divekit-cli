import shell from 'shelljs'
import { Logger } from 'tslog'
import { rootLogger } from './Logger'

export class GitCli {
  log: Logger = rootLogger.getChildLogger()
  private readonly prefix: string
  errorsOccurred = false

  constructor(prefix?: string) {
    this.log.settings.name = 'GitCli'
    if (prefix) this.prefix = prefix
    else this.prefix = 'https://github.com/divekit/'
  }

  clone(gitHubName: string, localName: string, isTool = false): boolean {
    let subDir = ''
    if (isTool) subDir = 'tools/'

    const command = 'git clone -q ' + this.prefix + gitHubName + '.git .divekit/' + subDir + localName

    this.log.debug('exec shell command: ' + command)
    const successful = shell.exec(command).code === 0
    this.errorsOccurred = this.errorsOccurred && successful
    return successful
  }
}
