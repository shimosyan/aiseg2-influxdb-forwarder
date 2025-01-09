import dayjs from 'dayjs';
import { AiSEG2 } from './AiSEG2';
import { Config } from './Config';
import { Influx } from './Influx';

Config.checkEnvFile();

const aiseg2Host = Config.getAisegHost();
const aiseg2User = Config.getAisegUser();
const aiseg2Password = Config.getAisegPassword();
const aiseg2UseHTTPS = Config.getAisegUseHTTPS();

const influxdbHost = Config.getInfluxdbHost();
const influxdbToken = Config.getInfluxdbToken();
const influxdbOrg = Config.getInfluxdbOrg();
const influxdbBucket = Config.getInfluxdbBucket();
const influxdbUseHTTPS = Config.getInfluxdbUseHTTPS();

console.log('aiseg2Host', aiseg2Host);
console.log('aiseg2User', aiseg2User);
console.log('aiseg2UseHTTPS', aiseg2UseHTTPS);
console.log('influxdbHost', influxdbHost);
console.log('influxdbOrg', influxdbOrg);
console.log('influxdbBucket', influxdbBucket);
console.log('influxdbUseHTTPS', influxdbUseHTTPS);

async function run() {
  async function main(now = dayjs()) {
    // AiSEG2 からデータを取得
    const aiseg2 = new AiSEG2(aiseg2Host, aiseg2User, aiseg2Password, aiseg2UseHTTPS);

    const powerSummary = await aiseg2.getPowerSummary();
    console.log(now.format('YYYY-MM-DD HH:mm:ss'), 'powerSummary', powerSummary);

    const detailsUsagePower = await aiseg2.getDetailsUsagePower();
    console.log(now.format('YYYY-MM-DD HH:mm:ss'), 'detailsUsagePower', detailsUsagePower);

    // influxdb へデータを送信
    const influx = new Influx(
      influxdbHost,
      influxdbToken,
      influxdbOrg,
      influxdbBucket,
      influxdbUseHTTPS,
    );
    influx.writePower(powerSummary, detailsUsagePower);
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
