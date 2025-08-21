import { Injectable, Logger } from '@nestjs/common';
import { Cron, CronExpression } from '@nestjs/schedule';
import { KeydbService } from '../keydb/keydb.service';
// import { HttpService } from '@nestjs/axios';
import * as fs from 'fs';
import * as path from 'path';
import { TopTraders } from 'src/common/types/topTraders.type';

const WHITELIST_KEY = 'whale_whitelist';

@Injectable()
export class UpdaterService {
  private readonly logger = new Logger(UpdaterService.name);
  private wallets: any;

  constructor(
    private readonly keyDBService: KeydbService,
    // private readonly httpService: HttpService,
  ) {
    const filePath = path.join(process.cwd(), 'db.json');
    const raw = fs.readFileSync(filePath, 'utf-8');
    this.wallets = JSON.parse(raw);
  }

  @Cron(CronExpression.EVERY_6_HOURS)
  async handleUpdate() {
    this.logger.log('Scheduled whale data update process started...');

    try {
      const topTraders = await this.fetchTopTraders();
      this.logger.log(
        `Successfully fetched data ${topTraders.length} top traders`,
      );

      if (topTraders.length === 0) {
        this.logger.log('No addresses found, skipping update.');
        return;
      }

      const whiteList = topTraders
        .filter(
          (trader) =>
            trader.account_value > 5_000_000 &&
            trader.perp_alltime_pnl > 0 &&
            trader.perp_month_pnl > 0 &&
            trader.perp_week_pnl > 0 &&
            trader.perp_day_pnl > 0 &&
            trader.main_position.value / trader.account_value >= 0.5 &&
            (trader.direction_bias >= 70 || trader.direction_bias <= 30),
        )
        .map((trader) => trader.address);

      await this.saveWhiteListToKeyDB(whiteList);
      this.logger.log(
        `Whitelist successfully updated in KeyDB with ${whiteList.length} addresses`,
      );
    } catch (error) {
      this.logger.log(`Failed to update whale data: ${error.stack}`);
    }
  }

  private async fetchTopTraders(): Promise<TopTraders[]> {
    this.logger.log('Fetching data from external analytics source...');

    // Simulator Data
    const mockData = this.wallets;

    return Promise.resolve(mockData);
  }

  private async saveWhiteListToKeyDB(addresses: string[]): Promise<void> {
    if (addresses.length === 0) {
      this.logger.log('Whitelist is empty, skipping KeyDB update.');
      return;
    }

    const jsonData = JSON.stringify(addresses);
    await this.keyDBService.client.set(WHITELIST_KEY, jsonData);
  }
}
