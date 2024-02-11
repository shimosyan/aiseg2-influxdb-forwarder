import dayjs from 'dayjs';
import { AiSEG2 } from './AiSEG2';
import { Config } from './Config';

Config.checkEnvFile();

const aiseg2Host = Config.getAisegHost();
const aiseg2User = Config.getAisegUser();
const aiseg2Password = Config.getAisegPassword();

console.log('aiseg2Host', aiseg2Host);
console.log('aiseg2User', aiseg2User);

async function run() {
  async function main(now = dayjs()) {
    const aiseg2 = new AiSEG2(aiseg2Host, aiseg2User, aiseg2Password);

    const powerSummary = await aiseg2.getPowerSummary();

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
}

run();
