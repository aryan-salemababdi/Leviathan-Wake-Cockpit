import { Injectable, OnModuleInit, Logger, Inject } from '@nestjs/common';
import { ConfigType } from '@nestjs/config';
import Redis from 'ioredis';
import keydbConfig from '../../config/keydb/keydb.config';

@Injectable()
export class KeydbService implements OnModuleInit {
  private readonly logger = new Logger(KeydbService.name);
  public client: Redis;

  constructor(
    @Inject(keydbConfig.KEY)
    private readonly config: ConfigType<typeof keydbConfig>,
  ) {}

  async onModuleInit() {
    this.client = new Redis({
      host: this.config.host,
      port: this.config.port,
    });

    this.client.on('connect', () => {
      this.logger.log(
        `Successfully connected to KeyDB at ${this.config.host}:${this.config.port}`,
      );
    });

    this.client.on('error', (err) => {
      this.logger.error(`Could not connect to KeyDB: ${err}`);
    });
  }
}
