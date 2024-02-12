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

  static getAisegHost(envName = 'AISEG2_HOST') {
    return Config.getEnvValue(envName);
  }

  static getAisegUser(envName = 'AISEG2_USER') {
    return Config.getEnvValue(envName);
  }

  static getAisegPassword(envName = 'AISEG2_PASSWORD') {
    return Config.getEnvValue(envName);
  }

  static getAisegUseHTTPS(envName = 'AISEG2_USE_HTTPS') {
    return process.env[envName] === '1';
  }

  static getInfluxdbHost(envName = 'INFLUXDB_HOST') {
    return Config.getEnvValue(envName);
  }

  static getInfluxdbToken(envName = 'INFLUXDB_TOKEN') {
    return Config.getEnvValue(envName);
  }

  static getInfluxdbOrg(envName = 'INFLUXDB_ORG') {
    return Config.getEnvValue(envName);
  }

  static getInfluxdbBucket(envName = 'INFLUXDB_BUCKET') {
    return Config.getEnvValue(envName);
  }
}

class ConfigError extends Error {
  static {
    this.prototype.name = 'ConfigError';
  }
}
