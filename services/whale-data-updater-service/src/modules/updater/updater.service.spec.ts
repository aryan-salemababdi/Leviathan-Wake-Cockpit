import { Test, TestingModule } from '@nestjs/testing';
import { UpdaterService } from './updater.service';
import { KeydbService } from '../keydb/keydb.service';

describe('UpdaterService', () => {
  let service: UpdaterService;
  let keyDBService: Partial<KeydbService>;

  const mockWallets = [
    {
      address: '0x1',
      account_value: 10_000_000,
      main_position: { coin: 'BTC', value: 6_000_000, side: 'LONG' },
      direction_bias: 80,
      perp_day_pnl: 1000,
      perp_week_pnl: 5000,
      perp_month_pnl: 10000,
      perp_alltime_pnl: 50000,
    },
    {
      address: '0x2',
      account_value: 3_000_000,
      main_position: { coin: 'ETH', value: 2_000_000, side: 'LONG' },
      direction_bias: 50,
      perp_day_pnl: 500,
      perp_week_pnl: 2000,
      perp_month_pnl: 4000,
      perp_alltime_pnl: 10000,
    },
  ];

  beforeEach(async () => {
    keyDBService = {
      client: {
        set: jest.fn().mockResolvedValue('OK'),
      } as unknown as import('ioredis').Redis,
    };

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        UpdaterService,
        { provide: KeydbService, useValue: keyDBService },
      ],
    }).compile();

    service = module.get<UpdaterService>(UpdaterService);
    service['wallets'] = mockWallets;
  });

  it('should be defined', () => {
    expect(service).toBeDefined();
  });

  it('should filter and save correct addresses to KeyDB', async () => {
    await service.handleUpdate();

    expect(keyDBService.client?.set).toHaveBeenCalledWith(
      'whale_whitelist',
      JSON.stringify(['0x1']),
    );
  });

  it('should skip update if no traders meet criteria', async () => {
    service['wallets'] = mockWallets.map((t) => ({
      ...t,
      account_value: 1000,
    }));

    await service.handleUpdate();

    expect(keyDBService.client?.set).not.toHaveBeenCalled();
  });
});