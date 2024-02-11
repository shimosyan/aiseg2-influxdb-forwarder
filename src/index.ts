import dayjs from 'dayjs';
import { AiSEG2 } from './AiSEG2';
import { Browser } from './Browser';
import { Config } from './Config';

Config.checkEnvFile();

const aiseg2IpAddress = Config.getAisegIpAddress();
const aiseg2User = Config.getAisegUser();
const aiseg2Password = Config.getAisegPassword();

const chromePath = Config.getChromePath();
const debugHeadlessMode = Config.getHeadlessMode();

console.log('aiseg2IpAddress', aiseg2IpAddress);
console.log('aiseg2User', aiseg2User);
console.log('chromePath', chromePath);
console.log('debugHeadlessMode', debugHeadlessMode);

const browser = new Browser(chromePath);
browser.setDebugOption(debugHeadlessMode);

async function run() {
  await browser.launch();

  async function main(now = dayjs()) {
    const page = await browser.initPage();
    const aiseg2 = new AiSEG2(page, aiseg2IpAddress, aiseg2User, aiseg2Password);

    const powerSummary = await aiseg2.getPowerSummary();

    await browser.closePage(page);

    console.log(now.format('YYYY-MM-DD HH:mm:ss'), powerSummary);
  }

  async function interval(microSeconds: number) {
    for (;;) {
      await new Promise((resolve) => setTimeout(resolve, microSeconds));
      await main();
    }
  }

  //await main();
  await interval(5000);
  await browser.close();
}

run();
