import * as fs from 'fs';

export class Config {
  static checkEnvFile() {
    if (!fs.existsSync('./.env')) {
      throw new ConfigError(
        '".env" ファイルがありません。".env.sample" をコピーして、".env" ファイルを作成してください。',
      );
    }
  }

  static getEnvValue(envName: string) {
    const value = process.env[envName];

    if (value === '' || typeof value === 'undefined') {
      throw new ConfigError(`".env"ファイルに "${envName}" が記載がないまたは空白です。`);
    }

    return value;
  }

  static getAisegIpAddress(envName = 'AISEG2_IP_ADDRESS') {
    return Config.getEnvValue(envName);
  }

  static getAisegUser(envName = 'AISEG2_USER') {
    return Config.getEnvValue(envName);
  }

  static getAisegPassword(envName = 'AISEG2_PASSWORD') {
    return Config.getEnvValue(envName);
  }

  static getChromePath(envName = 'CHROME_PATH') {
    const value = process.env[envName];
    return value === '' ? undefined : value;
  }

  static getHeadlessMode(envName = 'DEBUG_IS_HEADLESS'): boolean | 'new' {
    const value = process.env[envName];
    return value !== 'TRUE' ? 'new' : false;
  }
}

class ConfigError extends Error {
  static {
    this.prototype.name = 'ConfigError';
  }
}
