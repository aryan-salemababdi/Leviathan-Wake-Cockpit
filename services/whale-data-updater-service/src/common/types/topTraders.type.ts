export interface TopTraders {
  address: string;
  pnl: number;
  account_value: number;
  main_position: {
    coin: string;
    value: number;
    side: string;
  };
  direction_bias: number;
  perp_day_pnl: number;
  perp_week_pnl: number;
  perp_month_pnl: number;
  perp_alltime_pnl: number;
}
