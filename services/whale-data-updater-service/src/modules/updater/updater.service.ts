import { Injectable, Logger } from '@nestjs/common';
import { Cron, CronExpression } from '@nestjs/schedule';
import { KeydbService } from '../keydb/keydb.service';
import { HttpService } from '@nestjs/axios';

const WHITELIST_KEY = 'whale_whitelist';

@Injectable()
export class UpdaterService {
  private readonly logger = new Logger(UpdaterService.name);

  constructor(
    private readonly keyDBService: KeydbService,
    private readonly httpService: HttpService,
  ) {}

  @Cron(CronExpression.EVERY_10_MINUTES)
  async handleUpdate() {
    this.logger.log('Scheduled whale data update process started...');

    try {
      const topTraders = await this.fetchTopTraders();
      this.logger.log(
        `Successfully fetched data ${topTraders.length} top traders`,
      );

      const whiteList = topTraders.map((trader) => trader.address);
      if (whiteList.length === 0) {
        this.logger.log('No addresses found, skipping update.');
        return;
      }

      await this.saveWhiteListToKeyDB(whiteList);
      this.logger.log(
        `Whitelist successfully updated in KeyDB with ${whiteList.length} addresses`,
      );
    } catch (error) {
      this.logger.log(`Failed to update whale data: ${error.stack}`);
    }
  }

  private async fetchTopTraders(): Promise<{ address: string; pnl: number }[]> {
    this.logger.log('Fetching data from external analytics source...');
    // const url = 'https://hyperdash.info/api/top-traders';
    // const response = await firstValueFrom(this.httpService.get(url));
    // return response.data;

    // Simulator Data
    const mockData = [
      { address: '0x73BCEb1Cd57C711feCac222608109E6A6634550A', pnl: 5200000 },
      { address: '0x4862733B5FdDFd35f35ea8CCfB05490c56D4f302', pnl: 4100000 },
      { address: '0x4a1D225d32c4B173628d1D472c2191924B279282', pnl: 3850000 },
    ];

    return Promise.resolve(mockData);
  }

  private async saveWhiteListToKeyDB(address: string[]): Promise<void> {
    const jsonData = JSON.stringify(address);
    await this.keyDBService.client.set(WHITELIST_KEY, jsonData);
  }
}
