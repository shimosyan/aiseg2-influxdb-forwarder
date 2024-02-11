import * as Puppeteer from 'puppeteer';

export class Browser {
  private debugHeadlessMode: boolean | 'new' | undefined = 'new';

  private browser?: Puppeteer.Browser;

  constructor(private chromePath: string | undefined = undefined) {}

  public setDebugOption(debugHeadlessMode: boolean | 'new') {
    this.debugHeadlessMode = debugHeadlessMode;
  }

  async launch() {
    this.browser = await Puppeteer.launch({
      headless: this.debugHeadlessMode,
      executablePath: this.chromePath,
      args: ['--no-sandbox', '--disable-setuid-sandbox'],
    });
  }

  async initPage(): Promise<Puppeteer.Page> {
    const page = await this.browser!.newPage();
    await page.setViewport({ width: 720, height: 480 });

    return page;
  }

  async closePage(page: Puppeteer.Page) {
    page.close();
  }

  async close() {
    await this.browser!.close();
  }
}
